package services

import (
	"database/sql"
	"errors"

	"github.com/mgdavidd/server-Eme-Mar/internal/models"
)

type ProductService struct {
	DB *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{DB: db}
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

func (s *ProductService) GetById(id int) (models.Product, error) {
	var p models.Product

	err := s.DB.QueryRow(`
        SELECT id, nombre, costo_total, precio, foto 
        FROM productos WHERE id = ?
    `, id).Scan(&p.ID, &p.Name, &p.TotalCost, &p.Price, &p.Foto)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Product{}, ErrNotFound
	}
	if err != nil {
		return models.Product{}, err
	}

	insRows, err := s.DB.Query(`
		SELECT insumo_id, cantidad_insumo
		FROM producto_insumos
		WHERE producto_id = ?
		`, p.ID)
	if err != nil {
		return models.Product{}, err
	}

	insumos := []models.ProductInsumo{}

	for insRows.Next() {
		var ins models.ProductInsumo
		err := insRows.Scan(&ins.InsumoID, &ins.Quantity)
		if err != nil {
			insRows.Close()
			return models.Product{}, err
		}
		insumos = append(insumos, ins)
	}
	insRows.Close()

	p.Insumos = insumos

	return p, nil
}

func (s *ProductService) Create(p *models.Product) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	costoTotal := 0.0
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

func (s *ProductService) Update(p models.ProductSimple) error {
	res, err := s.DB.Exec(`
		UPDATE productos
		SET nombre = ?, precio = ?, foto = ?
		WHERE id = ?
	`, p.Name, p.Price, p.Foto, p.ID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *ProductService) Delete(id int) error {
	res, err := s.DB.Exec(`DELETE FROM productos WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *ProductService) UpdateInsumoQuantity(productID, insumoID int64, quantity float64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(`
        UPDATE producto_insumos
        SET cantidad_insumo = ?
        WHERE producto_id = ? AND insumo_id = ?
    `, quantity, productID, insumoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		tx.Rollback()
		return ErrNotFound
	}

	// No recalculamos aquí: el trigger en la DB actualizará productos.costo_total
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// RemoveInsumoFromProduct deletes the relation producto_insumos and recalculates costo_total.
func (s *ProductService) RemoveInsumoFromProduct(productID, insumoID int64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(`
        DELETE FROM producto_insumos WHERE producto_id = ? AND insumo_id = ?
    `, productID, insumoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		tx.Rollback()
		return ErrNotFound
	}

	// El trigger AFTER DELETE en producto_insumos se encargará de recalcular costo_total
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *ProductService) UpdateOrCreateInsumo(productID, insumoID int64, quantity float64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// Verify product exists
	var tmp int64
	err = tx.QueryRow("SELECT id FROM productos WHERE id = ?", productID).Scan(&tmp)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return ErrNotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	// Verify insumo exists
	err = tx.QueryRow("SELECT id FROM insumos WHERE id = ?", insumoID).Scan(&tmp)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return ErrNotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	// Upsert (requiere índice único en producto_id, insumo_id)
	_, err = tx.Exec(`
        INSERT INTO producto_insumos (producto_id, insumo_id, cantidad_insumo)
        VALUES (?, ?, ?)
        ON CONFLICT(producto_id, insumo_id) DO UPDATE
        SET cantidad_insumo = excluded.cantidad_insumo
    `, productID, insumoID, quantity)
	if err != nil {
		tx.Rollback()
		return err
	}

	// No recalculamos ni actualizamos productos.costo_total aquí: el trigger lo hará automáticamente.
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
