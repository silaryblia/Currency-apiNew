package service

import (
	"fmt"
	"strings"

	"Currency-apiNew/internal/domain"

	"go.uber.org/zap"
)

type CurrencyService struct {
	repo   domain.CurrencyRepository
	logger *zap.Logger
}

func NewCurrencyService(repo domain.CurrencyRepository, logger *zap.Logger) *CurrencyService {
	return &CurrencyService{
		repo:   repo,
		logger: logger,
	}
}

func (s *CurrencyService) GetAllCurrencies() domain.CurrenciesListResponse {
	s.logger.Debug("Получение всех валют (сервис)")

	rates := s.repo.GetAll()

	currencies := make([]domain.CurrencyResponse, 0, len(rates))
	for code, rate := range rates {
		currencies = append(currencies, domain.CurrencyResponse{
			Code: strings.ToUpper(code),
			Rate: rate,
		})
	}

	response := domain.CurrenciesListResponse{
		Currencies: currencies,
		Count:      len(currencies),
		Total:      len(currencies),
	}

	s.logger.Debug("Успешно сформирован ответ сервиса",
		zap.Int("currencies_count", len(currencies)))

	return response
}

func (s *CurrencyService) GetCurrency(code string) (*domain.CurrencyResponse, error) {
	s.logger.Debug("Получение валюты (сервис)", zap.String("code", code))

	rate, ok := s.repo.Get(code)
	if !ok {
		s.logger.Warn("Валюта не найдена (сервис)", zap.String("code", code))
		return nil, domain.ErrCurrencyNotFound
	}

	s.logger.Debug("Валюта найдена (сервис)",
		zap.String("code", code),
		zap.Float64("rate", rate))

	return &domain.CurrencyResponse{
		Code: strings.ToUpper(code),
		Rate: rate,
	}, nil
}

func (s *CurrencyService) CreateCurrency(req domain.CreateCurrencyRequest) (*domain.SuccessResponse, error) {
	s.logger.Debug("Создание валюты (сервис)",
		zap.String("code", req.Code),
		zap.Float64("rate", req.Rate))

	if err := domain.Validate.Struct(req); err != nil {
		s.logger.Warn("Ошибка валидации при создании валюты",
			zap.String("code", req.Code),
			zap.Error(err))
		return nil, fmt.Errorf("ошибка валидации: %w", err)
	}

	s.repo.AddOrUpdate(req.Code, req.Rate)

	s.logger.Info("Валюта успешно создана (сервис)",
		zap.String("code", req.Code),
		zap.Float64("rate", req.Rate))

	return &domain.SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Валюта %s создана с курсом %.2f", strings.ToUpper(req.Code), req.Rate),
	}, nil
}

func (s *CurrencyService) UpdateCurrency(code string, req domain.UpdateCurrencyRequest) (*domain.SuccessResponse, error) {
	s.logger.Debug("Обновление валюты (сервис)",
		zap.String("code", code),
		zap.Float64("rate", req.Rate))

	if err := domain.Validate.Struct(req); err != nil {
		s.logger.Warn("Ошибка валидации при обновлении валюты",
			zap.String("code", code),
			zap.Error(err))
		return nil, fmt.Errorf("ошибка валидации: %w", err)
	}

	if _, exists := s.repo.Get(code); !exists {
		s.logger.Warn("Попытка обновить несуществующую валюту", zap.String("code", code))
		return nil, domain.ErrCurrencyNotFound
	}

	s.repo.AddOrUpdate(code, req.Rate)

	s.logger.Info("Валюта успешно обновлена (сервис)",
		zap.String("code", code),
		zap.Float64("rate", req.Rate))

	return &domain.SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Валюта %s обновлена с курсом %.2f", strings.ToUpper(code), req.Rate),
	}, nil
}

func (s *CurrencyService) DeleteCurrency(code string) (*domain.SuccessResponse, error) {
	s.logger.Debug("Удаление валюты (сервис)", zap.String("code", code))

	deleted := s.repo.Delete(code)
	if !deleted {
		s.logger.Warn("Попытка удалить несуществующую валюту", zap.String("code", code))
		return nil, domain.ErrCurrencyNotFound
	}

	s.logger.Info("Валюта успешно удалена (сервис)", zap.String("code", code))

	return &domain.SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Валюта %s удалена", strings.ToUpper(code)),
	}, nil
}
