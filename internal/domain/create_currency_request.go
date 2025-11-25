package domain

// CreateCurrencyRequest представляет запрос на создание валюты
type CreateCurrencyRequest struct {
	Code string  `json:"code" validate:"required,alpha,lowercase"`
	Rate float64 `json:"rate" validate:"required,gt=0"`
}

// Validate валидирует запрос на создание валюты
func (r *CreateCurrencyRequest) Validate() error {
	return Validate.Struct(r)
}

// NewCreateCurrencyRequest создает новый запрос на создание валюты
func NewCreateCurrencyRequest(code string, rate float64) *CreateCurrencyRequest {
	return &CreateCurrencyRequest{
		Code: code,
		Rate: rate,
	}
}
