package services

import "errors"

var (
	ErrNotFound     = errors.New("no encontrado")
	ErrInvalidInput = errors.New("datos inválidos")
)
