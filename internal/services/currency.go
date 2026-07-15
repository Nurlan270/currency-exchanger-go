package services

import (
	"currency-exchanger/internal/models"
	"currency-exchanger/internal/stores"
)

type CurrencyService struct {
	store *stores.CurrencyStore
}

func NewCurrencyService(store *stores.CurrencyStore) *CurrencyService {
	return &CurrencyService{store: store}
}

func (s CurrencyService) GetAll() ([]models.Currency, error) {
	return s.store.GetAll()
}

func (s CurrencyService) GetByCode(code string) (models.Currency, error) {
	return s.store.GetByCode(code)
}

func (s CurrencyService) Create(name string, code string, sign string) (models.Currency, error) {
	return s.store.Create(name, code, sign)
}
