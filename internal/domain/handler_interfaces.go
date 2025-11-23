package domain

import (
	"net/http"

	"github.com/gorilla/mux"
)

// CurrencyHandler определяет интерфейс для HTTP обработчиков
type CurrencyHandler interface {
	// GetAllCurrencies обработчик для получения всех валют
	GetAllCurrencies(w http.ResponseWriter, r *http.Request)

	// GetCurrency обработчик для получения конкретной валюты
	GetCurrency(w http.ResponseWriter, r *http.Request)

	// CreateCurrency обработчик для создания валюты
	CreateCurrency(w http.ResponseWriter, r *http.Request)

	// UpdateCurrency обработчик для обновления валюты
	UpdateCurrency(w http.ResponseWriter, r *http.Request)

	// DeleteCurrency обработчик для удаления валюты
	DeleteCurrency(w http.ResponseWriter, r *http.Request)

	// RegisterRoutes регистрирует маршруты в роутере
	RegisterRoutes(router *mux.Router)
}

// HealthHandler определяет интерфейс для health-check обработчиков
type HealthHandler interface {
	// HealthCheck обработчик для проверки здоровья сервера
	HealthCheck(w http.ResponseWriter, r *http.Request)

	// Help обработчик для справки по API
	Help(w http.ResponseWriter, r *http.Request)

	// RegisterRoutes регистрирует маршруты в роутере
	RegisterRoutes(router *mux.Router)
}
