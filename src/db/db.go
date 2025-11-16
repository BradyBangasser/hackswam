package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	if port == "" {
		port = "5432"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname,
	)

	var err error
	for i := 0; i < 10; i++ { // retry 10 times
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
		}
		if err == nil {
			fmt.Println("Connected to DB!")
			return
		}
		fmt.Println("Waiting for DB...", err)
		time.Sleep(2 * time.Second)
	}

	panic(fmt.Sprintf("Cannot connect to DB: %v", err))
}
