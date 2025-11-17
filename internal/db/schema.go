package db

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) {
	queries := []string{
		// CLIENTES
		`CREATE TABLE IF NOT EXISTS clientes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL,
			telefono TEXT,
			deuda REAL DEFAULT 0
		);`,

		// INSUMOS
		`CREATE TABLE IF NOT EXISTS insumos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL,
			unidad_medida TEXT NOT NULL,
			stock_actual REAL NOT NULL DEFAULT 0,
			minimo_sugerido REAL NOT NULL DEFAULT 0,
			precio_unitario REAL NOT NULL
		);`,

		// MOVIMIENTOS
		`CREATE TABLE IF NOT EXISTS movimientos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			descripcion TEXT NOT NULL,
			tipo TEXT NOT NULL,
			monto REAL NOT NULL,
			fecha TEXT NOT NULL
		);`,

		// PRODUCTOS
		`CREATE TABLE IF NOT EXISTS productos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL,
			costo_total REAL NOT NULL,
			foto BLOB NULL
		);`,

		// PRODUCTO - INSUMOS (relación muchos-a-muchos)
		`CREATE TABLE IF NOT EXISTS producto_insumos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			producto_id INTEGER NOT NULL,
			insumo_id INTEGER NOT NULL,
			cantidad_insumo REAL NOT NULL,
			unidad_medida TEXT NOT NULL,

			FOREIGN KEY (producto_id) REFERENCES productos(id) ON DELETE CASCADE,
			FOREIGN KEY (insumo_id) REFERENCES insumos(id) ON DELETE CASCADE
		);`,

		// ÍNDICES para mejorar rendimiento de JOIN
		`CREATE INDEX IF NOT EXISTS idx_prod_ins_producto ON producto_insumos(producto_id);`,
		`CREATE INDEX IF NOT EXISTS idx_prod_ins_insumo ON producto_insumos(insumo_id);`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatal("Error creando tablas: ", err)
		}
	}

	log.Println("Migraciones ejecutadas ✔")
}
