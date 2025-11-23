package domain

// CurrencyService определяет интерфейс для бизнес-логики работы с валютами
type CurrencyService interface {
	// GetAllCurrencies возвращает все валюты
	GetAllCurrencies() CurrenciesListResponse

	// GetCurrency возвращает валюту по коду
	GetCurrency(code string) (*CurrencyResponse, error)

	// CreateCurrency создает новую валюту
	CreateCurrency(req CreateCurrencyRequest) (*SuccessResponse, error)

	// UpdateCurrency обновляет существующую валюту
	UpdateCurrency(code string, req UpdateCurrencyRequest) (*SuccessResponse, error)

	// DeleteCurrency удаляет валюту
	DeleteCurrency(code string) (*SuccessResponse, error)
}
