CREATE TABLE IF NOT EXISTS billings (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    loan_id INTEGER NOT NULL,
    loan_amount INTEGER NOT NULL,
    loan_weeks INTEGER NOT NULL,
    loan_interest INTEGER NOT NULL,
    outstanding INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_billings_loan_id UNIQUE (loan_id)
);
