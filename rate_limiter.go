package main

import (
	"net/http"
	"strconv"
	"strings"
)

const DailyCallLimit = "X-Daily-Call-Limit"
const SecondlyCallLimit = "X-Secondly-Call-Limit"
const DailyRetryAfter = "X-Daily-Retry-After"
const SecondlyRetryAfter = "X-Secondly-Retry-After"

type RateLimiter struct {
	usedSecondBasedCost      int
	secondBasedCostThreshold int
	usedDailyBasedCost       int
	dailyBasedCostThreshold  int
	dailyRetryAfter          int
	secondlyRetryAfter       int
}

func FillRateLimiter(resp *http.Response, rateLimiter *RateLimiter) {
	var dailyCallLimit = resp.Header.Get(DailyCallLimit)
	if len(dailyCallLimit) > 0 {
		var limit = strings.Split(dailyCallLimit, "/")
		rateLimiter.usedDailyBasedCost, _ = strconv.Atoi(limit[0])
		rateLimiter.dailyBasedCostThreshold, _ = strconv.Atoi(limit[1])
	}

	var secondlyCallLimit = resp.Header.Get(SecondlyCallLimit)
	if len(secondlyCallLimit) > 0 {
		var limit = strings.Split(secondlyCallLimit, "/")
		rateLimiter.usedSecondBasedCost, _ = strconv.Atoi(limit[0])
		rateLimiter.secondBasedCostThreshold, _ = strconv.Atoi(limit[1])
	}

	var dailyRetryAfter = resp.Header.Get(DailyRetryAfter)
	if len(dailyRetryAfter) > 0 {
		rateLimiter.dailyRetryAfter, _ = strconv.Atoi(dailyRetryAfter)
	}

	var secondlyRetryAfter = resp.Header.Get(SecondlyRetryAfter)
	if len(secondlyRetryAfter) > 0 {
		rateLimiter.secondlyRetryAfter, _ = strconv.Atoi(secondlyRetryAfter)
	}
}
