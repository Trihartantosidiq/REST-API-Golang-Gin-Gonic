package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "user=postgres dbname=products sslmode=disable password=admin123"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}

	fmt.Println("Connected to database")
}
