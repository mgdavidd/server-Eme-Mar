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
		WHERE date(fecha) >= date('now', '-30 days')
        ORDER BY fecha DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	moves := []models.Move{}
	for rows.Next() {
		var i models.Move
		if err := rows.Scan(&i.ID, &i.Description, &i.Type, &i.Amount, &i.Date); err != nil {
			return nil, err
		}
		moves = append(moves, i)
	}

	return moves, nil
}

func (s *MovementService) GetMovesByClient(clientID int) ([]models.Move, error) {
	rows, err := s.DB.Query(`
        SELECT id, descripcion, tipo, monto, fecha, cliente_id
        FROM movimientos
        WHERE cliente_id = ?
		  AND date(fecha) >= date('now', '-30 days')
        ORDER BY fecha DESC
    `, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	moves := []models.Move{}

	for rows.Next() {
		var m models.Move
		var clientIDNullable sql.NullInt64

		err := rows.Scan(
			&m.ID,
			&m.Description,
			&m.Type,
			&m.Amount,
			&m.Date,
			&clientIDNullable,
		)
		if err != nil {
			return nil, err
		}

		if clientIDNullable.Valid {
			clientIDValue := clientIDNullable.Int64
			m.ClientID = &clientIDValue
		} else {
			m.ClientID = nil
		}

		moves = append(moves, m)
	}

	return moves, nil
}

func (s *MovementService) GetRecent() ([]models.Move, error) {
	rows, err := s.DB.Query(`
		SELECT id, descripcion, tipo, monto, fecha
		FROM movimientos
		ORDER BY fecha DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	moves := []models.Move{}
	for rows.Next() {
		var i models.Move
		if err := rows.Scan(&i.ID, &i.Description, &i.Type, &i.Amount, &i.Date); err != nil {
			return nil, err
		}
		moves = append(moves, i)
	}

	return moves, nil
}

func (s *MovementService) GetBalance() (models.Account, error) {
	var b models.Account
	err := s.DB.QueryRow(`
        SELECT
            (SELECT saldo FROM caja WHERE id = 1),
            (SELECT COALESCE(SUM(deuda), 0) FROM clientes)
    `).Scan(&b.Balance, &b.AmountOwed)
	if err != nil {
		return models.Account{}, err
	}
	return b, nil
}

type Queryer interface {
	QueryRow(query string, args ...any) *sql.Row
}

func buildSaleDescription(items []models.SaleItem, db Queryer) (string, error) {

	var sb strings.Builder

	for _, item := range items {
		var name string
		err := db.QueryRow(`
			SELECT nombre 
			FROM productos 
			WHERE id = ?
		`, item.ProductID).Scan(&name)

		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		if err != nil {
			return "", err
		}

		sb.WriteString("- ")
		sb.WriteString(name)
		sb.WriteString(" x ")
		sb.WriteString(strconv.FormatInt(item.Quantity, 10))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (s *MovementService) Supply(supply models.Supply) (err error) {
	if supply.Amount <= 0 {
		return ErrInvalidInput
	}

	supply.Date = time.Now().Format("2006-01-02 15:04")

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			if cerr := tx.Commit(); cerr != nil {
				err = cerr
			}
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

	description := "Surtido de insumo: " +
		strings.ToTitle(nameInsumo) +
		" X " +
		strconv.FormatFloat(supply.Amount, 'f', -1, 64)

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

func (s *MovementService) Sell(sale models.Sale) (err error) {
	sale.Date = time.Now().Format("2006-01-02 15:04")

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			if cerr := tx.Commit(); cerr != nil {
				err = cerr
			}
		}
	}()

	var clientName string
	err = tx.QueryRow(`
        SELECT nombre FROM clientes WHERE id = ?
    `, sale.ClientId).Scan(&clientName)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	sale.Total = 0
	for _, item := range sale.Items {
		var price float64
		err = tx.QueryRow(`
			SELECT precio FROM productos WHERE id = ?
		`, item.ProductID).Scan(&price)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		sale.Total += price * float64(item.Quantity)
	}

	for _, item := range sale.Items {
		rows, err := tx.Query(`
            SELECT insumo_id, cantidad_insumo
            FROM producto_insumos
            WHERE producto_id = ?
        `, item.ProductID)
		if err != nil {
			return err
		}

		for rows.Next() {
			var insumoID int64
			var qtyPerProduct float64
			if err := rows.Scan(&insumoID, &qtyPerProduct); err != nil {
				rows.Close()
				return err
			}

			totalNeeded := qtyPerProduct * float64(item.Quantity)

			res, err := tx.Exec(`
                UPDATE insumos
                SET stock_actual = stock_actual - ?
                WHERE id = ? AND stock_actual >= ?
            `, totalNeeded, insumoID, totalNeeded)
			if err != nil {
				rows.Close()
				return err
			}
			ra, _ := res.RowsAffected()
			if ra == 0 {
				rows.Close()
				return ErrInvalidInput
			}
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			return err
		}
		rows.Close()
	}

	if sale.IsCredit {
		res, err := tx.Exec(`
			INSERT INTO credit_sales (client_id, total, remaining_balance, date)
			VALUES (?, ?, ?, ?)
		`, sale.ClientId, sale.Total, sale.Total, sale.Date)
		if err != nil {
			return err
		}

		creditID, _ := res.LastInsertId()

		for _, item := range sale.Items {
			_, err = tx.Exec(`
				INSERT INTO credit_sale_items (credit_sale_id, product_id, quantity)
				VALUES (?, ?, ?)
			`, creditID, item.ProductID, item.Quantity)
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(`
        UPDATE clientes
        SET deuda = deuda + ?
        WHERE id = ?`, sale.Total, sale.ClientId)
		if err != nil {
			return err
		}

		return nil
	}

	description, err := buildSaleDescription(sale.Items, tx)
	if err != nil {
		return err
	}

	description = strings.ToTitle(clientName) + ":\n" + description

	_, err = tx.Exec(`
    INSERT INTO movimientos (descripcion, tipo, monto, fecha, cliente_id)
    VALUES (?, 'ingreso', ?, ?, ?)
	`, description, sale.Total, sale.Date, sale.ClientId)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE caja
		SET saldo = saldo + ?
		WHERE id = 1
	`, sale.Total)
	if err != nil {
		return err
	}

	return nil
}

func (s *MovementService) PayCredit(creditSaleID int64, amount float64) (err error) {
	if amount <= 0 {
		return ErrInvalidInput
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			if cerr := tx.Commit(); cerr != nil {
				err = cerr
			}
		}
	}()

	var rem float64
	var clientID int64
	err = tx.QueryRow(`SELECT remaining_balance, client_id FROM credit_sales WHERE id = ?`, creditSaleID).Scan(&rem, &clientID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	if amount > rem {
		return ErrInvalidInput
	}

	res, err := tx.Exec(`UPDATE credit_sales SET remaining_balance = remaining_balance - ? WHERE id = ? AND remaining_balance >= ?`, amount, creditSaleID, amount)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return ErrInvalidInput
	}

	res, err = tx.Exec(`UPDATE clientes SET deuda = deuda - ? WHERE id = ? AND deuda >= ?`, amount, clientID, amount)
	if err != nil {
		return err
	}
	ra, _ = res.RowsAffected()
	if ra == 0 {
		return ErrInvalidInput
	}

	_, err = tx.Exec(`INSERT INTO credit_payments (credit_sale_id, amount, date) VALUES (?, ?, ?)`,
		creditSaleID, amount, time.Now().Format("2006-01-02 15:04"))
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
    INSERT INTO movimientos (descripcion, tipo, monto, fecha, cliente_id)
    VALUES (?, 'ingreso', ?, ?, ?)
	`, "Abono a crédito", amount, time.Now().Format("2006-01-02 15:04"), clientID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE caja SET saldo = saldo + ? WHERE id = 1`, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *MovementService) GetAllCreditSales() ([]models.CreditSale, error) {

	sales := []models.CreditSale{}

	rows, err := s.DB.Query(`
		SELECT 
			cs.id,
			cs.total,
			cs.remaining_balance,
			cs.date,
			c.nombre
		FROM credit_sales cs
		JOIN clientes c ON c.id = cs.client_id
		WHERE
			cs.remaining_balance > 0
			OR (
				cs.remaining_balance = 0
				AND date(cs.date) >= date('now', '-30 days')
			)
		ORDER BY cs.id
	`)
	if err != nil {
		return sales, err
	}
	defer rows.Close()

	for rows.Next() {
		var cs models.CreditSale
		var remain float64

		err := rows.Scan(
			&cs.SaleId,
			&cs.Total,
			&remain,
			&cs.Date,
			&cs.ClientName,
		)
		if err != nil {
			return sales, err
		}

		cs.TotalPaid = cs.Total - remain
		cs.Items = []models.SaleItem{}
		sales = append(sales, cs)
	}

	if err := rows.Err(); err != nil {
		return sales, err
	}

	for i := range sales {
		itemsRows, err := s.DB.Query(`
			SELECT product_id, quantity 
			FROM credit_sale_items
			WHERE credit_sale_id = ?
		`, sales[i].SaleId)
		if err != nil {
			return sales, err
		}

		for itemsRows.Next() {
			var item models.SaleItem
			err := itemsRows.Scan(&item.ProductID, &item.Quantity)
			if err != nil {
				itemsRows.Close()
				return sales, err
			}
			sales[i].Items = append(sales[i].Items, item)
		}
		itemsRows.Close()

		desc, err := buildSaleDescription(sales[i].Items, s.DB)
		if err != nil {
			return sales, err
		}
		sales[i].Description = desc
	}

	return sales, nil
}

