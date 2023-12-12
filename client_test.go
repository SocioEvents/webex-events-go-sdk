package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"testing"
)

type TestHttpClient struct {
	mock.Mock
}

func (c *TestHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestWhenAccessTokenIsBlank(t *testing.T) {
	ctx := context.Background()
	config := NewConfig()
	client := NewClient(config)

	var variables = make(map[string]any)
	headers := http.Header{}
	var query = "{ currenciesList{ isoCode }}"
	var operationName = "CurrenciesList"

	response, err := client.Query(ctx, &QueryRequest{
		Query:         query,
		OperationName: operationName,
		Variables:     variables,
		headers:       headers,
	})

	assert.Nil(t, response)
	assert.ErrorIs(t, err, AccessTokenRequiredError)
}

func TestQueryWith409StatusCode(t *testing.T) {
	ctx := context.Background()
	config := NewConfig()
	config.SetAccessToken("sk_live_some_access_token")
	config.SetMaxRetries(3)
	client := NewClient(config)

	var headers = http.Header{}
	headers.Set("Idempotency-Key", "61672155-56d3-4375-a864-52e7bba4f445")
	var variables = make(map[string]any)

	var query = "{ currenciesList{ isoCode }}"
	var operationName = "CurrenciesList"

	var responseBody = `{ "message": "Conflict", "extensions": { code: "CONFLICT" }}`

	httpClient := new(TestHttpClient)
	stringReader := strings.NewReader(responseBody)
	body := io.NopCloser(stringReader)

	reqUrl, err := url.Parse("http://localhost")
	assert.NoError(t, err)

	header := http.Header{
		"Authorization":      []string{fmt.Sprintf("Bearer %s", config.GetAccessToken())},
		"User-Agent":         []string{"Webex Go SDK"},
		"Content-Type":       []string{"application/json"},
		"X-Sdk-Name":         []string{"Go SDK"},
		"X-Sdk-Version":      []string{VERSION},
		"X-Sdk-Lang-Version": []string{runtime.Version()},
		"Idempotency-Key":    []string{"61672155-56d3-4375-a864-52e7bba4f445"},
	}
	var mockResponse = &http.Response{
		Request: &http.Request{
			Header: header,
			URL:    reqUrl,
		},
		Status:     "400",
		StatusCode: 409,
		Body:       body,
	}

	httpClient.On("Do", mock.Anything).Return(mockResponse, nil)
	client.SetHttpClient(httpClient)

	response, err := client.Query(ctx, &QueryRequest{
		Query:         query,
		OperationName: operationName,
		Variables:     variables,
		headers:       headers,
	})

	assert.Equal(t, 2, response.RetryCount)
}

func TestQueryWith200StatusCode(t *testing.T) {
	ctx := context.Background()
	config := NewConfig()
	config.SetAccessToken("sk_live_some_access_token")
	client := NewClient(config)

	var headers = http.Header{}
	headers.Set("Idempotency-Key", "61672155-56d3-4375-a864-52e7bba4f445")
	var variables = make(map[string]any)

	var query = "{ currenciesList{ isoCode }}"
	var operationName = "CurrenciesList"

	var responseBody = `{"data": { "currenciesList": { "isoCode": "USD" }}}`

	httpClient := new(TestHttpClient)
	stringReader := strings.NewReader(responseBody)
	body := io.NopCloser(stringReader)

	reqUrl, err := url.Parse("http://localhost")
	assert.NoError(t, err)

	header := http.Header{
		"Authorization":      []string{fmt.Sprintf("Bearer %s", config.GetAccessToken())},
		"User-Agent":         []string{"Webex Go SDK"},
		"Content-Type":       []string{"application/json"},
		"X-Sdk-Name":         []string{"Go SDK"},
		"X-Sdk-Version":      []string{VERSION},
		"X-Sdk-Lang-Version": []string{runtime.Version()},
		"Idempotency-Key":    []string{"61672155-56d3-4375-a864-52e7bba4f445"},
	}
	var mockResponse = &http.Response{
		Request: &http.Request{
			Header: header,
			URL:    reqUrl,
		},
		Status:     "200",
		StatusCode: 200,
		Body:       body,
		Header: http.Header{
			"Content-Type":           []string{"application/json"},
			"X-Daily-Call-Limit":     []string{"10/200"},
			"X-Secondly-Call-Limit":  []string{"5/50"},
			"X-Daily-Retry-After":    []string{"3"},
			"X-Secondly-Retry-After": []string{"5"},
		},
	}

	httpClient.On("Do", mock.Anything).Return(mockResponse, nil)
	client.SetHttpClient(httpClient)

	response, err := client.Query(ctx, &QueryRequest{
		Query:         query,
		OperationName: operationName,
		Variables:     variables,
		headers:       headers,
	})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	result := make(map[string]any)
	err = json.Unmarshal([]byte(response.RequestBody), &result)
	assert.NoError(t, err)
	assert.EqualValues(t, result["operationName"], operationName)
	assert.EqualValues(t, result["query"], query)
	assert.EqualValues(t, result["variables"], variables)
	assert.NoError(t, err)
	assert.EqualValues(t, response.Status, http.StatusOK)
	assert.EqualValues(t, response.Body, responseBody)
	assert.EqualValues(t, response.RequestHeaders["Authorization"][0], "Bearer "+config.GetAccessToken())
	assert.Contains(t, response.RequestHeaders["User-Agent"][0], "Webex Go SDK")
	assert.EqualValues(t, response.RequestHeaders["Content-Type"][0], "application/json")
	assert.EqualValues(t, response.RequestHeaders["X-Sdk-Name"][0], "Go SDK")
	assert.EqualValues(t, response.RequestHeaders["X-Sdk-Version"][0], VERSION)
	assert.EqualValues(t, response.RequestHeaders["X-Sdk-Lang-Version"][0], runtime.Version())
	assert.EqualValues(t, response.RequestHeaders["Idempotency-Key"][0], "61672155-56d3-4375-a864-52e7bba4f445")
}
