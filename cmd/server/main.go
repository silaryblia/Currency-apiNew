package main

import (
	"Currency-apiNew/internal/config"
	"log"

	"Currency-apiNew/internal/server"
	"Currency-apiNew/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Загружаем конфигурацию
	cfgLoader := config.NewViperConfigLoader(".", "config")
	cfg, err := cfgLoader.Load()

	if err != nil {
		log.Printf("Не удалось загрузить конфиг: %v", err)
		// Используем конфиг по умолчанию
		cfg = config.NewDefaultConfig()
		log.Println("Используется конфигурация по умолчанию")
	}

	// Инициализируем логгер с настройками из конфига
	if err := logger.InitLogger(&cfg.Logger); err != nil {
		log.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer logger.Logger.Sync()

	logger.Logger.Info("Запуск приложения Currency API",
		zap.String("version", "1.0.0"),
		zap.String("server_address", cfg.Server.Address))

	// Создаем и запускаем сервер
	srv := server.NewServer(cfg)

	if err := srv.Start(); err != nil {

		logger.Logger.Fatal("Ошибка при работе сервера", zap.Error(err))
	}

	logger.Logger.Info("Приложение завершено")
}
