package request

import (
	"github.com/getcohesive/marketplace_sdk_go/pkg/common/errors"
)

func APIError(message string) error {
	return errors.CohesiveError{Message: "APIError: " + message}
}

func APIConnectionError(message string) error {
	return errors.CohesiveError{Message: "APIConnectionError: " + message}
}
func APIClientError(message string) error {
	return errors.CohesiveError{Message: "APIClientError: " + message}
}
