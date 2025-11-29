package models

type BalanceAdjustment struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}
