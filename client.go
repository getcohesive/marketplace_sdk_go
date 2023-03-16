package cohesive_marketplace_sdk

import (
	"github.com/getcohesive/marketplace_sdk_go/pkg/authentication"
	"github.com/getcohesive/marketplace_sdk_go/pkg/common/errors"
	"github.com/getcohesive/marketplace_sdk_go/pkg/users"
	"net/url"
	"os"

	"github.com/getcohesive/marketplace_sdk_go/pkg/request"
	"github.com/getcohesive/marketplace_sdk_go/pkg/usage"
)

type client struct {
	config     *Config
	httpClient request.HTTPClient
}

func (c *client) Usage() usage.Usage {
	return usage.NewUsage(c.httpClient)
}

func (c *client) Users() users.Users {
	return users.NewUsers(c.httpClient)
}

func (c *client) ValidateToken(token string) (*authentication.AuthDetails, error) {
	return authentication.ValidateToken(token, c.config.CohesiveAppSecret)
}

type Client interface {
	Usage() usage.Usage
	Users() users.Users
	ValidateToken(token string) (*authentication.AuthDetails, error)
}

func NewClient(config *Config) (Client, error) {
	if config == nil {
		config = &Config{
			Config: request.Config{
				CohesiveApiKey: os.Getenv("COHESIVE_API_KEY"),
			},
			CohesiveAppSecret: os.Getenv("COHESIVE_APP_SECRET"),
			CohesiveAppID:     os.Getenv("COHESIVE_APP_ID"),
		}
		baseURL, err := url.Parse(os.Getenv("COHESIVE_BASE_URL"))
		if err != nil {
			return nil, errors.CohesiveError{Message: "Bad COHESIVE_BASE_URL"}
		}
		config.CohesiveBaseURL = baseURL
		if err := config.Validate(); err != nil {
			return nil, err
		}
	}
	httpClient, err := request.NewHTTPClient(&config.Config)
	if err != nil {
		return nil, err
	}
	return &client{
		config:     config,
		httpClient: httpClient,
	}, nil
}
