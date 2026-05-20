-- 1. Users Table: Stores the current state (balance) of each user account.
-- Note: BIGINT is used for the balance, representing the amount in cents.
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    user_name VARCHAR(100) NOT NULL,
    current_balance BIGINT NOT NULL DEFAULT 0, -- Stored in cents
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Function to update the updated_at timestamp automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to fire on user updates
DROP TRIGGER IF EXISTS set_updated_at ON users;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- 2. Transactions Table: Stores the immutable history of all transactions.
-- Note: BIGINT is used for the amount, representing the amount in cents.
CREATE TABLE IF NOT EXISTS transactions (
    transaction_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id),
    amount BIGINT NOT NULL, -- Stored in cents
    transaction_type VARCHAR(20) NOT NULL,
    description VARCHAR(255),
    merchant_name VARCHAR(100),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Initial Dummy Data for testing (Using UUIDs and cents):
-- Alice starts with $1000.00 (100000 cents)
INSERT INTO users (user_id, email, user_name, current_balance)
VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'alice@flowbank.com', 'Alice Smith', 100000),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'bob@flowbank.com', 'Bob Johnson', 50000)
ON CONFLICT (user_id) DO NOTHING;

-- 3. Fraud Scores Table: Stores fraud detection results for each transaction
CREATE TABLE IF NOT EXISTS fraud_scores(
    fraud_score_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(transaction_id),
    risk_score VARCHAR(10) NOT NULL CHECK (risk_score IN ('low', 'medium', 'high')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('clean', 'suspicious', 'flagged')),
    triggered_rules TEXT[] NOT NULL DEFAULT '{}',
    confidence DECIMAL(3,2) NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    processing_time_ms DECIMAL(10,2),
    scored_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for query performance
CREATE INDEX IF NOT EXISTS idx_fraud_scores_transaction_id ON fraud_scores(transaction_id);
CREATE INDEX IF NOT EXISTS idx_fraud_scores_status ON fraud_scores(status);
CREATE INDEX IF NOT EXISTS idx_fraud_scores_scored_at ON fraud_scores(scored_at DESC);

-- Index on transactions table for fraud detection queries
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_created_at ON transactions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_amount ON transactions(user_id, amount);