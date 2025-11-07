package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	log.Println("🏦 FlowBank Consumer starting...")

	// TODO: Connect to database
	// TODO: Validate and persist to DB

	// 1. Read connection details from environment variables
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	if brokers == "" || topic == "" || groupID == "" {
		log.Fatal("Missing required environment variables (KAFKA_BROKERS, KAFKA_TOPIC, KAFKA_GROUP_ID)")
	}

	// 2. Initialise Kafka consumer (Reader)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokers, ","),
		Topic:   topic,
		GroupID: groupID,
		MaxWait: 1 * time.Second, // Time to wait for new messages
	})
	defer r.Close()

	log.Printf("Kafka Consumer connetced to broker: %s, topic: %s, group: %s", brokers, topic, groupID)
	log.Println("Ready to consume messages...")

	// 3. Start consuming transactions
	for {
		ctx := context.Background()

		// ReadMessage is blocking until a message is received or the context is cancelled
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("ERROR reading messages: %v\n", err)
			continue
		}

		// 4. Placeholder logic for processing (replacing old heartbeat)
		log.Printf("✅ RECEIVED: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
