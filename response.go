package main

import "time"

type Response struct {
	status         int
	headers        map[string][]string
	body           string
	requestHeaders map[string][]string
	requestBody    string
	url            string
	retryCount     int
	timeSpentInMs  time.Duration
	rateLimiter    RateLimiter
}
