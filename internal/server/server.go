package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Currency-apiNew/internal/handler"
	"Currency-apiNew/internal/repository"
	"Currency-apiNew/internal/service"
	"Currency-apiNew/pkg/logger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	router *mux.Router
}

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
	}

	s.setupDependencies()

	s.server = &http.Server{
		Addr:         ":8082",
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) setupDependencies() {
	logger.Logger.Info("Настройка зависимостей сервера")

	// Middleware (порядок важен!)
	s.router.Use(handler.RecoveryMiddleware(logger.Logger))
	s.router.Use(handler.LoggingMiddleware(logger.Logger))

	// Инициализация слоев
	repo := repository.NewCurrencyRepoInMemory(logger.Logger)
	currencyService := service.NewCurrencyService(repo, logger.Logger)
	currencyHandler := handler.NewCurrencyHandler(currencyService, logger.Logger)
	healthHandler := handler.NewHealthHandler(logger.Logger) // ✅ Добавляем HealthHandler

	// Настройка маршрутов
	currencyHandler.RegisterRoutes(s.router)
	healthHandler.RegisterRoutes(s.router) // ✅ Регистрируем health routes

	logger.Logger.Info("Все зависимости настроены")
}

func (s *Server) Start() error {
	logger.Logger.Info("Запуск сервера", zap.String("address", ":8082"))

	serverErrors := make(chan error, 1)

	go func() {
		logger.Logger.Info("Сервер слушает на http://localhost:8082")
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("ошибка сервера: %w", err)

	case <-stop:
		logger.Logger.Info("Получен сигнал остановки")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.server.Close()
			return fmt.Errorf("ошибка graceful shutdown: %w", err)
		}

		logger.Logger.Info("Сервер корректно остановлен")
	}

	return nil
}
