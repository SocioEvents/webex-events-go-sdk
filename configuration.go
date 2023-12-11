package main

import "time"

var accessToken string
var timeout = time.Duration(30)
var maxRetries int

func SetAccessToken(token string) {
	accessToken = token
}

func SetTimeout(_timeout time.Duration) {
	timeout = _timeout
}

func SetMaxRetries(_maxRetries int) {
	maxRetries = _maxRetries
}

func GetAccessToken() string {
	return accessToken
}

func GetTimeout() time.Duration {
	return timeout
}

func GetMaxRetries() int {
	return maxRetries
}
