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
	UsedSecondBasedCost      int
	SecondBasedCostThreshold int
	UsedDailyBasedCost       int
	DailyBasedCostThreshold  int
	DailyRetryAfter          int
	SecondlyRetryAfter       int
}

func NewRateLimiter() *RateLimiter {
	return new(RateLimiter)
}

func (rateLimiter *RateLimiter) fill(resp *http.Response) {
	var dailyCallLimit = resp.Header.Get(DailyCallLimit)
	if len(dailyCallLimit) > 0 {
		var limit = strings.Split(dailyCallLimit, "/")
		rateLimiter.UsedDailyBasedCost, _ = strconv.Atoi(limit[0])
		rateLimiter.DailyBasedCostThreshold, _ = strconv.Atoi(limit[1])
	}

	var secondlyCallLimit = resp.Header.Get(SecondlyCallLimit)
	if len(secondlyCallLimit) > 0 {
		var limit = strings.Split(secondlyCallLimit, "/")
		rateLimiter.UsedSecondBasedCost, _ = strconv.Atoi(limit[0])
		rateLimiter.SecondBasedCostThreshold, _ = strconv.Atoi(limit[1])
	}

	var dailyRetryAfter = resp.Header.Get(DailyRetryAfter)
	if len(dailyRetryAfter) > 0 {
		rateLimiter.DailyRetryAfter, _ = strconv.Atoi(dailyRetryAfter)
	}

	var secondlyRetryAfter = resp.Header.Get(SecondlyRetryAfter)
	if len(secondlyRetryAfter) > 0 {
		rateLimiter.SecondlyRetryAfter, _ = strconv.Atoi(secondlyRetryAfter)
	}
}
