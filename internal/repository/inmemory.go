package repository

import (
	"sync"

	"Currency-apiNew/internal/domain"

	"go.uber.org/zap"
)

type CurrencyRepoInMemory struct {
	rates  map[string]float64
	mu     *sync.RWMutex
	logger *zap.Logger
}

func NewCurrencyRepoInMemory(logger *zap.Logger) domain.CurrencyRepository {
	logger.Debug("Создание нового репозитория валют")

	return &CurrencyRepoInMemory{
		rates: map[string]float64{
			"usd": 80.00,
			"eur": 85.00,
			"aed": 20.00,
		},
		mu:     &sync.RWMutex{},
		logger: logger,
	}
}

func (repo *CurrencyRepoInMemory) GetAll() map[string]float64 {
	repo.logger.Debug("Получение всех валют из репозитория")

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	result := make(map[string]float64)
	for k, v := range repo.rates {
		result[k] = v
	}

	repo.logger.Debug("Успешно получены валюты",
		zap.Int("count", len(result)))
	return result
}

func (repo *CurrencyRepoInMemory) Get(code string) (float64, bool) {
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

	return rate, exists
}

func (repo *CurrencyRepoInMemory) AddOrUpdate(code string, rate float64) {
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
}

func (repo *CurrencyRepoInMemory) Delete(code string) bool {
	repo.logger.Info("Удаление валюты из репозитория",
		zap.String("code", code))

	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, exists := repo.rates[code]
	if exists {
		delete(repo.rates, code)
		repo.logger.Info("Валюта успешно удалена",
			zap.String("code", code))
		return true
	}

	repo.logger.Warn("Попытка удалить несуществующую валюту",
		zap.String("code", code))
	return false
}
