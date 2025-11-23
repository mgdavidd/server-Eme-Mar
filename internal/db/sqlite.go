package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func ConnectDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠ No se pudo cargar .env, usando variables del sistema...")
	}
	url := os.Getenv("DATABASE_URL")
	token := os.Getenv("DATABASE_AUTH_TOKEN")

	var db *sql.DB

	if url != "" && token != "" {
		dsn := url + "?authToken=" + token

		db, err = sql.Open("libsql", dsn)
		if err == nil {
			log.Println("Conectado a Turso ✔")
			return db
		}

		log.Println("Error conectando a Turso, usando SQLite local:", err)
	}

	// Fallback SQLite local
	db, err = sql.Open("sqlite", "./db/local.db")
	if err != nil {
		log.Fatal("No se pudo conectar a SQLite local:", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Error enabling foreign keys:", err)
	}

	log.Println("Conectado a SQLite local ✔")
	return db
}
