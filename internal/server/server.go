package server

import (
	"Currency-apiNew/internal/config"
	"Currency-apiNew/internal/domain"
	"Currency-apiNew/internal/handler"
	"Currency-apiNew/internal/repository"
	"Currency-apiNew/internal/router"
	"Currency-apiNew/internal/service"
	"Currency-apiNew/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	config *config.AppConfig
	repo   domain.CurrencyRepository
}

func NewServer(cfg *config.AppConfig) *Server {
	s := &Server{
		config: cfg,
	}

	// Настраиваем зависимости
	if err := s.setupDependencies(); err != nil {
		logger.Logger.Fatal("Ошибка настройки зависимостей", zap.Error(err))
	}

	return s
}

func (s *Server) setupDependencies() error {
	logger.Logger.Info("Настройка зависимостей сервера",
		zap.String("server_address", s.config.Server.Address))

	repo, err := repository.NewCurrencyRepoPostgreSQL(s.config.Database, logger.Logger)
	if err != nil {
		return fmt.Errorf("ошибка создания PostgreSQL репозитория: %w", err)
	}
	s.repo = repo

	//// Middleware (порядок важен!)
	//s.router.Use(handler.RecoveryMiddleware(logger.Logger))
	//s.router.Use(handler.LoggingMiddleware(logger.Logger))
	//
	//// СОЗДАЕМ POSTGRESQL РЕПОЗИТОРИЙ
	//dbConfig := domain.DatabaseConfig{
	//	Host:         s.config.Database.Host,
	//	Port:         s.config.Database.Port,
	//	User:         s.config.Database.User,
	//	Password:     s.config.Database.Password,
	//	Name:         s.config.Database.Name,
	//	SSLMode:      s.config.Database.SSLMode,
	//	MaxConns:     s.config.Database.MaxConns,
	//	DefaultRates: s.config.Currency.DefaultRates,
	//}

	// Инициализация слоев
	// Repository слой (данные)
	//repo, err := repository.NewCurrencyRepoPostgreSQL(dbConfig, logger.Logger)
	//if err != nil {
	//	return fmt.Errorf("Ошибка создания PostgreSQL репозитория: %w", err)
	//}
	//s.repo = repo

	// Service слой (бизнес-логика)
	currencyService := service.NewCurrencyService(repo, logger.Logger)

	// Handler слой (HTTP транспорт)
	currencyHandler := handler.NewCurrencyHandler(currencyService, logger.Logger)
	healthHandler := handler.NewHealthHandler(logger.Logger) // Добавляем HealthHandler

	// НАСТРАИВАЕМ РОУТЕР
	appRouter := router.New(logger.Logger)
	httpRouter := appRouter.Setup(currencyHandler, healthHandler)

	// создаем HTTP сервер
	s.server = &http.Server{
		Addr:         s.config.Server.Address,
		Handler:      httpRouter,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	logger.Logger.Info("Все зависимости настроены",
		zap.String("database", fmt.Sprintf("%s:%d", s.config.Database.Host, s.config.Database.Port)))

	return nil
}

func (s *Server) Start() error {
	logger.Logger.Info("Запуск сервера",
		zap.String("address", s.config.Server.Address))

	// ОБЕСПЕЧИВАЕМ GRACEFUL SHUTDOWN ДЛЯ БД
	defer func() {
		if postgresRepo, ok := s.repo.(*repository.CurrencyRepoPostgreSQL); ok {
			if err := postgresRepo.Close(); err != nil {
				logger.Logger.Error("Ошибка закрытия соединения с БД",
					zap.Error(err))
			} else {
				logger.Logger.Info("Соединение с БД закрыто")
			}
		}
	}()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Logger.Info("Сервер слушает на http://localhost:8082",
			zap.String("url", fmt.Sprintf("http://%s", s.config.Server.Address)))

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

		ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.ShutdownTimeout)

		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.server.Close()
			return fmt.Errorf("ошибка graceful shutdown: %w", err)
		}

		logger.Logger.Info("Сервер корректно остановлен")
	}

	return nil
}
