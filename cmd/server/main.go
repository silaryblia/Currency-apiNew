package main

import (
	"log"

	"Currency-apiNew/internal/server"
	"Currency-apiNew/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Инициализируем логгер
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer logger.Logger.Sync()

	logger.Logger.Info("Запуск приложения Currency API")

	// Создаем и запускаем сервер
	srv := server.NewServer()

	if err := srv.Start(); err != nil {
		logger.Logger.Fatal("Ошибка при работе сервера", zap.Error(err))
	}

	logger.Logger.Info("Приложение завершено")
}
