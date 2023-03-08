package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

)

type Config struct{
	CohesiveBaseURL *url.URL
	CohesiveApiKey string
	CohesiveApiTimeout time.Duration
}

func (c Config) Validate() error {
	if c.CohesiveBaseURL == nil {
		return errors.New("empty base URL")
	}
	if c.CohesiveApiKey == "" {
		return errors.New("empty API key")
	}
	if c.CohesiveApiTimeout == 0{
		c.CohesiveApiTimeout = time.Second * 10
	}
	return nil
}

type HTTPClient interface {
	Request(method string, path string, body interface{}) ([]byte, error)
}
func NewHTTPClient(config *Config) (HTTPClient, error){
	if err := config.Validate(); err != nil{
		return nil, err
	}
	return &httpClient{
		config: config,
	}, nil
}


type httpClient struct{
	config *Config
}

func (c *httpClient) Request(method string, path string, body interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, APIClientError(err.Error())
	}
	request := &http.Request{
		Method: method,
		URL: c.config.CohesiveBaseURL.ResolveReference(&url.URL{Path: "/report-usage"}),
		Body: io.NopCloser(bytes.NewReader(requestBody)),
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.config.CohesiveApiKey)
	client := &http.Client{Timeout: c.config.CohesiveApiTimeout}
	response, err := client.Do(request)
	if err != nil {
		return nil, APIConnectionError(err.Error())
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil{
		return nil, APIClientError("failed to read response"+err.Error())
	}
	if response.StatusCode >= 400 {
		return nil, APIError(fmt.Sprintf("failed with code %d and response %s ", response.StatusCode,  buf.String()))
	}
	return buf.Bytes(), nil
}