package models

type Product struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Price     float64         `json:"price"` // precio al que se vende
	Foto      []byte          `json:"foto,omitempty"`
	Insumos   []ProductInsumo `json:"insumos"` // solo id + cantidad
	TotalCost float64         `json:"costo_total"`
}

type ProductInsumo struct {
	InsumoID int64   `json:"id_insumo"`
	Quantity float64 `json:"quantity"` // cantidad requerida para fabricar el producto (misma unidad de medida del insumo de la tabla)
}
