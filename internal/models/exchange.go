package models

type Exchange struct {
	ID             uint     `json:"id"`
	BaseCurrency   Currency `json:"baseCurrency"`
	TargetCurrency Currency `json:"targetCurrency"`
	Rate           float64  `json:"rate"`
}
