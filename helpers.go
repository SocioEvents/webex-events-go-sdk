package main

import (
	"fmt"
	"os"
	"runtime"
)

var (
	userAgent string
)

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
