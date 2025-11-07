package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("🏦 FlowBank API starting...")

	// TODO: Initialize database connection
	// TODO: Setup routes
	// TODO: Implement handlers

	// 1. Get port from environment variable and fix binding syntax
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	// Prepend the colon for Go net/http binding
	listenAddr := ":" + port

	// 2. Setup Health Check Route
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy", "service":"flowbank-api"}`)
	})

	// 3. Start the HTTP Server
	log.Printf("API listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
