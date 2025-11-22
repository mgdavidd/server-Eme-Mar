package services

import (
	"database/sql"

	"github.com/mgdavidd/server-Eme-Mar/internal/models"
)

type ClientService struct {
	DB *sql.DB
}

func NewClientService(db *sql.DB) *ClientService {
	return &ClientService{DB: db}
}

func (s *ClientService) GetAll() ([]models.Client, error) {
	rows, err := s.DB.Query(`
        SELECT id, nombre, telefono, deuda 
        FROM clientes
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Client{}
	for rows.Next() {
		var c models.Client
		rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Debt)
		list = append(list, c)
	}

	return list, nil
}

func (s *ClientService) GetById(id int) (models.Client, error) {
	var c models.Client

	err := s.DB.QueryRow(`
        SELECT id, nombre, telefono, deuda 
        FROM clientes WHERE id = ?
    `, id).Scan(&c.ID, &c.Name, &c.Phone, &c.Debt)

	if err == sql.ErrNoRows {
		return models.Client{}, ErrNotFound
	}
	if err != nil {
		return models.Client{}, err
	}
	return c, nil
}

func (s *ClientService) Create(c *models.Client) error {
	stmt, err := s.DB.Prepare(`
        INSERT INTO clientes (nombre, telefono, deuda)
        VALUES (?, ?, ?)
    `)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(c.Name, c.Phone, c.Debt)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	c.ID = id
	return nil
}

func (s *ClientService) UpdateClient(c *models.Client) error {
	res, err := s.DB.Exec(`
        UPDATE clientes
        SET nombre = ?, telefono = ?, deuda = ?
        WHERE id = ?
    `, c.Name, c.Phone, c.Debt, c.ID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *ClientService) DeleteClient(id int) error {
	res, err := s.DB.Exec(`DELETE FROM clientes WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// clientes que deben
func (s *ClientService) GetIndebtedClient() ([]models.Client, error) {
	rows, err := s.DB.Query(`
        SELECT id, nombre, telefono, deuda 
        FROM clientes WHERE deuda > 0
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Client{}
	for rows.Next() {
		var c models.Client
		rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Debt)
		list = append(list, c)
	}

	return list, nil
}
