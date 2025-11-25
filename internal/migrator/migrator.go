package migrator

import (
	"Currency-apiNew/internal/config"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
)

// Migrator управляет миграциями базы данных
type Migrator struct {
	migrate *migrate.Migrate
	logger  *zap.Logger
}

// NewMigrator создает новый экземпляр мигратора
func NewMigrator(dbConfig config.DatabaseConfig, logger *zap.Logger) (*Migrator, error) {
	// Строка подключения к PostgreSQL
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSLMode)

	// Создаем мигратор
	m, err := migrate.New(
		"file://migrations", // путь к файлам миграций
		connStr,             // строка подключения к БД
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания мигратора: %w", err)
	}

	return &Migrator{
		migrate: m,
		logger:  logger,
	}, nil
}

// RunMigrations выполняет все pending миграции
func (m *Migrator) RunMigrations() error {
	m.logger.Info("Запуск миграций базы данных")

	err := m.migrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка выполнения миграций: %w", err)
	}

	if err == migrate.ErrNoChange {
		m.logger.Info("Нет новых миграций для применения")
	} else {
		m.logger.Info("Все миграции успешно применены")
	}

	return nil
}

// RollbackMigration откатывает последнюю миграцию
func (m *Migrator) RollbackMigration() error {
	m.logger.Info("Откат последней миграции")

	err := m.migrate.Steps(-1)
	if err != nil {
		return fmt.Errorf("ошибка отката миграции: %w", err)
	}

	m.logger.Info("Миграция успешно откатана")
	return nil
}

// GetMigrationVersion возвращает текущую версию миграций
func (m *Migrator) GetMigrationVersion() (uint, bool, error) {
	version, dirty, err := m.migrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("ошибка получения версии миграции: %w", err)
	}

	return version, dirty, nil
}

// Close закрывает соединение мигратора
func (m *Migrator) Close() error {
	if m.migrate != nil {
		sourceErr, databaseErr := m.migrate.Close()
		if sourceErr != nil {
			return fmt.Errorf("ошибка закрытия source: %w", sourceErr)
		}
		if databaseErr != nil {
			return fmt.Errorf("ошибка закрытия database: %w", databaseErr)
		}
	}
	return nil
}
