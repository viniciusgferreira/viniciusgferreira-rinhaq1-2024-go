package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/db"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/handlers"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/repositories"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/services"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	pool, err := db.Connect()
	repo := repositories.New(pool)
	s := services.New(repo)
	h := handlers.NewClientHandler(s)

	mux := gin.New()
	mux.Use(gin.Recovery())
	//mux.Use(customLogger())
	gin.SetMode(gin.ReleaseMode)
	mux.Handle(http.MethodGet, "/clientes/:id/extrato", h.CreateStatement)
	mux.Handle(http.MethodPost, "/clientes/:id/transacoes", h.CreateTransaction)
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
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
}

func customLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		//// Read the request body
		//body, err := io.ReadAll(c.Request.Body)
		//if err != nil {
		//	log.Println("error reading body")
		//	// Handle error reading request body
		//	// You may choose to log the error or handle it in another way
		//}

		//log.Println(string(body))
		//defer c.Request.Body.Close() // Close the request body after reading

		// Restore the request body to its original state so it can be read again
		//c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Process request
		c.Next()

		// End timer
		end := time.Now()

		// Calculate latency
		latency := end.Sub(start)

		//fmt.Printf("[REQUEST] %s - %s %s\n", c.Request.Method, c.Request.URL, latency)

		// Log status code only if it's in the 4xx range
		status := c.Writer.Status()

		if status >= 400 && status < 600 {
			fmt.Printf("[CLIENT ERROR] %v - %s %s\n", status, c.Request.URL, latency)
		}
	}
}
