package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"Currency-apiNew/internal/domain"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type CurrencyHandler struct {
	service domain.CurrencyService
	logger  *zap.Logger
}

func NewCurrencyHandler(service domain.CurrencyService, logger *zap.Logger) domain.CurrencyHandler {
	return &CurrencyHandler{
		service: service,
		logger:  logger,
	}
}

// GetAllCurrencies обработчик для получения всех валют
func (h *CurrencyHandler) GetAllCurrencies(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Обработка HTTP GET /currencies")

	response := h.service.GetAllCurrencies()

	h.logger.Debug("Отправка списка валют клиенту",
		zap.Int("currencies_count", len(response.Currencies)))

	writeJSONResponse(w, http.StatusOK, response)
}

// GetCurrency обработчик для получения конкретной валюты
func (h *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := strings.ToLower(vars["code"])

	h.logger.Debug("Обработка HTTP GET /currencies/{code}",
		zap.String("code", code))

	currency, err := h.service.GetCurrency(code)
	if err != nil {
		h.logger.Warn("Валюта не найдена (HTTP)", zap.String("code", code))

		if err == domain.ErrCurrencyNotFound {
			writeErrorResponse(w, http.StatusNotFound, "Валюта не найдена")
		} else if err == domain.ErrDatabase {
			writeErrorResponse(w, http.StatusInternalServerError, "Ошибка базы данных")
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.logger.Debug("Отправка курса валюты клиенту",
		zap.String("code", code),
		zap.Float64("rate", currency.Rate))

	writeJSONResponse(w, http.StatusOK, currency)
}

// CreateCurrency обработчик для создания валюты
func (h *CurrencyHandler) CreateCurrency(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Обработка HTTP POST /currencies")

	var req domain.CreateCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Ошибка парсинга JSON запроса", zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный JSON в теле запроса")

		return
	}

	response, err := h.service.CreateCurrency(req)
	if err != nil {
		h.logger.Warn("Ошибка при создании валюты (HTTP)",
			zap.String("code", req.Code),
			zap.Error(err))

		if strings.Contains(err.Error(), "ошибка валидации") {
			// Извлекаем детальное сообщение об ошибке валидации
			validationMsg := extractValidationMessage(err.Error())
			writeErrorResponse(w, http.StatusBadRequest, validationMsg)
		} else if err == domain.ErrDuplicateCurrency {
			writeErrorResponse(w, http.StatusConflict, "Валюта уже существует")
		} else if err == domain.ErrDatabase {
			writeErrorResponse(w, http.StatusInternalServerError, "Ошибка базы данных")
		} else {
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.logger.Info("Валюта успешно создана (HTTP)",
		zap.String("code", req.Code),
		zap.Float64("rate", req.Rate))

	writeJSONResponse(w, http.StatusCreated, response)
}

// UpdateCurrency обработчик для обновления валюты
func (h *CurrencyHandler) UpdateCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := strings.ToLower(vars["code"])

	h.logger.Debug("Обработка HTTP PUT/PATCH /currencies/{code}",
		zap.String("code", code))

	var req domain.UpdateCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Ошибка парсинга JSON запроса",
			zap.String("code", code),
			zap.Error(err))
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный JSON в теле запроса")
		return
	}

	response, err := h.service.UpdateCurrency(code, req)
	if err != nil {
		h.logger.Warn("Ошибка при обновлении валюты (HTTP)",
			zap.String("code", code),
			zap.Error(err))

		if strings.Contains(err.Error(), "ошибка валидации") {
			// Извлекаем детальное сообщение об ошибке валидации
			validationMsg := extractValidationMessage(err.Error())
			writeErrorResponse(w, http.StatusBadRequest, validationMsg)
		} else if err == domain.ErrCurrencyNotFound {
			writeErrorResponse(w, http.StatusNotFound, "Валюта не найдена")
		} else if err == domain.ErrDatabase {
			writeErrorResponse(w, http.StatusInternalServerError, "Ошибка базы данных")
		} else {
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.logger.Info("Валюта успешно обновлена (HTTP)",
		zap.String("code", code),
		zap.Float64("rate", req.Rate))

	writeJSONResponse(w, http.StatusOK, response)
}

// DeleteCurrency обработчик для удаления валюты
func (h *CurrencyHandler) DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := strings.ToLower(vars["code"])

	h.logger.Debug("Обработка HTTP DELETE /currencies/{code}",
		zap.String("code", code))

	response, err := h.service.DeleteCurrency(code)
	if err != nil {
		h.logger.Warn("Ошибка при удалении валюты (HTTP)",
			zap.String("code", code),
			zap.Error(err))

		switch err {
		case domain.ErrCurrencyNotFound:
			writeErrorResponse(w, http.StatusNotFound, "Валюта не найдена")
		case domain.ErrDatabase:
			writeErrorResponse(w, http.StatusInternalServerError, "Ошибка базы данных")
		default:
			writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	h.logger.Info("Валюта успешно удалена (HTTP)", zap.String("code", code))
	writeJSONResponse(w, http.StatusOK, response)
}

// ДОБАВЛЯЕМ ВСПОМОГАТЕЛЬНУЮ ФУНКЦИЮ ДЛЯ ИЗВЛЕЧЕНИЯ СООБЩЕНИЯ ВАЛИДАЦИИ
func extractValidationMessage(err string) string {
	// Убираем префикс "ошибка валидации: " для более чистого сообщения
	if strings.HasPrefix(err, "ошибка валидации: ") {
		return strings.TrimPrefix(err, "ошибка валидации: ")
	}
	return err
}
