package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("🏦 FlowBank API starting...")

	// TODO: Initialize database connection
	// TODO: Setup routes
	// TODO: Implement handlers

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy", "service":"flowbank-api"}`)
	})

	port := "8080"
	log.Printf("API listening on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
