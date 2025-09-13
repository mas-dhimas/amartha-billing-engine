package loan

import (
	"database/sql"

	"github.com/mas-dhimas/amartha/internal/payment"
	"github.com/sirupsen/logrus"
)

// Repository defines the interface for loan data operations.
type Repository interface {
	GetLoanOutstanding(loanID string) (int64, error)
	InsertLoan(loan Loan, payments []payment.Payment) (string, error)
}

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository.
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// GetLoanOutstanding retrieves the outstanding amount for a given loan ID.
func (r *postgresRepository) GetLoanOutstanding(loanID string) (int64, error) {
	var amount int64
	query := `SELECT outstanding FROM loans WHERE id = $1`
	err := r.db.QueryRow(query, loanID).Scan(&amount)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// InsertLoan inserts a new loan and its associated payment schedules into the database.
func (r *postgresRepository) InsertLoan(loan Loan, payments []payment.Payment) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var loanID string
	query := `INSERT INTO loans (customer_id, principal, interest_rate, term_weeks, total_repayment, outstanding, status, current_week) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err = tx.QueryRow(query, loan.CustomerID, loan.Principal, loan.InterestRate, loan.TermWeeks, loan.TotalRepayment, loan.Outstanding, loan.Status, loan.CurrentWeek).Scan(&loanID)
	if err != nil {
		logrus.Errorf("Failed to insert loan: %v", err)
		return "", err
	}

	stmt, err := tx.Prepare(`INSERT INTO payments (loan_id, week, due_amount) VALUES ($1, $2, $3)`)
	if err != nil {
		logrus.Errorf("Failed to prepare payment insert statement: %v", err)
		return "", err
	}
	defer stmt.Close()

	for _, p := range payments {
		_, err := stmt.Exec(loanID, p.Week, p.DueAmount)
		if err != nil {
			logrus.Errorf("Failed to insert payment: %v", err)
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return loanID, nil
}
