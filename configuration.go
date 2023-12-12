package main

import "time"

var (
	accessToken string
	timeout     = time.Duration(30)
	maxRetries  uint
)

func SetAccessToken(token string) {
	accessToken = token
}

func SetTimeout(_timeout time.Duration) {
	timeout = _timeout
}

func SetMaxRetries(_maxRetries uint) {
	maxRetries = _maxRetries
}

func GetAccessToken() string {
	return accessToken
}

func GetTimeout() time.Duration {
	return timeout
}

func GetMaxRetries() uint {
	return maxRetries
}
