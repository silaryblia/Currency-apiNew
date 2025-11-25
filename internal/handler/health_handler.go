package handler

import (
	"Currency-apiNew/internal/domain"
	"net/http"
	"time"

	//"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HealthHandler struct {
	logger *zap.Logger
}

func NewHealthHandler(logger *zap.Logger) domain.HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Обработка запроса Health Check")

	healthResponse := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
		"service":   "Currency API",
	}

	writeJSONResponse(w, http.StatusOK, healthResponse)
}

func (h *HealthHandler) Help(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Обработка запроса справки")

	helpResponse := map[string]string{
		"message": "Currency API справка",
		"version": "1.0.0",
		"endpoints": `
GET  /api/v1/currencies          - получить все валюты
POST /api/v1/currencies          - создать новую валюту
GET  /api/v1/currencies/{code}   - получить конкретную валюту
PUT  /api/v1/currencies/{code}   - обновить валюту
PATCH /api/v1/currencies/{code}  - частично обновить валюту
DELETE /api/v1/currencies/{code} - удалить валюту
GET  /health                     - проверка здоровья сервера
GET  /                           - справка по API
`,
	}

	writeJSONResponse(w, http.StatusOK, helpResponse)
}
