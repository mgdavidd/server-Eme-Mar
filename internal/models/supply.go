package models

type Supply struct {
	IdInsumo    int64   `json:"id_insumo"`
	Amount      float64 `json:"amount"`       //cantidad
	TotalAmount float64 `json:"total_amount"` //precio total por el surtido
	Date        string  `json:"date"`
}
