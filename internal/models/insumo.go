package models

type Insumo struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Um        string  `json:"um"`
	Stock     float64 `json:"stock"`
	MinStock  float64 `json:"min_stock"`
	UnitPrice float64 `json:"unit_price"`
}
