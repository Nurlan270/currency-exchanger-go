package stores

import (
	"currency-exchanger/internal/db"
	"currency-exchanger/internal/models"
	"database/sql"
	"errors"
	"github.com/ncruces/go-sqlite3"
)

type ExchangeStore struct {
	db *sql.DB
}

func NewExchangeStore(db *sql.DB) *ExchangeStore {
	return &ExchangeStore{db: db}
}

func (s *ExchangeStore) GetAll() ([]models.Exchange, error) {
	rows, err := s.db.Query(`
		SELECT
			er.id,
			er.rate,
		
			bc.id,
			bc.full_name,
			bc.code,
			bc.sign,
		
			tc.id,
			tc.full_name,
			tc.code,
			tc.sign
		FROM exchange_rates er
		JOIN currencies bc
			ON bc.id = er.base_currency_id
		JOIN currencies tc
			ON tc.id = er.target_currency_id;
	`)
	if err != nil {
		return nil, err
	}

	var exchanges []models.Exchange
	for rows.Next() {
		var exchange models.Exchange
		var baseCurrency models.Currency
		var targetCurrency models.Currency

		if err := rows.Scan(
			&exchange.ID, &exchange.Rate,
			&baseCurrency.ID, &baseCurrency.Name, &baseCurrency.Code, &baseCurrency.Sign,
			&targetCurrency.ID, &targetCurrency.Name, &targetCurrency.Code, &targetCurrency.Sign,
		); err != nil {
			return nil, err
		}

		exchange.BaseCurrency = baseCurrency
		exchange.TargetCurrency = targetCurrency
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}

func (s *ExchangeStore) GetByCodes(baseCurrencyCode string, targetCurrencyCode string) (models.Exchange, error) {
	rows, err := s.db.Query(`
		SELECT id, full_name, code, sign
		FROM currencies
		WHERE code IN (?, ?);
	`, baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		return models.Exchange{}, err
	}

	currencies := make(map[string]models.Currency, 2)
	for rows.Next() {
		var c models.Currency

		if err := rows.Scan(&c.ID, &c.Name, &c.Code, &c.Sign); err != nil {
			return models.Exchange{}, err
		}

		currencies[c.Code] = c
	}

	baseCurrency := currencies[baseCurrencyCode]
	targetCurrency := currencies[targetCurrencyCode]

	row := s.db.QueryRow(`
		SELECT id, rate
		FROM exchange_rates
		WHERE base_currency_id = ? AND target_currency_id = ?;
	`, baseCurrency.ID, targetCurrency.ID)

	var exchange models.Exchange
	if err := row.Scan(&exchange.ID, &exchange.Rate); err != nil {
		return models.Exchange{}, err
	}

	exchange.BaseCurrency = baseCurrency
	exchange.TargetCurrency = targetCurrency

	return exchange, nil
}

func (s *ExchangeStore) GetByTargetCodes(
	targetCodeOne string, targetCodeTwo string,
) (models.Currency, models.Currency, float64, error) {
	row := s.db.QueryRow(`
		SELECT
			(er_to.rate / er_from.rate) AS rate,
		
			c_from.id, c_from.full_name, c_from.code, c_from.sign,
		
			c_to.id, c_to.full_name, c_to.code,c_to.sign
		FROM exchange_rates er_from
			 JOIN exchange_rates er_to
				  ON er_from.base_currency_id = er_to.base_currency_id
			 JOIN currencies c_from
				  ON c_from.id = er_from.target_currency_id
			 JOIN currencies c_to
				  ON c_to.id = er_to.target_currency_id
		WHERE c_from.code = ?
		  AND c_to.code = ?;
	`, targetCodeOne, targetCodeTwo)

	var rate float64
	var c1 models.Currency
	var c2 models.Currency
	if err := row.Scan(
		&rate,
		&c1.ID, &c1.Name, &c1.Code, &c1.Sign,
		&c2.ID, &c2.Name, &c2.Code, &c2.Sign,
	); err != nil {
		return models.Currency{}, models.Currency{}, 0, err
	}

	return c1, c2, rate, nil
}

func (s *ExchangeStore) Create(baseCurrencyCode string, targetCurrencyCode string, rate float64) (models.Exchange, error) {
	baseCurrency, targetCurrency, err := s.getCurrenciesByCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		return models.Exchange{}, err
	}

	result, err := s.db.Exec(`
		INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate) VALUES (?, ?, ?);
	`, baseCurrency.ID, targetCurrency.ID, rate)
	if err != nil {
		if errors.Is(err, sqlite3.CONSTRAINT_UNIQUE) {
			return models.Exchange{}, db.ErrRowAlreadyExists
		}
		return models.Exchange{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Exchange{}, err
	}

	return models.Exchange{
		ID:             uint(id),
		Rate:           rate,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
	}, nil
}

func (s *ExchangeStore) Update(baseCurrencyCode string, targetCurrencyCode string, rate float64) (models.Exchange, error) {
	baseCurrency, targetCurrency, err := s.getCurrenciesByCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		return models.Exchange{}, err
	}

	result, err := s.db.Exec(`
		UPDATE exchange_rates
		SET rate = ?
		WHERE base_currency_id = ? AND target_currency_id = ?;
	`, rate, baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		return models.Exchange{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Exchange{}, err
	}

	return models.Exchange{
		ID:             uint(id),
		Rate:           rate,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
	}, nil
}

func (s *ExchangeStore) getCurrenciesByCodes(
	baseCurrencyCode string,
	targetCurrencyCode string,
) (models.Currency, models.Currency, error) {
	rows, err := s.db.Query(`
		SELECT id, full_name, code, sign
		FROM currencies
		WHERE code IN (?, ?);
	`, baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		return models.Currency{}, models.Currency{}, err
	}

	found := make(map[string]models.Currency, 2)

	for rows.Next() {
		var c models.Currency

		if err := rows.Scan(&c.ID, &c.Name, &c.Code, &c.Sign); err != nil {
			return models.Currency{}, models.Currency{}, err
		}

		found[c.Code] = c
	}

	baseCurrency, ok := found[baseCurrencyCode]
	if !ok {
		return models.Currency{}, models.Currency{}, db.ErrNotFound
	}

	targetCurrency, ok := found[targetCurrencyCode]
	if !ok {
		return models.Currency{}, models.Currency{}, db.ErrNotFound
	}

	return baseCurrency, targetCurrency, nil
}
