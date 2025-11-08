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