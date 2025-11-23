package models

type Client struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Phone string  `json:"phone"`
	Debt  float64 `json:"debt"`
}
