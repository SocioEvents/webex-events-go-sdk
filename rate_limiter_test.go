package main

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_fillRateLimiter(t *testing.T) {
	var mockResponse = http.Response{
		Header: http.Header{
			"Content-Type":           []string{"application/json"},
			"X-Daily-Call-Limit":     []string{"10/200"},
			"X-Secondly-Call-Limit":  []string{"5/50"},
			"X-Daily-Retry-After":    []string{"3"},
			"X-Secondly-Retry-After": []string{"5"},
		},
	}

	var rateLimiter = RateLimiter{}
	rateLimiter.fill(&mockResponse)
	assert.Equal(t, rateLimiter.DailyBasedCostThreshold, 200)
	assert.Equal(t, rateLimiter.UsedDailyBasedCost, 10)
	assert.Equal(t, rateLimiter.SecondBasedCostThreshold, 50)
	assert.Equal(t, rateLimiter.UsedSecondBasedCost, 5)
	assert.Equal(t, rateLimiter.DailyRetryAfterInSecond, 3)
	assert.Equal(t, rateLimiter.SecondlyRetryAfterInMs, 5)
}
