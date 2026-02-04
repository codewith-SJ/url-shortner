
package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB


func InitDB() {
	connStr := "postgres://admin:admin@db:5432/urls?sslmode=disable"

	var err error

	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				break
			}
		}

		log.Println("Waiting for DB to be ready...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("DB not reachable after retries:", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS short_urls (
		id SERIAL PRIMARY KEY,
		code TEXT UNIQUE,
		original_url TEXT
	);`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connected and table ready")
}

