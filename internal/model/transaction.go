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

// DashboardStats contains aggregated statistics for the dashboard
type DashboardStats struct {
	TotalTransactions  int64              `json:"total_transactions"`
	TotalVolume        int64              `json:"total_volume"`
	AverageAmount      int64              `json:"average_amount"`
	CreditCount        int64              `json:"credit_count"`
	DebitCount         int64              `json:"debit_count"`
	MoneyIn            int64              `json:"money_in"`
	MoneyOut           int64              `json:"money_out"`
	TopMerchants       []MerchantStats    `json:"top_merchants"`
	RecentTransactions []TransactionDTO   `json:"recent_transactions"`
	HourlyVolume       []HourlyVolumeData `json:"hourly_volume"`
}

// MerchantStats represents transaction counts per merchant
type MerchantStats struct {
	MerchantName string `json:"merchant_name"`
	Count        int64  `json:"count"`
	TotalAmount  int64  `json:"total_amount"`
}

// TransactionDTO is a simplified transaction for the dashboard
type TransactionDTO struct {
	ID           string `json:"id"`
	UserName     string `json:"user_name"`
	Amount       int64  `json:"amount"`
	Type         string `json:"type"`
	MerchantName string `json:"merchant_name"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

// HourlyVolumeData represents transaction volume per hour
type HourlyVolumeData struct {
	Hour        string `json:"hour"`
	CreditTotal int64  `json:"credit_total"`
	DebitTotal  int64  `json:"debit_total"`
	Count       int64  `json:"count"`
}
