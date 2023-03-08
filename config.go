package cohesive_marketplace_sdk

import (
	"errors"

	"github.com/getcohesive/marketplace_sdk_go/pkg/request"
)

type Config struct {
	request.Config
	CohesiveAppSecret string
	CohesiveAppID     string
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("empty config")
	}
	if c.CohesiveAppSecret == "" {
		return errors.New("empty app secret")
	}
	if c.CohesiveAppID == "" {
		return errors.New("empty app ID")
	}
	if err := c.Config.Validate(); err != nil {
		return err
	}
	return nil
}
