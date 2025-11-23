package domain

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Currency - основная сущность (модель данных)
type Currency struct {
	Code string  `json:"code" validate:"required,alpha,lowercase"`
	Rate float64 `json:"rate" validate:"required,gt=0"`
}

//// CurrencyRepository определяет интерфейс для хранилища валют
//type CurrencyRepository interface {
//	GetAll() (map[string]float64, error)
//	Get(code string) (float64, bool, error)
//	AddOrUpdate(code string, rate float64) error
//	Delete(code string) (bool, error)
//	Exists(code string) (bool, error)
//}

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
	ErrCurrencyNotFound  = fmt.Errorf("валюта не найдена")
	ErrInvalidRate       = fmt.Errorf("курс должен быть положительным")
	ErrDatabase          = fmt.Errorf("ошибка базы данных")
	ErrDuplicateCurrency = fmt.Errorf("валюта уже существует")
)

/*
// CurrencyService интерфейс для бизнес-логики
type CurrencyService interface {
	GetAllCurrencies() (map[string]Currency, error)
	GetCurrency(code string) (*CurrencyResponse, error)
	CreateCurrency(req CreateCurrencyRequest) (*SuccessResponse, error)
	UpdateCurrency(code string, req UpdateCurrencyRequest) (*SuccessResponse, error)
	DeleteCurrency(code string) (*SuccessResponse, error)
}

// Структуры запросов и ответов
type CreateCurrencyRequest struct {
	Code string  `json:"code" validate:"required,alpha,lowercase"`
	Rate float64 `json:"rate" validate:"required,gt=0"`
}

type CurrencyResponse struct {
	Code string  `json:"code"`
	Rate float64 `json:"rate"`
}

// Validator instance
var Validate = validator.New()

/*

type UpdateCurrencyRequest struct {
	Rate float64 `json:"rate" validate:"required,gt=0"`
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



// Domain errors
var (
	ErrCurrencyNotFound = fmt.Errorf("валюта не найдена")
	ErrInvalidRate      = fmt.Errorf("курс должен быть положительным")
)
*/
