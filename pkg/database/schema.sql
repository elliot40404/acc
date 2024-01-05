CREATE TABLE IF NOT EXISTS transactions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
	description TEXT NOT NULL,
	amount FLOAT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS transactions_type_idx ON transactions (type);
CREATE INDEX IF NOT EXISTS transactions_created_at_idx ON transactions (created_at);
CREATE INDEX IF NOT EXISTS transactions_amount_idx ON transactions (amount);
CREATE INDEX IF NOT EXISTS transactions_description_idx ON transactions (description);