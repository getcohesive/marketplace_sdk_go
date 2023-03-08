package usage

import (
	"errors"
	"fmt"

	"github.com/getcohesive/marketplace_sdk_go/cohesive_marketplace_sdk/pkg/request"
)


type UsageParams struct {
	request.BaseParams
	WorkspaceId    int    `json:"workspace_id"`
	InstanceId     int    `json:"instance_id"`
	Units          int    `json:"units"`
	Timestamp      int    `json:"timestamp"`
}

type Usage interface {
	Report(params UsageParams, idempotencyKey string) error
}

func NewUsage(client request.HTTPClient) Usage {
	return &usageClient{client}
}

type usageClient struct {
	client request.HTTPClient
}

func (u *usageClient) Report(params UsageParams, idempotencyKey string) error {
	if idempotencyKey == "" {
		return errors.New("idempotency key is required")
	}
	params.IdempotencyKey = idempotencyKey
	response, err := u.client.Request("POST", "/report-usage", params)
	if err != nil {
		fmt.Printf("Report Usage Failed : res = %s, err = %e", string(response), err)
		return err
	}
	return nil
}
