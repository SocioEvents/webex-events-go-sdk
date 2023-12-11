package main

type RateLimiter struct {
	usedSecondBasedCost      int
	secondBasedCostThreshold int
	usedDailyBasedCost       int
	dailyBasedCostThreshold  int
	dailyRetryAfter          int
	secondlyRetryAfter       int
}
