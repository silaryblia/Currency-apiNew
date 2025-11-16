package handler

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"Currency-apiNew/internal/domain"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errorResp := domain.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}
	writeJSONResponse(w, statusCode, errorResp)
}

func RecoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Паника в обработчике",
						zap.Any("panic", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("remote_addr", r.RemoteAddr),
						zap.String("stack", string(debug.Stack())))

					writeErrorResponse(w, http.StatusInternalServerError,
						"Внутренняя ошибка сервера")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			logger.Info("Входящий HTTP запрос",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()))

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			logLevel := logger.Info
			if rw.statusCode >= 400 && rw.statusCode < 500 {
				logLevel = logger.Warn
			} else if rw.statusCode >= 500 {
				logLevel = logger.Error
			}

			logLevel("Исходящий HTTP ответ",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status_code", rw.statusCode),
				zap.Duration("duration", duration))
		})
	}
}
