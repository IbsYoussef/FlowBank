package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"flowbank/internal/db"
	internalkafka "flowbank/internal/kafka"
	"flowbank/internal/model"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func main() {
	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

		// Start minimal HTTP for Render service
	go func () {
		http.HandleFunc("/health", func (w http.ResponseWriter, r *http.Request)  {
			fmt.Fprintf(w, `{"status":"ok"}`)
		})
		http.ListenAndServe(":"+port, nil)
	}()

	log.Println("🏦 FlowBank Consumer starting...")
	ctx := context.Background()

	// 1. Construct the Database URL from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
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
	brokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")

	if brokers == "" || kafkaTopic == "" || kafkaGroupID == "" {
		log.Fatal("FATAL: Required Kafka environment variable not set.")
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// If SASL credentials are set, use SASL_SSL for Aiven
	saslUser := os.Getenv("KAFKA_SASL_USERNAME")
	saslPass := os.Getenv("KAFKA_SASL_PASSWORD")
	caCert := os.Getenv("KAFKA_CA_CERT")

	if saslUser != "" && saslPass != "" && caCert != "" {
		tlsConfig, err := internalkafka.NewTLSConfig(caCert)
		if err != nil {
			log.Fatalf("FATAL: Failed to create TLS config: %v", err)
		}
		dialer.TLS = tlsConfig
		dialer.SASLMechanism = plain.Mechanism{
			Username: saslUser,
			Password: saslPass,
		}
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokers},
		Topic:    kafkaTopic,
		GroupID:  kafkaGroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
		Dialer:   dialer,
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
