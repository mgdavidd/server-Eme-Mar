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

func (s *InsumoService) GetAll() ([]models.Insumo, error) {
	rows, err := s.DB.Query(`
        SELECT id, nombre, unidad_medida, stock_actual, minimo_sugerido, precio_unitario 
        FROM insumos
    `)
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

func (s *InsumoService) GetById(id int) (models.Insumo, error) {
	var i models.Insumo

	err := s.DB.QueryRow(`
        SELECT id, nombre, unidad_medida, stock_actual, minimo_sugerido, precio_unitario
        FROM insumos WHERE id = ?
    `, id).Scan(&i.ID, &i.Name, &i.Um, &i.Stock, &i.MinStock, &i.UnitPrice)

	if err == sql.ErrNoRows {
		return models.Insumo{}, ErrNotFound
	}
	if err != nil {
		return models.Insumo{}, err
	}
	return i, nil
}

func (s *InsumoService) Create(i *models.Insumo) error {
	stmt, err := s.DB.Prepare(`
        INSERT INTO insumos (nombre, unidad_medida, stock_actual, minimo_sugerido, precio_unitario)
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(i.Name, i.Um, i.Stock, i.MinStock, i.UnitPrice)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	i.ID = id
	return nil
}

func (s *InsumoService) Update(i *models.Insumo) error {
	res, err := s.DB.Exec(`
		UPDATE insumos
		SET nombre = ?, unidad_medida = ?, stock_actual = ?, minimo_sugerido = ?, precio_unitario = ?
		WHERE id = ?
	`, i.Name, i.Um, i.Stock, i.MinStock, i.UnitPrice, i.ID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
