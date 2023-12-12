package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

func TestWhenAccessTokenIsBlank(t *testing.T) {
	var variables = make(map[string]any)
	var headers = make(map[string]string)

	var query = "{ currenciesList{ isoCode }}"
	var operationName = "CurrenciesList"
	response, err := Query(query, operationName, variables, headers)

	assert.Equal(t, response, Response{})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Access Token is required.")
}

func TestQueryWith200StatusCode(t *testing.T) {
	var responseBody = `{"data": { "currenciesList": { "isoCode": "USD" }}}`
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(responseBody))
	}))
	defer func() { testServer.Close() }()

	SetAccessToken("sk_live_some_access_token")
	requestUrl[GetAccessToken()] = testServer.URL

	var variables = make(map[string]any)
	var headers = make(map[string]string)
	headers["Idempotency-Key"] = "61672155-56d3-4375-a864-52e7bba4f445"

	var query = "{ currenciesList{ isoCode }}"
	var operationName = "CurrenciesList"
	response, err := Query(query, operationName, variables, headers)

	var result map[string]any
	json.Unmarshal([]byte(response.RequestBody), &result)
	assert.Equal(t, result["operationName"], operationName)
	assert.Equal(t, result["query"], query)
	assert.Equal(t, result["variables"], variables)
	assert.Nil(t, err)
	assert.Equal(t, response.Status, http.StatusOK)
	assert.Equal(t, response.Body, responseBody)
	assert.Equal(t, response.RequestHeaders["Authorization"][0], "Bearer "+GetAccessToken())
	var exists = strings.HasPrefix(response.RequestHeaders["User-Agent"][0], "Webex Go SDK")
	assert.Equal(t, exists, true)
	assert.Equal(t, response.RequestHeaders["Content-Type"][0], "application/json")
	assert.Equal(t, response.RequestHeaders["X-Sdk-Name"][0], "Go SDK")
	assert.Equal(t, response.RequestHeaders["X-Sdk-Version"][0], VERSION)
	assert.Equal(t, response.RequestHeaders["X-Sdk-Lang-Version"][0], runtime.Version())
	assert.Equal(t, response.RequestHeaders["Idempotency-Key"][0], "61672155-56d3-4375-a864-52e7bba4f445")
}
