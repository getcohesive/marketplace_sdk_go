package auth

import "github.com/getcohesive/marketplace_sdk_go/cohesive_marketplace_sdk/pkg/common"


func AuthenticationError ( message string) error {
	return common.CohesiveError{Message: "AuthenticationError: "+ message}
}
