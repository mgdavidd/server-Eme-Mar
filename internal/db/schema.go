package db

import (
	"database/sql"
	"log"
)

// RunMigrations ejecuta todas las tablas usando la conexión pasada.
func RunMigrations(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS clientes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL,
			telefono TEXT,
			deuda REAL DEFAULT 0
		);`,

		`CREATE TABLE IF NOT EXISTS insumos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL,

			unidad_medida TEXT NOT NULL,          
			stock_actual REAL NOT NULL DEFAULT 0,
			minimo_sugerido REAL NOT NULL DEFAULT 0,

			precio_unitario REAL NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS movimientos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			descripcion TEXT NOT NULL,
			tipo TEXT NOT NULL,
			monto REAL NOT NULL,
			fecha TEXT NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS calculadora_costos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL,
			costo_insumos REAL NOT NULL,
			costo_total REAL NOT NULL
		);`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatal("Error creando tablas:", err)
		}
	}

	log.Println("Migraciones ejecutadas ✔")
}
