package domain

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Currency represents a currency entity
type Currency struct {
	Code string  `json:"code" validate:"required,alpha,lowercase"`
	Rate float64 `json:"rate" validate:"required,gt=0"`
}

// CurrencyRepository defines the interface for currency storage
type CurrencyRepository interface {
	GetAll() map[string]float64
	Get(code string) (float64, bool)
	AddOrUpdate(code string, rate float64)
	Delete(code string) bool
}

// Request structures
type CreateCurrencyRequest struct {
	Code string  `json:"code" validate:"required,alpha,lowercase"`
	Rate float64 `json:"rate" validate:"required,gt=0"`
}

type UpdateCurrencyRequest struct {
	Rate float64 `json:"rate" validate:"required,gt=0"`
}

// Response structures
type CurrencyResponse struct {
	Code string  `json:"code"`
	Rate float64 `json:"rate"`
}

type CurrenciesListResponse struct {
	Currencies []CurrencyResponse `json:"currencies"`
	Count      int                `json:"count"`
	Total      int                `json:"total"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Validator instance
var Validate = validator.New()

// Domain errors
var (
	ErrCurrencyNotFound = fmt.Errorf("валюта не найдена")
	ErrInvalidRate      = fmt.Errorf("курс должен быть положительным")
)
