package model

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents an immutable record of a financial event in the system.
// It is the message payload sent across the Kafka (RedPanda) stream.
type Transaction struct {
	ID           uuid.UUID `json:"id" db:"transaction_id"`           // Globally unique identifier for this specific transaction event.
	UserID       uuid.UUID `json:"user_id" db:"user_id"`             // The ID of the user whose account is affected by this transaction.
	Amount       int64     `json:"amount" db:"amount"`               // The monetary value, stored in the smallest currency unit (e.g., cents) for financial accuracy.
	Type         string    `json:"type" db:"transaction_type"`       // The nature of the movement: "debit" (money out) or "credit" (money in).
	Description  string    `json:"description" db:"description"`     // A brief description of the purpose of the transaction.
	MerchantName string    `json:"merchant_name" db:"merchant_name"` // The name of the merchant/counterparty (e.g., "Starbucks").
	Status       string    `json:"status" db:"status"`               // The current processing state: "pending", "completed", or "failed".
	CreatedAt    time.Time `json:"created_at" db:"created_at"`       // Timestamp when the transaction event was first created.
}

// User represents the persistent and mutable state of a user account, stored in PostgreSQL.
// This is the record the Consumer service updates after processing a transaction.
type User struct {
	ID        uuid.UUID `json:"id" db:"user_id"`              // Globally unique identifier for the user account.
	Email     string    `json:"email" db:"email"`             // The user's unique email address.
	Name      string    `json:"name" db:"user_name"`          // The user's full name.
	Balance   int64     `json:"balance" db:"current_balance"` // The user's current account balance, stored in cents (int64).
	CreatedAt time.Time `json:"created_at" db:"created_at"`   // Timestamp when the user account was created.
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`   // Timestamp of the last successful update (e.g., a change in balance).
}
