package services

import (
	"currency-exchanger/internal/helpers"
	"currency-exchanger/internal/models"
	"currency-exchanger/internal/stores"
	"database/sql"
	"errors"
)

type ExchangeService struct {
	store *stores.ExchangeStore
}

func NewExchangeService(store *stores.ExchangeStore) *ExchangeService {
	return &ExchangeService{store: store}
}

func (s ExchangeService) GetAll() ([]models.Exchange, error) {
	return s.store.GetAll()
}

func (s ExchangeService) GetByCodes(baseCurrencyCode string, targetCurrencyCode string) (models.Exchange, error) {
	exchange, err := s.store.GetByCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		return models.Exchange{}, err
	}

	exchange.Rate = helpers.Round(exchange.Rate)

	return exchange, nil
}

func (s ExchangeService) GetByCodesWithRevert(
	baseCurrencyCode string, targetCurrencyCode string,
) (models.Exchange, error) {
	// Scenario 1
	exchange, err := s.store.GetByCodes(baseCurrencyCode, targetCurrencyCode)
	if err == nil {
		return exchange, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return models.Exchange{}, err
	}

	// Revert Rate (Scenario 2)
	exchange, err = s.store.GetByCodes(targetCurrencyCode, baseCurrencyCode)
	if err != nil {
		return models.Exchange{}, err
	}

	tmp := exchange.BaseCurrency // To not lose the reference to the original BaseCurrency
	exchange.BaseCurrency = exchange.TargetCurrency
	exchange.TargetCurrency = tmp
	exchange.Rate = helpers.Round(1 / exchange.Rate)

	return exchange, nil
}

func (s ExchangeService) GetByForeignRate(
	baseCurrencyCode string, targetCurrencyCode string,
) (models.Exchange, error) {
	c1, c2, rate, err := s.store.GetByTargetCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		return models.Exchange{}, err
	}

	exchange := models.Exchange{
		BaseCurrency:   c1,
		TargetCurrency: c2,
		Rate:           helpers.Round(rate),
	}

	return exchange, nil
}

func (s ExchangeService) Create(baseCurrencyCode string, targetCurrencyCode string, rate float64) (models.Exchange, error) {
	return s.store.Create(baseCurrencyCode, targetCurrencyCode, helpers.Round(rate))
}

func (s ExchangeService) Update(baseCurrencyCode string, targetCurrencyCode string, rate float64) (models.Exchange, error) {
	return s.store.Update(baseCurrencyCode, targetCurrencyCode, helpers.Round(rate))
}

func (s ExchangeService) ConvertAmount(amount float64, rate float64) float64 {
	return helpers.Round(amount * rate)
}
