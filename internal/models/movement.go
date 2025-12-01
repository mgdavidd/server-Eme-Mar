package models

type Move struct {
	ID          int64   `json:"id"`
	Amount      float64 `json:"amount"` //cantidad
	Type        string  `json:"type"`
	Description string  `json:"descripcion"` //precio total por el surtido
	Date        string  `json:"date"`
	ClientID    *int64  `json:"client_id,omitempty"`
}
