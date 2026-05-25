package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"flowbank/internal/model"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// Producer abstracts the kafka client for publishing messages.
type Producer struct {
	writer *kafka.Writer
	topic  string
}

// NewProducer creates and initializes a new Kafka Producer instance.
func NewProducer(brokers []string, topic string) (*Producer, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// If SASL Credentials are set, use SASL_SSL for Aiven
	saslUser := os.Getenv("KAFKA_SASL_USERNAME")
	saslPass := os.Getenv("KAFKA_SASL_PASSWORD")
	caCert := os.Getenv("KAFKA_CA_CERT")

	if saslUser != "" && saslPass != "" && caCert != "" {
		tlsConfig, err := NewTLSConfig(caCert)
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}
		dialer.TLS = tlsConfig
		dialer.SASLMechanism = plain.Mechanism{
			Username: saslUser,
			Password: saslPass,
		}
	}

	// A kafka.Writer is used for producing messages.
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		BatchTimeout: 10 * time.Millisecond, // Wait up to 10ms to batch messages
		Dialer:       dialer,
	})

	fmt.Printf("Kafka Producer initialized for topic: %s on brokers: %v\n", topic, brokers)

	return &Producer{
		writer: w,
		topic:  topic,
	}, nil
}

// Close gracefully closes the producer connection.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// PublishTransaction serializes a model.Transaction struct and publishes it to the topic.
func (p *Producer) PublishTransaction(ctx context.Context, tx model.Transaction) error {
	// 1. Serialize the Transaction struct to a JSON byte array
	value, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction to JSON: %w", err)
	}

	// 2. Create the Kafka Message
	msg := kafka.Message{
		// Use the Transaction ID as the key to ensure all transactions for the same ID
		// land on the same partition (good for ordering and idempotency checks).
		Key:   []byte(tx.ID.String()),
		Value: value,
		Time:  tx.CreatedAt,
	}

	// 3. Write the message to the Kafka stream (RedPanda)
	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	return nil
}
