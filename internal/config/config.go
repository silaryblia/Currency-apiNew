package config

import (
	"time"
)

// AppConfig представляет основную конфигурацию приложения
type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Currency CurrencyConfig `mapstructure:"currency"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig содержит настройки HTTP сервера
type ServerConfig struct {
	Address         string        `mapstructure:"address"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// LoggerConfig содержит настройки логгера
type LoggerConfig struct {
	Level       string `mapstructure:"level"`    // debug, info, warn, error
	Encoding    string `mapstructure:"encoding"` // json, console
	Development bool   `mapstructure:"development"`
}

// CurrencyConfig содержит настройки для валют
type CurrencyConfig struct {
	DefaultRates map[string]float64 `mapstructure:"default_rates"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	MaxConns int    `mapstructure:"max_conns"`
}

// NewDefaultConfig возвращает конфигурацию по умолчанию
func NewDefaultConfig() *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Address:         "localhost:8082",
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Logger: LoggerConfig{
			Level:       "debug",
			Encoding:    "console",
			Development: true,
		},
		Currency: CurrencyConfig{
			DefaultRates: map[string]float64{
				"usd": 80.00,
				"eur": 85.00,
				"aed": 20.00,
			},
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "currrency_db",
			SSLMode:  "disable",
			MaxConns: 10,
		},
	}
}
