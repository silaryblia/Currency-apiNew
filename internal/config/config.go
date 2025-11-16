package config

import (
	"time"
)

type Config struct {
	ServerAddress   string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func New() *Config {
	return &Config{
		ServerAddress:   "localhost:8080",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	}
}
