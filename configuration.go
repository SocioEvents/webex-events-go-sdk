package main

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	accessToken string
	timeout     time.Duration
	maxRetries  uint
	logger      *slog.Logger
}

func NewConfig() *Config {
	var config = &Config{
		timeout:    time.Duration(30),
		maxRetries: 5,
	}

	var handlerOptions = &slog.HandlerOptions{Level: slog.LevelWarn}
	config.logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	return config
}

func (c *Config) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *Config) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *Config) SetMaxRetries(maxRetries uint) {
	c.maxRetries = maxRetries
}

func (c *Config) SetLoggerHandler(handler slog.Handler) {
	if handler == nil {
		var handlerOptions = &slog.HandlerOptions{Level: slog.LevelWarn}
		c.logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	} else {
		c.logger = slog.New(handler)
	}
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
