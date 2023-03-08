package request

import "github.com/getcohesive/marketplace_sdk_go/cohesive_marketplace_sdk/pkg/common"


func APIError(message string) error {
	return common.CohesiveError{Message: "APIError: "+message}
}

func APIConnectionError(message string) error {
	return common.CohesiveError{Message: "APIConnectionError: "+message}
}
func APIClientError(message string) error {
	return common.CohesiveError{Message: "APIClientError: "+message}
}
