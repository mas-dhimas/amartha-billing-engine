package payment

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Repository interface {
	InsertPayment(loanID string, amount, weeksToPay int64) error
	GetPaymentOutstanding(loanID string) (weeksToPay, amount int64, err error)
	GetPaymentSchedule(loanID string) ([]Payment, error)
}

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository.
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// InsertPayment records a payment for a given loan ID and updates the loan's outstanding amount and current week.
// It marks the next 'weeksToPay' unpaid payments as paid.
func (r *postgresRepository) InsertPayment(loanID string, amount, weeksToPay int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE payments SET paid = true, paid_amount = $1, paid_at = $2 `
	query += `WHERE loan_id = $3 AND paid = false ORDER BY week ASC LIMIT $4`
	res, err := tx.Exec(query, amount, time.Now(), loanID, weeksToPay)
	if err != nil {
		logrus.Errorf("Failed to update payments: %v", err)
		return err
	}

	rowsUpdated, _ := res.RowsAffected()
	if rowsUpdated != weeksToPay {
		return fmt.Errorf("only %d payments could be applied, expected %d", rowsUpdated, weeksToPay)
	}

	query = `UPDATE loans SET outstanding = outstanding - $1, current_week = current_week + $2, `
	query += `status = CASE WHEN (outstanding - $1) <= 0 THEN 'closed' ELSE status END, updated_at = $3 `
	query += `WHERE id = $4`
	_, err = tx.Exec(query, int64(amount), weeksToPay, time.Now(), loanID)
	if err != nil {
		logrus.Errorf("Failed to update loan: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetPaymentOutstanding retrieves the number of unpaid weeks for a given loan ID.
func (r *postgresRepository) GetPaymentOutstanding(loanID string) (weeksToPay, amount int64, err error) {
	query := `SELECT COUNT(p.id), SUM(p. due_amount) FROM payments p JOIN loans l ON l.id = p.loan_id `
	query += `WHERE p.loan_id = $1 AND p.paid = false AND p.week <= l.current_week ORDER BY p.week ASC`
	err = r.db.QueryRow(query, loanID).Scan(&weeksToPay, &amount)
	if err != nil {
		logrus.Errorf("Error querying payment outstanding: %v\n", err)
		return 0, 0, err
	}
	return
}

// GetPaymentSchedule retrieves the payment schedule for a given loan ID.
func (r *postgresRepository) GetPaymentSchedule(loanID string) ([]Payment, error) {
	query := `SELECT id, loan_id, week, due_amount, paid_amount, paid, paid_at FROM payments WHERE loan_id = $1`

	rows, err := r.db.Query(query, loanID)
	if err != nil {
		logrus.Errorf("Failed to query payment schedule: %v", err)
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.LoanID, &p.Week, &p.DueAmount, &p.PaidAmount, &p.Paid, &p.PaidAt); err != nil {
			logrus.Errorf("Failed to scan payment row: %v", err)
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}
