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

type BadResponseError struct {
	response []byte
}

func (err BadResponseError) Error() string {
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
	Query          string
	OperationName  string
	Variables      map[string]any
	IdempotencyKey string
}

func NewClient(config *Config) *Client {
	c := &Client{
		config: config,
	}

	c.SetHttpClient(&http.Client{Timeout: config.GetTimeout() * time.Second})
	return c
}

func (c *Client) SetHttpClient(httpClient HTTPClient) {
	c.client = httpClient
}

func (c *Client) DoIntrospectionQuery(ctx context.Context) (*Response, error) {
	client := NewClient(c.config)
	var variables = make(map[string]any)
	var r = QueryRequest{
		OperationName: "IntrospectionQuery",
		Variables:     variables,
		Query:         getIntrospectionQuery(),
	}

	var response, err = client.Query(ctx, &r)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) Query(ctx context.Context, queryRequest *QueryRequest) (*Response, error) {
	if c.config.GetAccessToken() == "" {
		return nil, AccessTokenRequiredError
	}

	var requestBody = make(map[string]any)
	requestBody["operationName"] = queryRequest.OperationName
	requestBody["variables"] = queryRequest.Variables
	requestBody["query"] = queryRequest.Query
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		c.config.logger.Error(err.Error())
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.getRequestUrl(), bytes.NewBuffer(jsonBody))
	if err != nil {
		c.config.logger.Error(err.Error())
		return nil, err
	}

	req.Header = http.Header{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.GetAccessToken())
	req.Header.Set("X-Sdk-Name", "Go SDK")
	req.Header.Set("X-Sdk-Version", VERSION)
	req.Header.Set("X-Sdk-Lang-Version", runtime.Version())
	req.Header.Set("User-Agent", getUserAgent())

	if len(queryRequest.IdempotencyKey) > 0 {
		req.Header.Set("Idempotency-Key", queryRequest.IdempotencyKey)
	}
	var (
		start    = time.Now()
		wait     = 250.0
		waitRate = 1.4
		resp     *http.Response
		retries  = uint(0)
	)

	c.config.logger.Info("The request to" + c.getRequestUrl() + " endpoint has been started.")
	resp, err = c.client.Do(req)
	if isErrorOrServerError(err, resp) {
		if c.config.maxRetries > 0 {
			if err != nil {
				c.config.logger.Error("The request is going to be retried. The error message is: " + err.Error())
			} else {
				c.config.logger.Error("The request is going to be retried due to the fact that the server returned " + resp.Status + " status code.")
			}
		}

		// Retry loop
		for ; retries < c.config.maxRetries; retries++ {
			wait *= waitRate
			c.config.logger.Info(fmt.Sprintf("Sleeping for %d ms...", int(wait)))
			time.Sleep(time.Duration(wait) * time.Millisecond)
			resp, err = c.client.Do(req)
			if isErrorOrServerError(err, resp) {
				if (retries + 1) < c.config.maxRetries { // Has another retry.
					if err != nil {
						c.config.logger.Error("The request is going to be retried. The error message is: " + err.Error())
					} else {
						c.config.logger.Error("The request is going to be retried due to the fact that the server returned " + resp.Status + " status code.")
					}
				}
				continue
			}
			break
		}
	}

	if resp != nil && resp.StatusCode > 299 {
		c.config.logger.Error(fmt.Sprintf("The request is failed. The request retried %d times but server returned %d status code", retries, resp.StatusCode))
	}

	if err != nil {
		c.config.logger.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var elapsed = time.Since(start)

	if err != nil {
		c.config.logger.Error(err.Error())
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.config.logger.Error(err.Error())
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
		RetryCount:     int(retries),
		TimeSpentInMs:  elapsed,
		RateLimiter:    rateLimiter,
	}
	if response.Status >= http.StatusBadRequest && response.Status <= http.StatusInternalServerError {
		errorResponse, err := c.errorResponse(body)
		if err != nil {
			c.config.logger.Error(err.Error())
			return response, err
		}
		response.ErrorResponse = errorResponse
	}

	return response, nil
}

func isErrorOrServerError(err error, resp *http.Response) bool {
	return err != nil || slices.Contains(RetriableHttpStatuses, resp.StatusCode)
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
		return nil, BadResponseError{response: body}
	}
}
