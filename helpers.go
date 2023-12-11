package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	userAgent  string
	requestUrl = make(map[string]string)
)

func GetRequestUrl() string {
	var value = requestUrl[GetAccessToken()]
	if len(value) > 0 {
		return value
	}

	var path = "/graphql"
	var url string
	if strings.HasPrefix(GetAccessToken(), "sk_live") {
		url = "http://localhost:3000"
	} else {
		url = "https://public.sandbox-api.socio.events"
	}

	var uri = url + path

	requestUrl[GetAccessToken()] = uri

	return uri
}

func GetUserAgent() string {
	if len(userAgent) > 0 {
		return userAgent
	}

	hostname, _ := os.Hostname()
	userAgent = fmt.Sprintf(
		"Webex Go SDK(v%s) - OS(%s) - hostname(%s) - Go Version(%s)",
		VERSION,
		runtime.GOOS,
		hostname,
		runtime.Version())

	return userAgent
}
