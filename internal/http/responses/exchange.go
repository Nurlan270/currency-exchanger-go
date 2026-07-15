package responses

import "currency-exchanger/internal/models"

type ExchangeResponse struct {
	BaseCurrency    models.Currency `json:"baseCurrency"`
	TargetCurrency  models.Currency `json:"targetCurrency"`
	Rate            float64         `json:"rate"`
	Amount          float64         `json:"amount"`
	ConvertedAmount float64         `json:"convertedAmount"`
}
