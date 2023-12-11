package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"runtime"
	"time"
)

const VERSION = "0.1.0"

func Query(query string, operationName string, variables map[string]any, headers map[string]string) (Response, error) {

	if len(GetAccessToken()) < 1 {
		return Response{}, errors.New("Access Token is required")
	}

	var requestBody = make(map[string]any)
	requestBody["operationName"] = operationName
	requestBody["variables"] = variables
	requestBody["query"] = query
	var jsonBody, _ = json.Marshal(requestBody)

	req, err := http.NewRequest("POST", GetRequestUrl(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return Response{}, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+GetAccessToken())
	req.Header.Set("X-Sdk-Name", "Go SDK")
	req.Header.Set("X-Sdk-Version", VERSION)
	req.Header.Set("X-Sdk-Lang-Version", runtime.Version())
	req.Header.Set("User-Agent", GetUserAgent())

	var client = http.Client{Timeout: GetTimeout() * time.Second}
	var start = time.Now()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	var elapsed = time.Since(start)

	if err != nil {
		return Response{}, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return Response{}, err
	}

	var response = Response{
		status:         resp.StatusCode,
		headers:        resp.Header,
		body:           string(bodyBytes),
		requestHeaders: resp.Request.Header,
		requestBody:    string(jsonBody),
		url:            resp.Request.URL.String(),
		retryCount:     0,
		timeSpentInMs:  elapsed,
		rateLimiter:    RateLimiter{},
	}
	return response, nil
}
