package db

import (
	"SE-CarRentalService/services"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DATABASE *sql.DB
var GRCPConnection *services.ConverterClient

func Connect(grcpConnection *services.ConverterClient) {
	GRCPConnection = grcpConnection
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
