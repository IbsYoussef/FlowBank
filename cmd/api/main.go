package main

import (
	"context"
	"encoding/json"
	"flowbank/internal/db"
	"flowbank/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	log.Println("🏦 FlowBank API starting...")

	ctx := context.Background()

	// 1. Initialise database connection
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	if dbUser == "" {
		log.Fatal("FATAL: Database connection environment variables not set.")
	}

	database, err := db.NewConnection(ctx, databaseURL)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to database: %v", err)
	}
	defer database.Close()

	// Initialise new service layer
	svc := service.NewService(database)

	// 2. Setup routes and 3. Implement handlers
	// 2. Setup routes and 3. Implement handlers
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy", "service":"flowbank-api"}`)
	})

	// Main API Endpoint: /api/v1/accounts/{userID}
	mux.HandleFunc("/api/v1/accounts/", handleGetAccountDetails(svc))

	port := "8080"
	log.Printf("API listening on %s", port)

	// --- CORRECTED SERVER STARTUP ---
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

// handleGetAccountDetails is the HTTP handler for fetching account balance and transactions.
func handleGetAccountDetails(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1. Extract the userID from the URL path.
		// The path is expected to be exactly "/api/v1/accounts/{userID}"
		// Since the route is registered as "/api/v1/accounts/", this path will contain the UUID after the final slash.
		userID := strings.TrimPrefix(r.URL.Path, "/api/v1/accounts/")

		// Quick validation check to ensure we got a valid UUID string
		if len(userID) != 36 {
			http.Error(w, "Invalid URL format or missing User ID. Expected: /api/v1/accounts/{userID}", http.StatusBadRequest)
			return
		}

		// 2. Call the service layer to get the data
		details, err := svc.GetAccountDetails(r.Context(), userID)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			log.Printf("ERROR: Failed to fetch account details for %s: %v", userID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// 3. Write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(details); err != nil {
			log.Printf("ERROR: Failed to write response JSON: %v", err)
		}
	}
}
