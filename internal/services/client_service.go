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
	rows, err := s.DB.Query("SELECT id, nombre, telefono, deuda FROM clientes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := []models.Client{}
	for rows.Next() {
		var c models.Client
		rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Debt)
		clients = append(clients, c)
	}

	return clients, nil
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
