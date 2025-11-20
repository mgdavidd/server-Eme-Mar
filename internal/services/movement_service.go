package services

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/mgdavidd/server-Eme-Mar/internal/models"
)

type MovementService struct {
	DB *sql.DB
}

func NewMoveService(db *sql.DB) *MovementService {
	return &MovementService{DB: db}
}

func (s *MovementService) GetAll() ([]models.Move, error) {
	rows, err := s.DB.Query(`
        SELECT id, descripcion, tipo, monto, fecha 
        FROM movimientos
    `)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	moves := []models.Move{}
	for rows.Next() {
		var i models.Move
		rows.Scan(&i.ID, &i.Description, &i.Type, &i.Amount, &i.Date)
		moves = append(moves, i)
	}

	return moves, nil
}

func (s *MovementService) GetBalance() (models.Account, error) {
	var b models.Account
	err := s.DB.QueryRow(`
    SELECT
        (SELECT saldo FROM caja WHERE id = 1) AS balance,
        (SELECT COALESCE(SUM(deuda), 0) FROM clientes) AS total_receivable
    `).Scan(&b.Balance, &b.AmountOwed)
	if err != nil {
		return models.Account{}, err
	}
	return b, nil
}

func (s *MovementService) Supply(supply models.Supply) error {

	if supply.Amount <= 0 {
		return ErrInvalidInput
	}

	supply.Date = time.Now().Format("2006-01-02 15:04")

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
	var unitPrice float64

	err = tx.QueryRow(`
        SELECT stock_actual, nombre, precio_unitario 
        FROM insumos 
        WHERE id = ?
    `, supply.IdInsumo).Scan(&actualStock, &nameInsumo, &unitPrice)

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

	supply.TotalAmount = unitPrice * supply.Amount

	description := "Surtido de insumo: " + strings.ToUpper(nameInsumo) + " X " + strconv.FormatFloat(supply.Amount, 'f', -1, 64)

	_, err = tx.Exec(`
        INSERT INTO movimientos (descripcion, tipo, monto, fecha)
        VALUES (?, 'egreso', ?, ?)
    `, description, supply.TotalAmount, supply.Date)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE caja
		SET saldo = saldo - ?
		WHERE id = 1
	`, supply.TotalAmount)
	if err != nil {
		return err
	}

	return nil
}

//func (s *MovementService) Sell()
