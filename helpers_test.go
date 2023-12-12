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

	assert.Equal(t, want, true)

	// From cache

	got = getUserAgent()
	want = strings.HasPrefix(got, "Webex Go SDK")
	assert.Equal(t, want, true)
}

func TestGetRequestUrl(t *testing.T) {
	SetAccessToken("sk_test_some_access_token1")
	var got = GetRequestUrl()
	var want = "https://public.sandbox-api.socio.events/graphql"

	assert.Equal(t, got, want)

	// Retry again. The second call will retrieve from cache.
	got = GetRequestUrl()
	want = "https://public.sandbox-api.socio.events/graphql"

	assert.Equal(t, got, want)

	SetAccessToken("sk_live_some_access_token1")
	got = GetRequestUrl()
	want = "https://public.api.socio.events/graphql"

	assert.Equal(t, got, want)
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
	var response = Response{
		Status: 400,
		Body:   json,
	}

	errorResponse := ErrorResponse{}

	var err = fillErrorResponse(&response, &errorResponse)
	assert.Nil(t, err)
	assert.Equal(t, errorResponse.Message, "Something went wrong")

}

func Test_fillErrorResponseIfStatusCodeIs200(t *testing.T) {
	var json = `
	{
		"data": {
			"currenciesList": {
				"isoCode": "USD"
			}
		}
	}
`
	var response = Response{
		Status: 200,
		Body:   json,
	}

	errorResponse := ErrorResponse{}
	var p1 = &errorResponse

	var err = fillErrorResponse(&response, &errorResponse)
	var p2 = response.ErrorResponse

	assert.Nil(t, err)
	assert.NotEqual(t, p1, p2)
}

func Test_fillErrorResponseIfJSONisInvalid(t *testing.T) {
	var json = "malformed json"
	var response = Response{
		Status: 400,
		Body:   json,
	}

	errorResponse := ErrorResponse{}

	var err = fillErrorResponse(&response, &errorResponse)
	assert.NotNil(t, err)
	assert.Errorf(t, err, "The provided JSON is not valid. Provided body is: malformed json")
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

	var response = Response{
		Status: 400,
		Body:   json,
	}

	errorResponse := ErrorResponse{}

	var err = fillErrorResponse(&response, &errorResponse)
	assert.Nil(t, err)
	assert.Equal(t, response.ErrorResponse, errorResponse)
	assert.Equal(t, errorResponse.Message, "Something went wrong")
	assert.Equal(t, errorResponse.Extensions.Code, "TOKEN_IS_EXPIRED")
	assert.Equal(t, errorResponse.Extensions.ReferenceId, "reference id")
	assert.Equal(t, errorResponse.Extensions.Cost, 51)
	assert.Equal(t, errorResponse.Extensions.AvailableCost, 50)
	assert.Equal(t, errorResponse.Extensions.Threshold, 50)
	assert.Equal(t, errorResponse.Extensions.DailyThreshold, 200)
	assert.Equal(t, errorResponse.Extensions.DailyAvailableCost, 190)

	_, ok := errorResponse.Extensions.Errors["first_name"]
	assert.Equal(t, ok, true)

	msg, ok := errorResponse.GraphqlErrors[0]["message"]
	assert.Equal(t, ok, true)
	assert.Equal(t, msg, "graphql error")
}
