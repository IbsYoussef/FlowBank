package model

import "time"

// Trasaction represents a single financial movement between accounts.
type Transaction struct {
	ID            string    `json:"id" db:"transaction_id"`     // Unique ID for the transaction
	Type          string    `json:"type" db:"transaction_type"` // e.g., "deposit", "withdrawel", "transfer"
	SourceAccount string    `json:"source_account" db:"source_account_id"`
	TargetAccount string    `json:"target_account" db:"target_account_id"` // Only used for transfer type
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"` // e.g., "USD", "GBP"
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	Status        string    `json:"status" db:"status"`
}

// Account represents a bank account's current state (for PostgreSQL)
type Account struct {
	ID      string  `json:"id" db:"account_id"`
	Balance float64 `json:"balance" db:"current_balance"`
	Owner   string  `json:"owner" db:"owner_name"`
}
