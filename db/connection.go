package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func Connect() (*pgxpool.Pool, error) {
	//connStr := "host=localhost port=5432 user=postgres password=postgres dbname=rinhadb sslmode=disable"
	connStr := "host=rinha-db port=5432 user=postgres password=postgres dbname=rinhadb sslmode=disable"
	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = 75
	poolCfg.MinConns = 5
	db, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}
	for {
		err = db.Ping(context.Background())
		if err == nil {
			log.Println("Database connected")
			break
		}
		log.Println("database connection error:", err.Error())
		log.Println("retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
	return db, nil
}