func (s *MovementService) GetCreditSalesClients(client_id int) ([]models.CreditSale, error) {

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			if cerr := tx.Commit(); cerr != nil {
				err = cerr
			}
		}
	}()

	var tmpID int
	err = tx.QueryRow(`
		SELECT id FROM clientes WHERE id = ?
	`, client_id).Scan(&tmpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	rows, err := tx.Query(`
		SELECT cs.id, cs.total, cs.remaining_balance, cs.date, csi.product_id, csi.quantity
		FROM credit_sales cs
		JOIN credit_sale_items csi ON csi.credit_sale_id = cs.id
		WHERE cs.client_id = ?
		  AND (
			  cs.remaining_balance > 0
			  OR (
				  cs.remaining_balance = 0
				  AND date(cs.date) >= date('now', '-30 days')
			  )
		  )
		ORDER BY cs.id
	`, client_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	salesMap := make(map[int64]int)
	var creditSalesArr []models.CreditSale

	for rows.Next() {
		var saleID int64
		var total float64
		var remaining float64
		var dateStr string
		var productID int64
		var qty float64

		if err := rows.Scan(&saleID, &total, &remaining, &dateStr, &productID, &qty); err != nil {
			return nil, err
		}
		idx, exists := salesMap[saleID]
		if !exists {
			cs := models.CreditSale{
				SaleId:    saleID,
				Items:     []models.SaleItem{},
				Total:     total,
				Date:      dateStr,
				TotalPaid: total - remaining,
			}
			creditSalesArr = append(creditSalesArr, cs)
			idx = len(creditSalesArr) - 1
			salesMap[saleID] = idx
		}

		creditSalesArr[idx].Items = append(creditSalesArr[idx].Items, models.SaleItem{
			ProductID: productID,
			Quantity:  int64(qty),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range creditSalesArr {
		desc, err := buildSaleDescription(creditSalesArr[i].Items, tx)
		if err != nil {
			return nil, err
		}
		creditSalesArr[i].Description = desc
	}

	return creditSalesArr, nil
}

func (s *MovementService) GetCreditPayments(sale_id int) ([]models.Payments, error) {
	var tmpID int64
	err := s.DB.QueryRow(`
		SELECT id FROM credit_sales WHERE id = ?
	`, sale_id).Scan(&tmpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var paymentsArr []models.Payments

	rows, err := s.DB.Query(`
		SELECT id, date, amount FROM credit_payments WHERE credit_sale_id = ?
	`, sale_id)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var payment models.Payments
		if err := rows.Scan(&payment.ID, &payment.Date, &payment.Amount); err != nil {
			return nil, err
		}
		paymentsArr = append(paymentsArr, payment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return paymentsArr, nil
}

func (s *MovementService) AdjustBalance(req models.BalanceAdjustment) (err error) {
	if req.Amount < 0 {
		return ErrInvalidInput
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
		if cerr := tx.Commit(); cerr != nil {
			err = cerr
		}
	}()

	var currentBalance float64
	err = tx.QueryRow(`
		SELECT saldo FROM caja WHERE id = 1
	`).Scan(&currentBalance)
	if err != nil {
		return err
	}
	diff := req.Amount - currentBalance
	var movementType string
	if diff > 0 {
		movementType = "ingreso"
	}
	if diff < 0 {
		movementType = "egreso"
		diff = -diff
	}
	_, err = tx.Exec(`
		INSERT INTO movimientos (descripcion, tipo, monto, fecha)
		VALUES (?, ?, ?, ?)
	`, req.Description, movementType, diff, time.Now().Format("2006-01-02 15:04"))
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE caja
		SET saldo = ?
		WHERE id = 1
	`, req.Amount)
	if err != nil {
		return err
	}
	return nil
}
