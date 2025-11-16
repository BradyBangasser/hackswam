package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	host := os.Getenv("DB_HOST") // should be "db" (Podman Compose service name)
	port := os.Getenv("DB_PORT") // default "5432"
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	if port == "" {
		port = "5432"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	// Retry until DB is ready
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
		}
		if err == nil {
			break
		}
		fmt.Println("Waiting for DB to be ready...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic(err)
	}
}
