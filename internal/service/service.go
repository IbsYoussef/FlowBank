package service

import (
	"context"
	"encoding/json"
	"flowbank/internal/db"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// Amount is a custom type used for marshalling cents (int64) into dollar strings for the API.
// e.g., 123456 cents becomes "1234.56"
type Amount int64

// MarshalJSON converts the int64 (cents) into a float-formatted JSON string.
func (a Amount) MarshalJSON() ([]byte, error) {
	// Convert cents to a float representing dollars, rounded to 2 decimal places
	dollars := float64(a) / 100.0
	// Format to two decinal places
	s := fmt.Sprintf("%.2f", dollars)
	// Marshall the string as JSON
	return json.Marshal(s)
}

// AccountTransactionDTO is the simplified transacation data transfer object for the API response.
type AccountTransactionDTO struct {
	ID           uuid.UUID `json:"transaction_id"`
	Amount       Amount    `json:"amount"` // Will marshal as a formatted dollars string
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	MerchantName string    `json:"merchant_name"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"created_at"` // Format as string for readability
}

// AccountDetails is the main structure for the API's account query response.
type AccountDetails struct {
	UserID  uuid.UUID               `json:"user_id"`
	Email   string                  `json:"email"`
	Name    string                  `json:"user_name"`
	Balance Amount                  `json:"current_balance"` // Will marshal as a formatted dollar string
	Updated string                  `json:"last_updated"`
	History []AccountTransactionDTO `json:"transaction_history"`
}

// Service contains the business logic and dependencies (like the DB).
type Service struct {
	db *db.DB
}

// NewService creates a new instance of the Service layer.
func NewService(db *db.DB) *Service {
	return &Service{db: db}
}

// GetAccountDetails fetches the full account and transaction history for the API.
func (s *Service) GetAccountDetails(ctx context.Context, userID string) (*AccountDetails, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	// 1. Get User Balance
	user, err := s.db.GetUser(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	// 2. Get Transaction History
	transactions, err := s.db.GetUserTransactions(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction history: %w", err)
	}

	// 3. Map to DTOs and format
	historyDTOs := make([]AccountTransactionDTO, len(transactions))
	for i, tx := range transactions {
		historyDTOs[i] = AccountTransactionDTO{
			ID:           tx.ID,
			Amount:       Amount(math.Abs(float64(tx.Amount))), // Absolute value for display
			Type:         tx.Type,
			Description:  tx.Description,
			MerchantName: tx.MerchantName,
			Status:       tx.Status,
			CreatedAt:    tx.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	// 4. Construct the final response structure
	details := &AccountDetails{
		UserID:  user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Balance: Amount(user.Balance),
		Updated: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		History: historyDTOs,
	}

	return details, nil
}
