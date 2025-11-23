package domain

// CurrencyRepository определяет интерфейс для работы с хранилищем валют
type CurrencyRepository interface {
	// GetAll возвращает все валюты из хранилища
	GetAll() (map[string]float64, error)

	// Get возвращает курс валюты по коду
	Get(code string) (float64, bool, error)

	// AddOrUpdate добавляет новую валюту или обновляет существующую
	AddOrUpdate(code string, rate float64) error

	// Delete удаляет валюту по коду
	Delete(code string) (bool, error)

	// Exists проверяет существование валюты
	Exists(code string) (bool, error)
}
