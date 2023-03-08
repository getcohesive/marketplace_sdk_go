package dist_lock

import (
	"fmt"
	"github.com/getcohesive/marketplace_sdk_go/pkg/common/errors"
	"time"
)

var (
	LockFailedError = errors.CohesiveError{Message: "lock acquiring failed"}
)

type JobLockRepository interface {
	Save(model *JobLock) error
	DeleteLockIfExpired(key string, expiryInSeconds int64) (*JobLock, error)
	Delete(key string) error
}

type JobLock struct {
	// JobID make sure that this is unique in the db
	JobID string
	TS    time.Time
}

type Client interface {
	RunJobWithLock(
		lockKey string,
		jobFunction func(),
		lockExpiryInSecs int64,
		durationBetweenJobsInSecs int64,
		lockRetryTimeInSecs int64,
	)
	AcquireLock(lockKey string, lockExpiryInSeconds int64) (*JobLock, error)
	ReleaseLock(lock *JobLock) error
	RefreshLock(lock *JobLock) error
}

type client struct {
	JobLockRepository
}

func NewClient(jobLockRepository JobLockRepository) Client {
	return &client{JobLockRepository: jobLockRepository}
}

func (client *client) AcquireLock(lockKey string, lockExpiryInSeconds int64) (*JobLock, error) {
	lock := &JobLock{JobID: lockKey, TS: time.Now()}
	if err := client.Save(lock); err == nil {
		return lock, nil
	}
	lock, err := client.DeleteLockIfExpired(lockKey, lockExpiryInSeconds)
	if err != nil {
		return nil, LockFailedError
	}
	if time.Since(lock.TS).Seconds() > float64(lockExpiryInSeconds) {
		if err := client.Delete(lock.JobID); err != nil {
			return nil, LockFailedError
		}
		return client.AcquireLock(lockKey, lockExpiryInSeconds)
	}
	return nil, LockFailedError
}

func (client *client) RefreshLock(lock *JobLock) error {
	lock.TS = time.Now()
	return client.Save(lock)
}

func (client *client) ReleaseLock(lock *JobLock) error {
	return client.Delete(lock.JobID)
}

func (client *client) RunJobWithLock(
	lockKey string,
	jobFunction func(),
	lockExpiryInSecs int64,
	durationBetweenJobsInSecs int64,
	lockRetryTimeInSecs int64,
) {
	for {
		fmt.Println("CH_JOB_WITH_LOCK: checking for lock")
		lock, err := client.AcquireLock(lockKey, lockExpiryInSecs)
		if err != nil {
			fmt.Printf("CH_JOB_WITH_LOCK: lock busy %e. sleeping for %d secs", err, lockRetryTimeInSecs)
			time.Sleep(time.Duration(lockRetryTimeInSecs) * time.Second)
			continue
		}
		if lock != nil {
			fmt.Println("CH_JOB_WITH_LOCK: lock acquired")
			for {
				fmt.Println("CH_JOB_WITH_LOCK: refreshing lock")
				if err := client.RefreshLock(lock); err != nil {
					fmt.Printf("WARNING: CH_JOB_WITH_LOCK: job lock refresh failed. %e", err)
				}
				jobFunction()
				fmt.Println("CH_JOB_WITH_LOCK: job complete")
				time.Sleep(time.Duration(durationBetweenJobsInSecs) * time.Second)
			}
		}
	}
}
