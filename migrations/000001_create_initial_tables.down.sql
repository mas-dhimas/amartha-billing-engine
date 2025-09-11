-- Drop indexes first
DROP INDEX IF EXISTS idx_payments_paid;
DROP INDEX IF EXISTS idx_payments_loan_id;

DROP INDEX IF EXISTS idx_loans_status;
DROP INDEX IF EXISTS idx_loans_customer_id;

-- Drop tables (in reverse dependency order)
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS customers;
