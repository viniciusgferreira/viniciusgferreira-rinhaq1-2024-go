package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func Connect() *sql.DB {
	//connStr := "host=localhost port=5432 user=postgres password=postgres dbname=rinhadb sslmode=disable"
	connStr := "host=rinha-db port=5432 user=postgres password=postgres dbname=rinhadb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		err = db.Ping()
		if err == nil {
			log.Println("Database connected")
			break
		}
		log.Println("database connection error:", err.Error())
		log.Println("retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
	return db
}
