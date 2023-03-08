package cohesive_marketplace_sdk

import (
	"github.com/getcohesive/marketplace_sdk_go/pkg/common/errors"
	"github.com/getcohesive/marketplace_sdk_go/pkg/request"
)

type Config struct {
	request.Config    `json:"HTTP" yaml:"HTTP" mapstructure:"HTTP"`
	CohesiveAppSecret string `json:"COHESIVE_APP_SECRET" yaml:"COHESIVE_APP_SECRET" mapstructure:"COHESIVE_APP_SECRET"`
	CohesiveAppID     string `json:"COHESIVE_APP_ID" yaml:"COHESIVE_APP_ID" mapstructure:"COHESIVE_APP_ID"`
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
