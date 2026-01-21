package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/avirooppal/gosysutil/api"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	// Simple logger middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		mux.ServeHTTP(w, r)
	})

	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
