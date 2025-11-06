-- FlowBank Database Schema

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0, -- Balance in cents
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL, -- Amount in cents (positive for credit, negative for debit)
    type VARCHAR(20) NOT NULL CHECK (type IN ('debit', 'credit')),
    description TEXT,
    merchant_name VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_created ON transactions(user_id, created_at DESC);

-- Insert some seed users for testing
INSERT INTO users (id, email, name, balance) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 'alice@flowbank.com', 'Alice Johnson', 100000),
    ('550e8400-e29b-41d4-a716-446655440001', 'bob@flowbank.com', 'Bob Smith', 50000),
    ('550e8400-e29b-41d4-a716-446655440002', 'charlie@flowbank.com', 'Charlie Davis', 75000)
ON CONFLICT (email) DO NOTHING;

-- Function to update user balance
CREATE OR REPLACE FUNCTION update_user_balance()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'completed' THEN
        UPDATE users 
        SET balance = balance + NEW.amount,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.user_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update balance when transaction is completed
CREATE TRIGGER transaction_completed_trigger
    AFTER INSERT OR UPDATE OF status ON transactions
    FOR EACH ROW
    WHEN (NEW.status = 'completed')
    EXECUTE FUNCTION update_user_balance();

-- Create a view for easier transaction queries
CREATE OR REPLACE VIEW user_transaction_history AS
SELECT 
    t.id,
    t.user_id,
    u.name as user_name,
    u.email,
    t.amount,
    t.type,
    t.description,
    t.merchant_name,
    t.status,
    t.created_at
FROM transactions t
JOIN users u ON t.user_id = u.id
ORDER BY t.created_at DESC;