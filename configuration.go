package main

import (
	"time"
)

type Config struct {
	accessToken string
	timeout     time.Duration
	maxRetries  uint
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
	if timeout < time.Duration(1)*time.Second {
		panic("timeout must be greater than or equal to 1 second")
	}
	c.timeout = timeout
}

func (c *Config) SetMaxRetries(maxRetries uint) {
	c.maxRetries = maxRetries
}

func (c *Config) GetAccessToken() string {
	return c.accessToken
}

func (c *Config) GetTimeout() time.Duration {
	return c.timeout
}

func (c *Config) GetMaxRetries() uint {
	return c.maxRetries
}
