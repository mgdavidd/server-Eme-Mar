package models

type CreditSale struct {
	SaleId      int64      `json:"sale_id"`
	Items       []SaleItem `json:"items"` //[{"product_id": 1, "Quantity": 7},{....}]
	Total       float64    `json:"total"` // precio por el que se compro todo
	Description string     `json:"description"`
	TotalPaid   float64    `json:"total_paid"` //lo que lleva el cliente pagado
	Date        string     `json:"date"`       // cuando se hizo la venta
	ClientName  string     `json:"client_name"`
}

type Payments struct {
	ID     int64   `json:"id"`
	Date   string  `json:"date"` // cuando se hizo la venta
	Amount float64 `json:"amount"`
}
