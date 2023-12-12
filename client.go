package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"slices"
	"strings"
	"time"
)

const VERSION = "0.1.0"

var AccessTokenRequiredError = errors.New("access token is required")
var RetriableHttpStatuses = []int{408, 409, 429, 502, 503, 504}

type BadRequestError struct {
	response []byte
}

func (err BadRequestError) Error() string {
	return fmt.Sprintf("the provided JSON is not valid. provided body is:%s", string(err.response))
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	config *Config
	client HTTPClient
	url    *string
}

type QueryRequest struct {
	Query         string
	OperationName string
	Variables     map[string]any
	headers       http.Header
}

func NewClient(config *Config) *Client {
	c := &Client{
		config: config,
	}

	if c.config.GetTimeout() < (time.Duration(1) * time.Second) {
		c.config.SetTimeout(time.Duration(30) * time.Second)
	}
	c.SetHttpClient(&http.Client{Timeout: config.GetTimeout() * time.Second})
	return c
}

func (c *Client) SetHttpClient(httpClient HTTPClient) {
	c.client = httpClient
}
func (c *Client) Query(ctx context.Context, r *QueryRequest) (*Response, error) {

	if c.config.GetAccessToken() == "" {
		return nil, AccessTokenRequiredError
	}

	var requestBody = make(map[string]any)
	requestBody["operationName"] = r.OperationName
	requestBody["variables"] = r.Variables
	requestBody["query"] = r.Query
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.getRequestUrl(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header = r.headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.GetAccessToken())
	req.Header.Set("X-Sdk-Name", "Go SDK")
	req.Header.Set("X-Sdk-Version", VERSION)
	req.Header.Set("X-Sdk-Lang-Version", runtime.Version())
	req.Header.Set("User-Agent", getUserAgent())

	var start = time.Now()
	// Retry loop
	var wait = 250.0
	var waitRate = 1.4
	var resp *http.Response
	var retries = uint(0)

	if c.config.maxRetries < 1 {
		c.config.SetMaxRetries(5)
	}
	for ; retries < c.config.maxRetries; retries++ {
		resp, err = c.client.Do(req)
		if err != nil || slices.Contains(RetriableHttpStatuses, resp.StatusCode) {
			time.Sleep(time.Duration(wait*waitRate) * time.Millisecond)
			continue
		}
		break
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var elapsed = time.Since(start)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rateLimiter = NewRateLimiter()
	rateLimiter.fill(resp)

	response := &Response{
		Status:         resp.StatusCode,
		Headers:        resp.Header,
		Body:           string(body),
		RequestHeaders: resp.Request.Header,
		RequestBody:    string(jsonBody),
		Url:            resp.Request.URL.String(),
		RetryCount:     int(retries) - 1,
		TimeSpentInMs:  elapsed,
		RateLimiter:    rateLimiter,
	}
	if response.Status >= http.StatusBadRequest && response.Status <= http.StatusInternalServerError {
		errorResponse, err := c.errorResponse(body)
		if err != nil {
			return response, err
		}
		response.ErrorResponse = errorResponse
	}

	return response, nil
}

func (c *Client) getRequestUrl() string {
	if c.url != nil {
		return *c.url
	}
	var path = "/graphql"
	var url string
	if strings.HasPrefix(c.config.GetAccessToken(), "sk_live") {
		url = "https://public.api.socio.events"
	} else {
		url = "https://public.sandbox-api.socio.events"
	}
	uri := url + path
	c.url = &uri
	return uri
}

func (c *Client) errorResponse(body []byte) (*ErrorResponse, error) {
	errorResponse := NewErrorResponse()
	if json.Valid(body) {
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, nil
	} else {
		return nil, BadRequestError{response: body}
	}
}
