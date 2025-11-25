package main

import (
	"Currency-apiNew/internal/config"
	"Currency-apiNew/internal/server"
	"Currency-apiNew/pkg/logger"

	"go.uber.org/zap"
)

// Структура для настроек логгера
type LoggerSettings struct {
	Level       string
	Encoding    string
	Development bool
}

func main() {
	// Сначала загружаем конфигурацию БЕЗ логгера
	cfgLoader := config.NewViperConfigLoader()
	cfg, err := cfgLoader.Load()
	if err != nil {
		// Используем стандартный лог для ошибок загрузки конфига
		panic("Не удалось загрузить конфиг: " + err.Error())
	}

	// Инициализируем логгер
	if err := logger.InitLogger(
		cfg.Logger.Level,
		cfg.Logger.Encoding,
		cfg.Logger.Development,
	); err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}
	defer logger.Logger.Sync()

	logger.Logger.Info("Запуск приложения Currency API",
		zap.String("server_address", cfg.Server.Address))

	// Создаем и запускаем сервер с конфигом
	srv := server.NewServer(cfg)

	if err := srv.Start(); err != nil {
		logger.Logger.Fatal("Ошибка при работе сервера", zap.Error(err))
	}

	logger.Logger.Info("Приложение завершено")
}
