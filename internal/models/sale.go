package models

type Sale struct {
	ClientId int64      `json:"client_id"` // id del cliente que compra
	Items    []SaleItem `json:"items"`     //[{"product_id": 1, "Quantity": 7},{....}]
	Total    float64    `json:"total"`     // precio por el que se compro todo
	Date     string     `json:"date"`      // cuando se hizo la venta
	IsCredit bool       `json:"is_credit"` // true = fiado
}

type SaleItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}
