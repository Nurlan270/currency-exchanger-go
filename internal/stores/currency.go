package stores

import (
	"currency-exchanger/internal/db"
	"currency-exchanger/internal/models"
	"database/sql"
	"errors"
	"github.com/ncruces/go-sqlite3"
)

type CurrencyStore struct {
	db *sql.DB
}

func NewCurrencyStore(db *sql.DB) *CurrencyStore {
	return &CurrencyStore{db: db}
}

func (s *CurrencyStore) GetAll() ([]models.Currency, error) {
	rows, err := s.db.Query("SELECT id, full_name, code, sign FROM currencies")
	if err != nil {
		return nil, err
	}

	var currencies []models.Currency
	for rows.Next() {
		var currency models.Currency
		if err := rows.Scan(&currency.ID, &currency.Name, &currency.Code, &currency.Sign); err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func (s *CurrencyStore) GetByCode(code string) (models.Currency, error) {
	row := s.db.QueryRow("SELECT id, full_name, code, sign FROM currencies WHERE code = ?", code)

	var currency models.Currency
	if err := row.Scan(&currency.ID, &currency.Name, &currency.Code, &currency.Sign); err != nil {
		return models.Currency{}, err
	}

	return currency, nil
}

func (s *CurrencyStore) Create(name string, code string, sign string) (models.Currency, error) {
	result, err := s.db.Exec("INSERT INTO currencies (full_name, code, sign) VALUES (?, ?, ?)", name, code, sign)
	if err != nil {
		if errors.Is(err, sqlite3.CONSTRAINT_UNIQUE) {
			return models.Currency{}, db.ErrRowAlreadyExists
		}
		return models.Currency{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Currency{}, err
	}

	return models.Currency{
		ID:   uint(id),
		Name: name,
		Code: code,
		Sign: sign,
	}, nil
}
