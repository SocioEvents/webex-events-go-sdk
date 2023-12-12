package main

import "time"

type Response struct {
	Status         int
	Headers        map[string][]string
	Body           string
	RequestHeaders map[string][]string
	RequestBody    string
	Url            string
	RetryCount     int
	TimeSpentInMs  time.Duration
	RateLimiter    RateLimiter
}
