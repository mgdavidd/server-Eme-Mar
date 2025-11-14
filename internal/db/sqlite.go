package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite", "./db/local.db")
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Error enabling foreign keys:", err)
	}

	return db
}
