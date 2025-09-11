CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE customers (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name	VARCHAR not null
);

CREATE TABLE loans (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id     UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    principal       BIGINT NOT NULL CHECK (principal > 0),
    interest_rate   NUMERIC(5,2) NOT NULL,     -- store % (e.g. 10.00)
    total_repayment BIGINT NOT NULL CHECK (total_repayment > 0),
    outstanding     BIGINT NOT NULL CHECK (outstanding >= 0),
    term_weeks      INT NOT NULL CHECK (term_weeks > 0),
    current_week    INT NOT NULL DEFAULT 0 CHECK (current_week >= 0),
    status          TEXT NOT NULL CHECK (status IN ('active','closed','defaulted')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_loans_customer_id ON loans(customer_id);
CREATE INDEX idx_loans_status ON loans(status);

CREATE TABLE payments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id     UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    week        INT NOT NULL CHECK (week >= 0),
    due_amount  BIGINT NOT NULL CHECK (due_amount > 0),
    paid_amount BIGINT NOT NULL DEFAULT 0,
    paid        BOOLEAN NOT NULL DEFAULT false,
    paid_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (loan_id, week)
);

CREATE INDEX idx_payments_loan_id ON payments(loan_id);
CREATE INDEX idx_payments_paid ON payments(paid);

