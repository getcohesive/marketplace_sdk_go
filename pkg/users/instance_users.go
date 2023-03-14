package users

import (
	"encoding/json"
	"fmt"
	"github.com/getcohesive/marketplace_sdk_go/pkg/request"
)

type InstanceUsersParams struct {
	WorkspaceId int `json:"workspace_id"`
	InstanceId  int `json:"instance_id"`
}

type ListInstanceUsersResponse struct {
	Users []*struct {
		UserId               uint    `json:"user_id"`
		Name                 string  `json:"name"`
		Email                *string `json:"email"`
		ProfilePhoto         string  `json:"profile_photo"`
		InstanceMembershipId uint    `json:"instance_membership_id"`
		RoleId               uint    `json:"role_id"`
	} `json:"users"`
}

type Users interface {
	ListInstanceUsers(params InstanceUsersParams) (*ListInstanceUsersResponse, error)
}

func NewUsers(client request.HTTPClient) Users {
	return &usersClient{client}
}

type usersClient struct {
	client request.HTTPClient
}

func (u *usersClient) ListInstanceUsers(params InstanceUsersParams) (*ListInstanceUsersResponse, error) {
	response, err := u.client.Request("POST", "/list-instance-users", params)
	if err != nil {
		fmt.Printf("Report Usage Failed : res = %s, err = %e", string(response), err)
		return nil, err
	}

	users := &ListInstanceUsersResponse{}
	err = json.Unmarshal(response, users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
