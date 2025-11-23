package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ConfigLoader загружает конфигурацию из различных источников
type ConfigLoader interface {
	Load() (*AppConfig, error)
}

type ViperConfigLoader struct {
	configPath string
	configName string
}

func NewViperConfigLoader(configPath, configName string) *ViperConfigLoader {
	return &ViperConfigLoader{
		configPath: configPath,
		configName: configName,
	}
}

func (l *ViperConfigLoader) Load() (*AppConfig, error) {
	// Устанавливаем значения по умолчанию
	viper.SetDefault("server.address", "localhost:8080")
	viper.SetDefault("server.read_timeout", 15*time.Second)
	viper.SetDefault("server.write_timeout", 15*time.Second)
	viper.SetDefault("server.idle_timeout", 60*time.Second)
	viper.SetDefault("server.shutdown_timeout", 30*time.Second)
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.encoding", "console")
	viper.SetDefault("logger.development", true)
	viper.SetDefault("currency.default_rates", map[string]float64{
		"usd": 80.00,
		"eur": 85.00,
		"aed": 20.00,
	})

	// Настройка Viper
	if l.configPath != "" {
		viper.AddConfigPath(l.configPath)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigName(l.configName)
	viper.SetConfigType("yaml")

	// Чтение переменных окружения
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CURRENCY_API")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Чтение конфиг файла
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("ошибка чтения конфиг файла: %w", err)
		}
		// Файл не найден, используем значения по умолчанию + env vars
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига: %w", err)
	}

	return &config, nil
}

// LoadFromEnv загружает конфигурацию только из переменных окружения
func LoadFromEnv() *AppConfig {
	config := NewDefaultConfig()

	if addr := os.Getenv("CURRENCY_API_SERVER_ADDRESS"); addr != "" {
		config.Server.Address = addr
	}

	if level := os.Getenv("CURRENCY_API_LOGGER_LEVEL"); level != "" {
		config.Logger.Level = level
	}

	// Парсим дефолтные курсы из env
	for code := range config.Currency.DefaultRates {
		envVar := fmt.Sprintf("CURRENCY_API_DEFAULT_RATE_%s", strings.ToUpper(code))
		if rateStr := os.Getenv(envVar); rateStr != "" {
			if rate, err := strconv.ParseFloat(rateStr, 64); err == nil {
				config.Currency.DefaultRates[code] = rate
			}
		}
	}

	return config
}
