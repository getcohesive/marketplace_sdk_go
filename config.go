package cohesive_marketplace_sdk

import (
	"github.com/getcohesive/marketplace_sdk_go/pkg/common/errors"
	"github.com/getcohesive/marketplace_sdk_go/pkg/request"
)

type Config struct {
	request.Config
	CohesiveAppSecret string
	CohesiveAppID     string
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.CohesiveError{Message: "empty config"}
	}
	if c.CohesiveAppSecret == "" {
		return errors.CohesiveError{Message: "empty app secret"}
	}
	if c.CohesiveAppID == "" {
		return errors.CohesiveError{Message: "empty app ID"}
	}
	if err := c.Config.Validate(); err != nil {
		return err
	}
	return nil
}
