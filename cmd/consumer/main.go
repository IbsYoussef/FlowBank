package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println("🏦 FlowBank Consumer starting...")

	// TODO: Initialize Kafka consumer
	// TODO: Connect to database
	// TODO: Start consuming transactions
	// TODO: Validate and persist to DB

	for {
		fmt.Println("Consumer heartbeat:", time.Now().Format(time.RFC3339))
		time.Sleep(5 * time.Second)
	}
}
