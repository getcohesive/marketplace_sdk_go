package users

import (
	"encoding/json"
	"fmt"
)

type InviteUsersRequest struct {
	WorkspaceId   int      `json:"workspace_id"`
	InstanceId    int      `json:"instance_id"`
	UserId        int      `json:"user_id"`
	InvitedEmails []string `json:"invited_emails"`
}

type InviteUsersResponse struct {
	SuccessfulInvitations []string `json:"successful_invitations"`
	AddedToInstance       []string `json:"added_to_instance"`
	FailedInvitations     []struct {
		Email     string    `json:"email"`
		Error     string    `json:"error"`
		ErrorCode ErrorCode `json:"error_code"`
	} `json:"failed_invitations"`
}

type ErrorCode int

const (
	UserAddedToInstance        ErrorCode = 1
	UserPartOfAnotherWorkspace ErrorCode = 2
	InvalidEmailAddress        ErrorCode = 3
	PersonalWorkspaceError     ErrorCode = 4
	PrivateWorkspaceError      ErrorCode = 5
	AlreadyAMember             ErrorCode = 6
	AlreadyInvited             ErrorCode = 7
)

func (e ErrorCode) String() string {
	switch e {
	case UserAddedToInstance:
		return "user added to instance"
	case UserPartOfAnotherWorkspace:
		return "user part of another workspace"
	case InvalidEmailAddress:
		return "invalid email address"
	case PersonalWorkspaceError:
		return "invitee email should be a personal email since current workspace is a personal workspace"
	case PrivateWorkspaceError:
		return "you can only invite email belongs to current workspace domain"
	case AlreadyAMember:
		return "user already a member in current instance"
	case AlreadyInvited:
		return "email already invited to current instance"
	default:
		return "unknown error code"
	}
}

func (u *usersClient) InviteUsers(params InviteUsersRequest) (*InviteUsersResponse, error) {
	response, err := u.client.Request("POST", "/list-instance-users", params)
	if err != nil {
		fmt.Printf("Invite users to instance failed : res = %s, err = %e", string(response), err)
		return nil, err
	}

	inviteUsersResponse := &InviteUsersResponse{}
	err = json.Unmarshal(response, inviteUsersResponse)
	if err != nil {
		return nil, err
	}
	return inviteUsersResponse, nil
}
