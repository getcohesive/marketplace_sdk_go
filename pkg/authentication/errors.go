package authentication

import (
	"github.com/getcohesive/marketplace_sdk_go/pkg/common/errors"
)

func AuthError(message string) error {
	return errors.CohesiveError{Message: "AuthenticationError: " + message}
}
