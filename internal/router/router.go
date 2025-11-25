package router

import (
	"Currency-apiNew/internal/domain"
	"Currency-apiNew/internal/handler"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Router настраивает и возвращает HTTP роутер
type Router struct {
	router *mux.Router
	logger *zap.Logger
}

// New создает новый экземпляр роутера
func New(logger *zap.Logger) *Router {
	return &Router{
		router: mux.NewRouter(),
		logger: logger,
	}
}

// Setup настраивает все маршруты и middleware
func (r *Router) Setup(
	currencyHandler domain.CurrencyHandler,
	healthHandler domain.HealthHandler,
) *mux.Router {

	r.logger.Info("Настройка маршрутов API")

	// Middleware (порядок важен!)
	r.router.Use(handler.RecoveryMiddleware(r.logger))
	r.router.Use(handler.LoggingMiddleware(r.logger))

	//  РЕГИСТРИРУЕМ МАРШРУТЫ В РОУТЕРЕ
	r.registerCurrencyRoutes(currencyHandler)
	r.registerHealthRoutes(healthHandler)

	r.logger.Info("Все маршруты успешно настроены")
	return r.router
}

// registerCurrencyRoutes регистрирует маршруты для валют
func (r *Router) registerCurrencyRoutes(currencyHandler domain.CurrencyHandler) {
	r.logger.Debug("Регистрация маршрутов для валют")

	api := r.router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/currencies", currencyHandler.GetAllCurrencies).Methods("GET")
	api.HandleFunc("/currencies", currencyHandler.CreateCurrency).Methods("POST")
	api.HandleFunc("/currencies/{code}", currencyHandler.GetCurrency).Methods("GET")
	api.HandleFunc("/currencies/{code}", currencyHandler.UpdateCurrency).Methods("PUT", "PATCH")
	api.HandleFunc("/currencies/{code}", currencyHandler.DeleteCurrency).Methods("DELETE")
}

// registerHealthRoutes регистрирует health-check маршруты
func (r *Router) registerHealthRoutes(healthHandler domain.HealthHandler) {
	r.logger.Debug("Регистрация health-check маршрутов")

	r.router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	r.router.HandleFunc("/", healthHandler.Help).Methods("GET")
}

// GetRouter возвращает базовый роутер
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
