package main

import (
	"context"
	"flowbank/internal/kafka"
	"flowbank/internal/model"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// A list of pre-defined User IDs
var userIDs = []uuid.UUID{
	uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), // Alice
	uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"), // Bob
}

// Create a transaction event
func generateFakeTransaction() model.Transaction {
	// 1. Pick a random user to debit/credit
	userID := userIDs[rand.Intn(len(userIDs))]

	// 2. Determine transaction type (debit or credit)
	txType := "debit"
	if rand.Intn(2) == 0 {
		txType = "credit"
	}

	// 3. Generate a random amount (between $1.00 and $500.00)
	amount := rand.Int63n(49901) + 100

	// 4. Set status and context
	return model.Transaction{
		ID:           uuid.New(),
		UserID:       userID,
		Amount:       amount,
		Type:         txType,
		Description:  "Simulated Purchase/Deposit",
		MerchantName: fmt.Sprintf("Vendor %d", rand.Intn(100)),
		Status:       "PENDING", // Initial status before processing
		CreatedAt:    time.Now().UTC(),
	}
}

func main() {
	// Start HTTP server first before anything else
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

	// Randomly pick user
	rand.Seed(time.Now().UnixNano())

	// --- 1. Config setup ---
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")
	intervalStr := os.Getenv("PRODUCER_INTERVAL")

	if brokers == "" || topic == "" || intervalStr == "" {
		log.Fatalf("FATAL: Required environment variables (KAFKA_BROKERS, KAFKA_TOPIC, PRODUCER_INTERVAL) not set.")
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalf("FATAL: Failed to parse PRODUCER_INTERVAL '%s':%v", intervalStr, err)
	}

	log.Printf("🏦 FlowBank Producer starting with interval: %v", interval)
	log.Printf("Connecting to Kafka brokers: %s, Topic: %s", brokers, topic)

	// -- 2. Initialise Kafka Producer ---
	producer, err := kafka.NewProducer([]string{brokers}, topic)
	if err != nil {
		log.Fatalf("FATAL: Could not initialise Kafka producer: %v", err)
	}
	defer producer.Close()

	log.Println("✅ Kafka Producer successfully initialised.")

	// --- 3. Start Publishing Loop ---
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			// 1. Generate new transaction
			tx := generateFakeTransaction()

			// 2. Publish to Kafka
			err := producer.PublishTransaction(ctx, tx)
			if err != nil {
				log.Printf("ERROR: Failed to publish transaction %s: %v", tx.ID.String(), err)
				continue
			}

			// 3. Log success
			log.Printf("PUBLISHED: TxID %s | User %s | Type %s | Amount %d cents",
				tx.ID.String(), tx.UserID.String(), tx.Type, tx.Amount)
		}
	}
}
