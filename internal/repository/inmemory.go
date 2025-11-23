package repository

import (
	"Currency-apiNew/internal/domain"
	"sync"

	"go.uber.org/zap"
)

type CurrencyRepoInMemory struct {
	rates  map[string]float64
	mu     *sync.RWMutex
	logger *zap.Logger
}

func NewCurrencyRepoInMemory(defaultRates map[string]float64, logger *zap.Logger) domain.CurrencyRepository {
	logger.Debug("Создание нового in-memory репозитория валют",
		zap.Int("default_rates_count", len(defaultRates)))

	rates := make(map[string]float64)
	for k, v := range defaultRates {
		rates[k] = v
	}

	return &CurrencyRepoInMemory{
		rates:  rates,
		mu:     &sync.RWMutex{},
		logger: logger,
	}
}

func (repo *CurrencyRepoInMemory) GetAll() (map[string]float64, error) {
	repo.logger.Debug("Получение всех валют из репозитория")

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	// Копируем данные чтобы не менять оригинал
	result := make(map[string]float64)
	for k, v := range repo.rates {
		result[k] = v
	}

	repo.logger.Debug("Успешно получены валюты",
		zap.Int("count", len(result)))

	return result, nil
}

func (repo *CurrencyRepoInMemory) Get(code string) (float64, bool, error) {
	repo.logger.Debug("Получение валюты из репозитория",
		zap.String("code", code))

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	rate, exists := repo.rates[code]

	if exists {
		repo.logger.Debug("Валюта найдена",
			zap.String("code", code),
			zap.Float64("rate", rate))
	} else {
		repo.logger.Warn("Валюта не найдена в репозитории",
			zap.String("code", code))
	}

	return rate, exists, nil
}

func (repo *CurrencyRepoInMemory) AddOrUpdate(code string, rate float64) error {
	repo.logger.Info("Добавление/обновление валюты в репозитории",
		zap.String("code", code),
		zap.Float64("rate", rate))

	repo.mu.Lock()
	defer repo.mu.Unlock()

	oldRate, existed := repo.rates[code]
	if existed {
		repo.logger.Info("Обновление существующей валюты",
			zap.String("code", code),
			zap.Float64("old_rate", oldRate),
			zap.Float64("new_rate", rate))
	} else {
		repo.logger.Info("Добавление новой валюты",
			zap.String("code", code),
			zap.Float64("rate", rate))
	}

	repo.rates[code] = rate

	repo.logger.Debug("Валюта успешно сохранена",
		zap.String("code", code),
		zap.Float64("rate", rate))

	return nil
}

func (repo *CurrencyRepoInMemory) Delete(code string) (bool, error) {
	repo.logger.Info("Удаление валюты из репозитория",
		zap.String("code", code))

	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, exists := repo.rates[code]
	if exists {
		delete(repo.rates, code)
		repo.logger.Info("Валюта успешно удалена",
			zap.String("code", code))

		return true, nil
	}

	repo.logger.Warn("Попытка удалить несуществующую валюту",
		zap.String("code", code))

	return false, nil
}

func (repo *CurrencyRepoInMemory) Exists(code string) (bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	_, exists := repo.rates[code]
	return exists, nil
}
