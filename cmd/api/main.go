package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/avirooppal/gosysutil/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
