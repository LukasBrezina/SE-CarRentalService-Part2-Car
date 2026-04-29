package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DATABASE *sql.DB

func Connect() {
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		log.Println("DATABASE_URL is empty")
	} else {
		log.Println("DATABASE_URL is set")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB connection error:", err)
	}

	DATABASE = db
}
