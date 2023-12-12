package main

import (
	"encoding/json"
	"errors"
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
		url = "https://public.api.socio.events"
	} else {
		url = "https://public.sandbox-api.socio.events"
	}

	var uri = url + path

	requestUrl[GetAccessToken()] = uri

	return uri
}

func getUserAgent() string {
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

func fillErrorResponse(response *Response, errorResponse *ErrorResponse) error {
	if !(response.Status >= 400 && response.Status <= 500) {
		return nil
	}

	if json.Valid([]byte(response.Body)) {
		err := json.Unmarshal([]byte(response.Body), &errorResponse)
		if err != nil {
			return err
		}
		response.ErrorResponse = *errorResponse
	} else {
		return errors.New(fmt.Sprintf("The provided JSON is not valid. Provided body is: %s", response.Body))
	}

	return nil
}
