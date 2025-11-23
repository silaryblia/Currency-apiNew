package repository

import (
	"Currency-apiNew/internal/domain"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// CurrencyRepoPostgreSQL реализует domain.CurrencyRepository для PostgreSQL
type CurrencyRepoPostgreSQL struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewCurrencyRepoPostgreSQL создает новый экземпляр PostgreSQL репозитория
func NewCurrencyRepoPostgreSQL(cfg domain.DatabaseConfig, logger *zap.Logger) (domain.CurrencyRepository, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping БД: %w", err)
	}

	// Устанавливаем максимальное количество соединений
	db.SetMaxOpenConns(cfg.MaxConns)

	logger.Info("Успешное подключение к PostgreSQL",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("dbname", cfg.Name))

	repo := &CurrencyRepoPostgreSQL{
		db:     db,
		logger: logger}

	// Инициализируем таблицу
	if err := repo.initTable(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации таблицы: %w", err)
	}

	// Загружаем дефолтные данные
	if err := repo.loadDefaultRates(cfg.DefaultRates); err != nil {
		return nil, fmt.Errorf("ошибка загрузки дефолтных курсов: %w", err)
	}

	return repo, nil
}

// initTable создает таблицу если она не существует
func (repo *CurrencyRepoPostgreSQL) initTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS currencies (
			code VARCHAR(3) PRIMARY KEY,
			rate DECIMAL(10,4) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_currencies_code ON currencies(code);
	`

	_, err := repo.db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	repo.logger.Info("Таблица currencies успешно инициализирована")
	return nil
}

// loadDefaultRates загружает дефолтные курсы валют
func (repo *CurrencyRepoPostgreSQL) loadDefaultRates(defaultRates map[string]float64) error {
	for code, rate := range defaultRates {
		// Используем UPSERT (INSERT ... ON CONFLICT)
		query := `
			INSERT INTO currencies (code, rate) 
			VALUES ($1, $2)
			ON CONFLICT (code) DO NOTHING
		`

		_, err := repo.db.Exec(query, code, rate)
		if err != nil {
			return fmt.Errorf("ошибка загрузки дефолтного курса %s: %w", code, err)
		}
	}

	repo.logger.Info("Дефолтные курсы валют загружены",
		zap.Int("count", len(defaultRates)))
	return nil
}

// GetAll возвращает все валюты из БД
func (repo *CurrencyRepoPostgreSQL) GetAll() (map[string]float64, error) {
	repo.logger.Debug("Получение всех валют из PostgreSQL")

	query := "SELECT code, rate FROM currencies"
	rows, err := repo.db.Query(query)
	if err != nil {
		repo.logger.Error("Ошибка выполнения запроса GetAll", zap.Error(err))
		return nil, domain.ErrDatabase
	}
	defer rows.Close()

	rates := make(map[string]float64)
	for rows.Next() {
		var code string
		var rate float64

		if err := rows.Scan(&code, &rate); err != nil {
			repo.logger.Error("Ошибка сканирования строки", zap.Error(err))
			return nil, domain.ErrDatabase
		}
		rates[code] = rate
	}

	if err := rows.Err(); err != nil {
		repo.logger.Error("Ошибка итерации по результатам", zap.Error(err))

		return nil, domain.ErrDatabase
	}

	repo.logger.Debug("Успешно получены валюты из БД",
		zap.Int("count", len(rates)))

	return rates, nil
}

// Get возвращает курс валюты по коду
func (repo *CurrencyRepoPostgreSQL) Get(code string) (float64, bool, error) {
	repo.logger.Debug("Получение валюты из PostgreSQL", zap.String("code", code))

	query := "SELECT rate FROM currencies WHERE code = $1"
	var rate float64

	err := repo.db.QueryRow(query, code).Scan(&rate)
	if err != nil {
		if err == sql.ErrNoRows {
			repo.logger.Debug("Валюта не найдена в БД", zap.String("code", code))

			return 0, false, nil
		}

		repo.logger.Error("Ошибка выполнения запроса Get",
			zap.String("code", code), zap.Error(err))

		return 0, false, domain.ErrDatabase
	}

	repo.logger.Debug("Валюта найдена в БД",
		zap.String("code", code), zap.Float64("rate", rate))
	return rate, true, nil
}

// AddOrUpdate добавляет или обновляет валюту в БД
func (repo *CurrencyRepoPostgreSQL) AddOrUpdate(code string, rate float64) error {
	repo.logger.Info("Добавление/обновление валюты в PostgreSQL",
		zap.String("code", code), zap.Float64("rate", rate))

	query := `
		INSERT INTO currencies (code, rate) 
		VALUES ($1, $2)
		ON CONFLICT (code) 
		DO UPDATE SET rate = $2, updated_at = CURRENT_TIMESTAMP
		RETURNING code
	`

	var resultCode string
	err := repo.db.QueryRow(query, code, rate).Scan(&resultCode)
	if err != nil {
		repo.logger.Error("Ошибка добавления/обновления валюты",
			zap.String("code", code), zap.Error(err))

		return domain.ErrDatabase
	}

	repo.logger.Info("Валюта успешно сохранена в БД",
		zap.String("code", code), zap.Float64("rate", rate))

	return nil
}

// Delete удаляет валюту из БД
func (repo *CurrencyRepoPostgreSQL) Delete(code string) (bool, error) {
	repo.logger.Info("Удаление валюты из PostgreSQL",
		zap.String("code", code))

	query := `DELETE FROM currencies WHERE code = $1`
	result, err := repo.db.Exec(query, code)
	if err != nil {
		repo.logger.Error("Ошибка удаления валюты",
			zap.String("code", code), zap.Error(err))

		return false, domain.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("Ошибка получения количества удаленных строк",
			zap.String("code", code), zap.Error(err))
		return false, domain.ErrDatabase
	}

	deleted := rowsAffected > 0
	if deleted {
		repo.logger.Info("Валюта успешно удалена из БД",
			zap.String("code", code))
	} else {
		repo.logger.Warn("Попытка удалить несуществующую валюту",
			zap.String("code", code))
	}

	return deleted, nil
}

// Exists проверяет существование валюты в БД
func (repo *CurrencyRepoPostgreSQL) Exists(code string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM currencies WHERE code = $1)"
	var exists bool

	err := repo.db.QueryRow(query, code).Scan(&exists)
	if err != nil {
		repo.logger.Error("Ошибка проверки существования валюты",
			zap.String("code", code), zap.Error(err))

		return false, domain.ErrDatabase
	}

	return exists, nil
}

// Close закрывает соединение с БД
func (repo *CurrencyRepoPostgreSQL) Close() error {
	if repo.db != nil {
		return repo.db.Close()
	}
	return nil
}
