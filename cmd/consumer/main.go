package main

import (
	"context"
	"encoding/json"
	"flowbank/internal/db"
	"flowbank/internal/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// TODO: Connect to database
// TODO: Validate and persist to DB

const (
	topic         = "transactions"
	brokerAddress = "redpanda:9092" // Hostname of RedPanda container
	groupID       = "flowbank-consumer-group"
)

func main() {
	log.Println("🏦 FlowBank Consumer starting...")
	ctx := context.Background()

	// 1. Construct the Database URL from environment variables
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

	// 2. Connect to the database
	database, err := db.NewConnection(ctx, databaseURL)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to database: %v", err)
	}
	defer database.Close()

	// 3. Initialise Kafka consumer
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})

	// 4. Start consuming transactions
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Printf("ERROR fetching message: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("➡️ RECEIVED: Topic=%s, Partition=%d, Offset=%d, Key=%s", m.Topic, m.Partition, m.Offset, string(m.Key))

		var tx model.Transaction
		if err := json.Unmarshal(m.Value, &tx); err != nil {
			log.Printf("⛔ ERROR unmarshalling message (Offset %d): %v", m.Offset, err)
			r.CommitMessages(ctx, m)
			continue
		}

		// Proces the transaction (atomic DB update)
		if err := database.UpdateUserBalanceAndSaveTransaction(ctx, &tx); err != nil {
			log.Printf("❌ FAILED TO PROCESS (Offset %d, TxID %s): %v", m.Offset, tx.ID, err)
			r.CommitMessages(ctx, m)
			continue
		}

		// Commit the message manually only after successfully processing
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatalf("FATAL: Failed to commit message (Offset %d). Restarting consumer: %v", m.Offset, err)
		}

		log.Printf("✅ SUCESS: TxID %s processed and marked as %s. User %s balance updated.", tx.ID, tx.Status, tx.UserID)
	}
}
