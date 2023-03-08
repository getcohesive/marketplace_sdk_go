package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	cohesiveMarketplaceSDK "github.com/getcohesive/marketplace_sdk_go"
	"github.com/getcohesive/marketplace_sdk_go/pkg/dist_lock"
	"github.com/getcohesive/marketplace_sdk_go/pkg/request"
	"github.com/getcohesive/marketplace_sdk_go/pkg/usage"
)

var (
	NotFoundError      = fmt.Errorf("not found")
	AlreadyExistsError = fmt.Errorf("already exists")
)

func uniqueKey(workspaceId int, instanceId int) string {
	return fmt.Sprintf("%d,%d", workspaceId, instanceId)
}

func getWorkspaceIDAndInstanceIDFromUniqueKey(uniqueKey string) (int, int, error) {
	items := strings.Split(uniqueKey, ",")
	workspaceId, err := strconv.Atoi(items[0])
	if err != nil {
		return 0, 0, err
	}
	instanceId, err := strconv.Atoi(items[1])
	if err != nil {
		return workspaceId, 0, err
	}
	return workspaceId, instanceId, nil
}

type UsageTrackerRepository interface {
	Create(tracker *UsageTracker) error

	// FindAndUpdateCountBy
	// should increment used_items_count by @usageCountDiff for
	// all the records which match @uniqueKey
	// UPDATE TABLE SET used_items_count = used_items_count + @usageCountDiff WHERE unique_key = @uniqueKey
	FindAndUpdateCountBy(uniqueKey string, usageCountDiff int) error

	// FindAndUpdateReportedCountBy
	// should increment reported_item_count by @reportCountDiff for
	// all the records which match @uniqueKey
	// UPDATE TABLE SET reported_item_count = reported_item_count + @reportCountDiff WHERE unique_key = @uniqueKey
	FindAndUpdateReportedCountBy(uniqueKey string, reportCountDiff int) error

	// GetAllUnreportedUsageTrackers
	// should return all the usage trackers where used_items_count >  reported_items_count + free_trial_items
	GetAllUnreportedUsageTrackers() ([]*UsageTracker, error)
}

type UsageTracker struct {
	UniqueKey          string    `json:"unique_key"`
	UsedItemsCount     int       `json:"used_items_count"`
	ReportedItemsCount int       `json:"reported_items_count"`
	FreeTrialItems     int       `json:"free_trial_items"`
	ItemsPerUnit       int       `json:"items_per_unit"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Client struct {
	UsageTrackerRepository
	sdkClient      cohesiveMarketplaceSDK.Client
	distLockClient dist_lock.Client
}

func (c *Client) TrackUsage(
	workspaceId int,
	instanceId int,
	usedItems int,
	freeTrialItems int,
	itemsPerUnit int,
) error {
	uniqueKeyForUsage := uniqueKey(workspaceId, instanceId)
	lock, err := c.distLockClient.AcquireLock(uniqueKeyForUsage, 1)
	if err != nil {
		return err
	}

	defer func(distLockClient dist_lock.Client, lock *dist_lock.JobLock) {
		err := distLockClient.ReleaseLock(lock)
		if err != nil {
			fmt.Printf("ERROR: TrackUsage: release lock failed: %e", err)
		}
	}(c.distLockClient, lock)

	err = c.FindAndUpdateCountBy(uniqueKeyForUsage, usedItems)
	if err == nil || err != NotFoundError {
		return err
	}

	err = c.distLockClient.ReleaseLock(lock)
	if err != nil {
		fmt.Printf("ERROR: TrackUsage: refresh lock failed: %e", err)
		return err
	}
	usageTracker := &UsageTracker{
		UniqueKey:          uniqueKeyForUsage,
		UsedItemsCount:     usedItems,
		ReportedItemsCount: 0,
		FreeTrialItems:     freeTrialItems,
		ItemsPerUnit:       itemsPerUnit,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	err = c.Create(usageTracker)
	if err == nil {
		return nil
	}
	if err == AlreadyExistsError {
		err2 := c.FindAndUpdateCountBy(uniqueKeyForUsage, usedItems)
		return err2
	}
	return err
}

func (c *Client) ReportUsageOnce() {
	usageTrackers, err := c.GetAllUnreportedUsageTrackers()
	if err != nil {
		fmt.Printf("failed to fetch trackers to be updated %e", err)
		return
	}

	for _, usageTracker := range usageTrackers {
		paidUsage := usageTracker.UsedItemsCount - usageTracker.FreeTrialItems
		delta := paidUsage - usageTracker.ReportedItemsCount

		if delta > 0 {
			fmt.Printf("CH_USAGE_TRACKING: reporting for %s : %d\n", usageTracker.UniqueKey, (usageTracker.ReportedItemsCount/usageTracker.ItemsPerUnit)+1)

			workspaceId, instanceId, err := getWorkspaceIDAndInstanceIDFromUniqueKey(usageTracker.UniqueKey)
			if err != nil {
				fmt.Printf("WARNING: unique key parsing failed for key %s: %e", usageTracker.UniqueKey, err)
			}
			idempotencyKey := fmt.Sprintf("%d", (usageTracker.ReportedItemsCount/usageTracker.ItemsPerUnit)+1)
			// Report
			params := usage.Params{
				WorkspaceId: workspaceId,
				InstanceId:  instanceId,
				Units:       1,
				Timestamp:   int(time.Now().Unix() * 1000),
				BaseParams: request.BaseParams{
					IdempotencyKey: idempotencyKey,
				},
			}

			err = c.sdkClient.Usage().Report(params, idempotencyKey)

			if err == nil {
				err := c.FindAndUpdateReportedCountBy(usageTracker.UniqueKey, usageTracker.ItemsPerUnit)
				if err != nil {
					fmt.Printf("WARNING: failed to update reported data %s: %e", usageTracker.UniqueKey, err)
				}
			}
		}
	}
}

func StartReporting(client dist_lock.Client, usageTrackerClient Client) {
	go client.RunJobWithLock(
		"report_usage",
		usageTrackerClient.ReportUsageOnce,
		60*60,
		10,
		5*60,
	)
}
