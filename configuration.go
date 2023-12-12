package main

import "time"

type Config struct {
	accessToken string
	timeout     time.Duration
	maxRetries  int
}

func NewConfig() *Config {
	return &Config{
		timeout: time.Duration(30),
	}
}

func (c *Config) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *Config) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *Config) SetMaxRetries(maxRetries int) {
	c.maxRetries = maxRetries
}

func (c *Config) GetAccessToken() string {
	return c.accessToken
}

func (c *Config) GetTimeout() time.Duration {
	return c.timeout
}

func (c *Config) GetMaxRetries() int {
	return c.maxRetries
}
