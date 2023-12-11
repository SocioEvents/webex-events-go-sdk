package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const UuidRegexPattern = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"

func GetRequestUrl() string {
	var path = "/graphql"
	var url string
	if strings.HasPrefix(GetAccessToken(), "sk_live") {
		url = "http://localhost:3000"
	} else {
		url = "https://public.sandbox-api.socio.events"
	}

	return url + path
}

func GetUserAgent() string {

	hostname, _ := os.Hostname()
	return fmt.Sprintf(
		"Webex Go SDK(v%s) - OS(%s) - hostname(%s) - Go Version(%s)",
		VERSION,
		runtime.GOOS,
		hostname,
		runtime.Version())
}
