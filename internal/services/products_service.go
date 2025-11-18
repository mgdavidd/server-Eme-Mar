package services

import (
	"database/sql"

	"github.com/mgdavidd/server-Eme-Mar/internal/models"
)

type ProductService struct {
	DB *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) Create(p *models.Product) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	costoTotal := 9000.0
	for _, ins := range p.Insumos {
		var precio float64
		err := tx.QueryRow(`
			SELECT precio_unitario
			FROM insumos 
			WHERE id = ?
		`, ins.InsumoID).Scan(&precio)

		if err != nil {
			tx.Rollback()
			return err
		}

		costoTotal += precio * ins.Quantity
	}

	res, err := tx.Exec(`
		INSERT INTO productos (nombre, costo_total, precio, foto)
		VALUES (?, ?, ?, ?)
	`, p.Name, costoTotal, p.Price, p.Foto)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, _ := res.LastInsertId()
	p.ID = id
	p.TotalCost = costoTotal

	for _, ins := range p.Insumos {
		_, err := tx.Exec(`
			INSERT INTO producto_insumos (producto_id, insumo_id, cantidad_insumo)
			VALUES (?, ?, ?)
		`, id, ins.InsumoID, ins.Quantity)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	rows, err := s.DB.Query(`
        SELECT id, nombre, costo_total, precio, foto
        FROM productos
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Product{}

	for rows.Next() {
		var p models.Product

		err := rows.Scan(&p.ID, &p.Name, &p.TotalCost, &p.Price, &p.Foto)
		if err != nil {
			return nil, err
		}

		// Obtener insumos del producto
		insRows, err := s.DB.Query(`
			SELECT insumo_id, cantidad_insumo
			FROM producto_insumos
			WHERE producto_id = ?
		`, p.ID)
		if err != nil {
			return nil, err
		}

		insumos := []models.ProductInsumo{}

		for insRows.Next() {
			var ins models.ProductInsumo
			err := insRows.Scan(&ins.InsumoID, &ins.Quantity)
			if err != nil {
				insRows.Close()
				return nil, err
			}
			insumos = append(insumos, ins)
		}

		insRows.Close()

		p.Insumos = insumos
		list = append(list, p)
	}

	return list, nil
}
