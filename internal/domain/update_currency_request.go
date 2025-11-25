package domain

// UpdateCurrencyRequest представляет запрос на обновление валюты
type UpdateCurrencyRequest struct {
	Rate float64 `json:"rate" validate:"required,gt=0"`
}

// Validate валидирует запрос на обновление валюты
func (r *UpdateCurrencyRequest) Validate() error {
	return Validate.Struct(r)
}

// NewUpdateCurrencyRequest создает новый запрос на обновление валюты
func NewUpdateCurrencyRequest(rate float64) *UpdateCurrencyRequest {
	return &UpdateCurrencyRequest{
		Rate: rate,
	}
}
