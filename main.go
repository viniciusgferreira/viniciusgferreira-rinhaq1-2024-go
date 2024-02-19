package main

import (
	"database/sql"
	"fmt"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/db"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/handlers"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/repositories"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/services"
	"log"
	"net/http"
	"os"
)

func main() {
	pg := db.Connect()
	defer func(pg *sql.DB) {
		err := pg.Close()
		if err != nil {
			log.Println("error while closing db connection")
		}
	}(pg)

	repo := repositories.New(pg)
	s := services.New(repo)
	h := handlers.NewClientHandler(s)

	mux := http.NewServeMux()
	mux.Handle("/clientes/", h)
	port, ok := os.LookupEnv("APP_PORT")
	var addr string
	if !ok {
		addr = "localhost:8080"
	} else {
		addr = ":" + port
	}
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Printf("Starting server at %v\n", addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
}
