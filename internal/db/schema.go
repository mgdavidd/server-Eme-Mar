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
			fecha TEXT NOT NULL,
			cliente_id INTEGER NULL
		);`,

		// PRODUCTOS
		`CREATE TABLE IF NOT EXISTS productos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL,
			costo_total REAL NOT NULL,
			precio REAL NOT NULL,
			foto BLOB NULL
		);`,

		// PRODUCTO - INSUMOS
		`CREATE TABLE IF NOT EXISTS producto_insumos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			producto_id INTEGER NOT NULL,
			insumo_id INTEGER NOT NULL,
			cantidad_insumo REAL NOT NULL,

			FOREIGN KEY (producto_id) REFERENCES productos(id) ON DELETE CASCADE,
			FOREIGN KEY (insumo_id) REFERENCES insumos(id) ON DELETE CASCADE
		);`,

		// CAJA
		`CREATE TABLE IF NOT EXISTS caja (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			saldo REAL NOT NULL
		);`,

		/* ───────────────────────────────────────────── */
		/*        NUEVAS TABLAS DE VENTAS A CRÉDITO       */
		/* ───────────────────────────────────────────── */

		// CREDIT SALES (ventas fiadas)
		`CREATE TABLE IF NOT EXISTS credit_sales (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			client_id INTEGER NOT NULL,
			total REAL NOT NULL,
			remaining_balance REAL NOT NULL,
			date TEXT NOT NULL,

			FOREIGN KEY (client_id) REFERENCES clientes(id)
		);`,

		// ITEMS DE LAS VENTAS FIADAS
		`CREATE TABLE IF NOT EXISTS credit_sale_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			credit_sale_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity REAL NOT NULL,

			FOREIGN KEY (credit_sale_id) REFERENCES credit_sales(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES productos(id)
		);`,

		// ABONOS A VENTAS FIADAS
		`CREATE TABLE IF NOT EXISTS credit_payments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			credit_sale_id INTEGER NOT NULL,
			amount REAL NOT NULL,
			date TEXT NOT NULL,

			FOREIGN KEY (credit_sale_id) REFERENCES credit_sales(id) ON DELETE CASCADE
		);`,

		`CREATE INDEX IF NOT EXISTS idx_prod_ins_producto ON producto_insumos(producto_id);`,
		`CREATE INDEX IF NOT EXISTS idx_prod_ins_insumo ON producto_insumos(insumo_id);`,

		`CREATE INDEX IF NOT EXISTS idx_credit_sales_client ON credit_sales(client_id);`,
		`CREATE INDEX IF NOT EXISTS idx_credit_sale_items_sale ON credit_sale_items(credit_sale_id);`,
		`CREATE INDEX IF NOT EXISTS idx_credit_payments_sale ON credit_payments(credit_sale_id);`,

		`CREATE TRIGGER IF NOT EXISTS recalc_after_insert_prod_ins
		AFTER INSERT ON producto_insumos
		BEGIN
			UPDATE productos SET costo_total = (
				SELECT COALESCE(SUM(i.precio_unitario * pi.cantidad_insumo), 0)
				FROM producto_insumos pi JOIN insumos i ON pi.insumo_id = i.id
				WHERE pi.producto_id = NEW.producto_id
			) WHERE id = NEW.producto_id;
		END;`,

		`CREATE TRIGGER IF NOT EXISTS recalc_after_update_prod_ins
		AFTER UPDATE ON producto_insumos
		BEGIN
			UPDATE productos SET costo_total = (
				SELECT COALESCE(SUM(i.precio_unitario * pi.cantidad_insumo), 0)
				FROM producto_insumos pi JOIN insumos i ON pi.insumo_id = i.id
				WHERE pi.producto_id = NEW.producto_id
			) WHERE id = NEW.producto_id;
		END;`,

		`CREATE TRIGGER IF NOT EXISTS recalc_after_delete_prod_ins
		AFTER DELETE ON producto_insumos
		BEGIN
			UPDATE productos SET costo_total = (
				SELECT COALESCE(SUM(i.precio_unitario * pi.cantidad_insumo), 0)
				FROM producto_insumos pi JOIN insumos i ON pi.insumo_id = i.id
				WHERE pi.producto_id = OLD.producto_id
			) WHERE id = OLD.producto_id;
		END;`,

		`CREATE TRIGGER IF NOT EXISTS recalc_after_update_insumo_precio
		AFTER UPDATE OF precio_unitario ON insumos
		BEGIN
			UPDATE productos SET costo_total = (
				SELECT COALESCE(SUM(i.precio_unitario * pi.cantidad_insumo), 0)
				FROM producto_insumos pi JOIN insumos i ON pi.insumo_id = i.id
				WHERE pi.producto_id = productos.id
			) WHERE id IN (SELECT producto_id FROM producto_insumos WHERE insumo_id = NEW.id);
		END;`,

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_producto_insumo_unique
		ON producto_insumos(producto_id, insumo_id);`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatal("Error creando tablas: ", err)
		}
	}

	log.Println("Migraciones ejecutadas ✔")
}
