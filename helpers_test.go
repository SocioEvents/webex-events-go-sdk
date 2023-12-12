package main

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_getUserAgent(t *testing.T) {
	var got = getUserAgent()
	var want = strings.HasPrefix(got, "Webex Go SDK")

	assert.EqualValues(t, want, true)

	// From cache

	got = getUserAgent()
	want = strings.HasPrefix(got, "Webex Go SDK")
	assert.EqualValues(t, want, true)
}

func TestGetRequestUrl(t *testing.T) {
	config := NewConfig()
	config.SetAccessToken("sk_test_some_access_token1")
	client := NewClient(config)
	var got = client.getRequestUrl()
	var want = "https://public.sandbox-api.socio.events/graphql"

	assert.EqualValues(t, got, want)

	// Retry again. The second call will retrieve from cache.
	got = client.getRequestUrl()
	want = "https://public.sandbox-api.socio.events/graphql"

	assert.EqualValues(t, got, want)

	config = NewConfig()
	config.SetAccessToken("sk_live_some_access_token1")
	client = NewClient(config)

	got = client.getRequestUrl()
	want = "https://public.api.socio.events/graphql"

	assert.EqualValues(t, got, want)
}

func Test_fillErrorResponseIfStatusCodeIsInRange(t *testing.T) {
	var json = `
	{
		"message": "Something went wrong",
		"extensions": {
			"code": "TOKEN_IS_EXPIRED"
		}
	}
`
	config := NewConfig()
	config.SetAccessToken("sk_live_some_access_token")
	client := NewClient(config)

	errorResponse, err := client.errorResponse([]byte(json))
	assert.NoError(t, err)
	assert.EqualValues(t, errorResponse.Message, "Something went wrong")

}

func Test_fillErrorResponseIfJSONisInvalid(t *testing.T) {
	var json = "malformed json"
	config := NewConfig()
	config.SetAccessToken("sk_live_some_access_token")
	client := NewClient(config)

	_, err := client.errorResponse([]byte(json))
	assert.Error(t, err, BadResponseError{
		response: []byte("malformed json"),
	})
	//assert.Errorf(t, err, "The provided JSON is not valid. Provided body is: malformed json")
}

func Test_fillErrorResponseWithAllValues(t *testing.T) {
	var json = `
		{
			"message": "Something went wrong",
			"extensions": {
				"code": "TOKEN_IS_EXPIRED",
				"referenceId": "reference id",
				"errors": {
					"first_name": [
						"invalid"
					]
				},
				"cost": 51,
				"availableCost": 50,
				"threshold": 50,
				"dailyThreshold": 200,
				"dailyAvailableCost": 190
			},
			"errors": [
				{
					"message": "graphql error"
				}
			]
		}
`
	config := NewConfig()
	config.SetAccessToken("sk_live_some_access_token")
	client := NewClient(config)

	errorResponse, err := client.errorResponse([]byte(json))
	assert.NoError(t, err)
	assert.EqualValues(t, errorResponse.Message, "Something went wrong")
	assert.EqualValues(t, errorResponse.Extensions.Code, "TOKEN_IS_EXPIRED")
	assert.EqualValues(t, errorResponse.Extensions.ReferenceId, "reference id")
	assert.EqualValues(t, errorResponse.Extensions.Cost, 51)
	assert.EqualValues(t, errorResponse.Extensions.AvailableCost, 50)
	assert.EqualValues(t, errorResponse.Extensions.Threshold, 50)
	assert.EqualValues(t, errorResponse.Extensions.DailyThreshold, 200)
	assert.EqualValues(t, errorResponse.Extensions.DailyAvailableCost, 190)

	_, ok := errorResponse.Extensions.Errors["first_name"]
	assert.EqualValues(t, ok, true)

	msg, ok := errorResponse.GraphqlErrors[0]["message"]
	assert.EqualValues(t, ok, true)
	assert.EqualValues(t, msg, "graphql error")
}
