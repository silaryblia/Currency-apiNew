package domain

// DatabaseConfig содержит настройки подключения к БД
// Вынесено в domain, чтобы репозиторий мог использовать

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxConns     int
	DefaultRates map[string]float64
}

// МОЖНО УДАЛИТЬ
