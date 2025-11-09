package db

import (
	"context"
	"flowbank/internal/model"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB encapsulate the PostgreSQL connection pool
type DB struct {
	pool *pgxpool.Pool
}

// NewConnection initialised a database connection pool.
func NewConnection(ctx context.Context, databaseURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL database.")
	return &DB{pool: pool}, nil
}

// Close gracefully closes the database connection pool.
func (d *DB) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

// UpdateUserBalanceAndSaveTransaction performs an atomic database transaction.
func (d *DB) UpdateUserBalanceAndSaveTransaction(ctx context.Context, tx *model.Transaction) error {
	// Start a new databse transaction
	dbTx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not beging transaction: %w", err)
	}
	// Defer a rollback in case the transaction is not commited or panics
	defer func() {
		if err != nil {
			dbTx.Rollback(ctx)
		}
	}()

	// 1. Get the current user balance and LOCK the row (FOR UPDATE)
	var currentBalance int64
	err = dbTx.QueryRow(ctx, `
		SELECT current_balance 
		FROM users 
		WHERE user_id = $1 
		FOR UPDATE`,
		tx.UserID).Scan(&currentBalance)

	if err != nil {
		if err == pgx.ErrNoRows {
			tx.Status = "FAILED"
			log.Printf("⛔ User not found for transaction %s (%s). Marking as FAILED.", tx.ID, tx.UserID)
		} else {
			return fmt.Errorf("failed to retrieve user %s for update: %w", tx.UserID, err)
		}
	} else {
		// 2. Calculate the new balance
		var newBalance int64
		if tx.Type == "credit" {
			newBalance = currentBalance + tx.Amount
		} else if tx.Type == "debit" {
			newBalance = currentBalance - tx.Amount
		}

		// 3. Check for overdraft
		if tx.Type == "debit" && newBalance < 0 {
			tx.Status = "FAILED"
			log.Printf("⛔ Overdraft detected for user %s on transaction %s (Current Balance %d). Marking as failed", tx.UserID, tx.ID, currentBalance)
			// Do NOT update the user's balance
		} else {
			// 4. Update the user's balance
			_, err = dbTx.Exec(ctx, `
				UPDATE users 
				SET current_balance = $1, updated_at = $2 
				WHERE user_id = $3`,
				newBalance, time.Now(), tx.UserID)

			if err != nil {
				return fmt.Errorf("failed to update the user balance: %w", err)
			}
			tx.Status = "COMPLETED"
		}
	}

	// 5. Save the transaction record with its final status
	if saveErr := d.SaveTransaction(ctx, dbTx, tx); saveErr != nil {
		err = fmt.Errorf("failed to save transaction record: %w", saveErr)
		return err // Return the erro to trigger the defer rollback
	}

	// 6. Commit the entire operation
	err = dbTx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SaveTransaction inserts a transaction record into the database.
// It is called within the larger UpdateUserBalanceAndSaveTransaction function.
func (d *DB) SaveTransaction(ctx context.Context, dbTx pgx.Tx, tx *model.Transaction) error {
	_, err := dbTx.Exec(ctx, `
		INSERT INTO transactions (transaction_id, user_id, amount, transaction_type, description, merchant_name, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		tx.ID, tx.UserID, tx.Amount, tx.Type, tx.Description, tx.MerchantName, tx.Status, tx.CreatedAt,
	)
	return err
}

// GetUser retrieves a user's current account details (including balance).
func (d *DB) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User

	row := d.pool.QueryRow(ctx, `
		SELECT user_id, email, user_name, current_balance, created_at, updated_at 
		FROM users 
		WHERE user_id = $1`, userID)

	err := row.Scan(
		&userID,
		&user.Email,
		&user.Name,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return &user, nil
}

// GetUserTransactions retrieves a list of transaction records for a given user.
func (d *DB) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]model.Transaction, error) {
	rows, err := d.pool.Query(ctx, `
		SELECT 
			transaction_id, user_id, amount, transaction_type, description, merchant_name, status, created_at
		FROM transactions 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT 20`, // Limit to 20 for a reasonable history response
		userID)

	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction

	for rows.Next() {
		var tx model.Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Amount,
			&tx.Type,
			&tx.Description,
			&tx.MerchantName,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction row: %w", err)
		}
		transactions = append(transactions, tx)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error after iterating through transaction rows: %w", rows.Err())
	}

	return transactions, nil
}
