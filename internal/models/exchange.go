package models

type Exchange struct {
	ID             uint     `json:"id"`
	BaseCurrency   Currency `json:"base_currency"`
	TargetCurrency Currency `json:"target_currency"`
	Rate           float64  `json:"rate"`
}
