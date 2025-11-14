package services

import (
	"database/sql"

	"github.com/mgdavidd/server-Eme-Mar/internal/models"
)

type InsumoService struct {
	DB *sql.DB
}

func NewInsumoService(db *sql.DB) *InsumoService {
	return &InsumoService{DB: db}
}

func (i *InsumoService) GetAll() ([]models.Insumo, error) {
	rows, err := i.DB.Query("SELECT * FROM insumos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	insumos := []models.Insumo{}

	for rows.Next() {
		var i models.Insumo
		rows.Scan(&i.ID, &i.Name, &i.Um, &i.Stock, &i.MinStock, &i.UnitPrice)
		insumos = append(insumos, i)
	}
	return insumos, nil
}
