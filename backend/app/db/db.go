package db

import (
	"database/sql"
	"log"
	"os"
)

var DB *sql.DB

func Init() {
	db, err := sql.Open("pgx", os.Getenv("DB"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	DB = db
}
