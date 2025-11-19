package models

type ProductSimple struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"` // precio al que se vende
	Foto  []byte  `json:"foto,omitempty"`
}
