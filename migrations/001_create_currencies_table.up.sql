-- Создаем таблицу currencies
CREATE TABLE currencies (
                            code VARCHAR(3) PRIMARY KEY,
                            rate DECIMAL(10,4) NOT NULL,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс для быстрого поиска по коду валюты
CREATE INDEX idx_currencies_code ON currencies(code);

-- Вставляем начальные данные
INSERT INTO currencies (code, rate) VALUES
                                        ('usd', 80.00),
                                        ('eur', 85.00),
                                        ('aed', 20.00);