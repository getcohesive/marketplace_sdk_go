package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getcohesive/marketplace_sdk_go/pkg/common/errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	CohesiveBaseURL    *url.URL      `json:"COHESIVE_BASE_URL" yaml:"COHESIVE_BASE_URL" mapstructure:"COHESIVE_BASE_URL"`
	CohesiveApiKey     string        `json:"COHESIVE_API_KEY" yaml:"COHESIVE_API_KEY" mapstructure:"COHESIVE_API_KEY"`
	CohesiveApiTimeout time.Duration `json:"COHESIVE_API_TIMEOUT" yaml:"COHESIVE_API_TIMEOUT" mapstructure:"COHESIVE_API_TIMEOUT"`
}

func (c Config) Validate() error {
	if c.CohesiveBaseURL == nil {
		return errors.CohesiveError{Message: "empty base URL"}
	}
	if c.CohesiveApiKey == "" {
		return errors.CohesiveError{Message: "empty API key"}
	}
	if c.CohesiveApiTimeout == 0 {
		c.CohesiveApiTimeout = time.Second * 10
	}
	return nil
}

type HTTPClient interface {
	Request(method string, path string, body interface{}) ([]byte, error)
}

func NewHTTPClient(config *Config) (HTTPClient, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &httpClient{
		config: config,
	}, nil
}

type httpClient struct {
	config *Config
}

func (c *httpClient) Request(method string, path string, body interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, APIClientError(err.Error())
	}
	request, err := http.NewRequest(method, c.config.CohesiveBaseURL.ResolveReference(&url.URL{Path: path}).String(), io.NopCloser(bytes.NewReader(requestBody)))
	if err != nil {
		return nil, APIConnectionError(err.Error())
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.config.CohesiveApiKey)
	client := &http.Client{Timeout: c.config.CohesiveApiTimeout}
	response, err := client.Do(request)
	if err != nil {
		return nil, APIConnectionError(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body %e", err)
		}
	}(response.Body)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return nil, APIClientError("failed to read response" + err.Error())
	}
	if response.StatusCode >= 400 {
		return nil, APIError(fmt.Sprintf("failed with code %d and response %s ", response.StatusCode, buf.String()))
	}
	return buf.Bytes(), nil
}
