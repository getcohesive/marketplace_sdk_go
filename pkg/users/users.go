package users

import "github.com/getcohesive/marketplace_sdk_go/pkg/request"

type Users interface {
	ListInstanceUsers(params ListInstanceUsersParams) (*ListInstanceUsersResponse, error)
	InviteUsers(params InviteUsersRequest) (*InviteUsersResponse, error)
}

func NewUsers(client request.HTTPClient) Users {
	return &usersClient{client}
}

type usersClient struct {
	client request.HTTPClient
}
