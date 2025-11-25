package main

import (
	"Currency-apiNew/internal/config"
	"Currency-apiNew/internal/migrator"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func main() {
	// Инициализируем простой логгер
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	if len(os.Args) < 2 {
		fmt.Println("Использование:")
		fmt.Println("  migrate up     - применить миграции")
		fmt.Println("  migrate down   - откатить последнюю миграцию")
		fmt.Println("  migrate version - показать текущую версию")
		os.Exit(1)
	}

	// Загружаем конфигурацию
	cfgLoader := config.NewViperConfigLoader()
	cfg, err := cfgLoader.Load()
	if err != nil {
		logger.Fatal("Не удалось загрузить конфиг", zap.Error(err))
	}

	// Создаем мигратор
	mig, err := migrator.NewMigrator(cfg.Database, logger)
	if err != nil {
		logger.Fatal("Не удалось создать мигратор", zap.Error(err))
	}
	defer mig.Close()

	command := os.Args[1]
	switch command {
	case "up":
		if err := mig.RunMigrations(); err != nil {
			logger.Fatal("Ошибка применения миграций", zap.Error(err))
		}
		logger.Info("Миграции успешно применены")

	case "down":
		if err := mig.RollbackMigration(); err != nil {
			logger.Fatal("Ошибка отката миграции", zap.Error(err))
		}
		logger.Info("Миграция успешно откатана")

	case "version":
		version, dirty, err := mig.GetMigrationVersion()
		if err != nil {
			logger.Fatal("Ошибка получения версии", zap.Error(err))
		}
		logger.Info("Текущая версия миграций",
			zap.Uint("version", version),
			zap.Bool("dirty", dirty))

	default:
		fmt.Println("Неизвестная команда:", command)
		os.Exit(1)
	}
}
