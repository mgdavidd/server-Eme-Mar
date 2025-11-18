package services

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mgdavidd/server-Eme-Mar/internal/models"
)

type MovementService struct {
	DB *sql.DB
}

func NewMoveService(db *sql.DB) *MovementService {
	return &MovementService{DB: db}
}

func (s *MovementService) Supply(supply models.Supply) error {

	if supply.Amount <= 0 || supply.TotalAmount <= 0 {
		return ErrInvalidInput
	}

	// Transacción
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var actualStock float64
	var nameInsumo string

	err = tx.QueryRow(`
        SELECT stock_actual, nombre 
        FROM insumos 
        WHERE id = ?
    `, supply.IdInsumo).Scan(&actualStock, &nameInsumo)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	newStock := actualStock + supply.Amount

	_, err = tx.Exec(`
        UPDATE insumos
        SET stock_actual = ?
        WHERE id = ?
    `, newStock, supply.IdInsumo)
	if err != nil {
		return err
	}

	description := "Surtido de insumo: " + strings.ToUpper(nameInsumo)

	_, err = tx.Exec(`
        INSERT INTO movimientos (descripcion, tipo, monto, fecha)
        VALUES (?, 'egreso', ?, ?)
    `, description, supply.TotalAmount, supply.Date)
	if err != nil {
		return err
	}

	return nil
}

//func (s *MovementService) Sell()
