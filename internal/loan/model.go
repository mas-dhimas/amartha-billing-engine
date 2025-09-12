package loan

import (
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	ID             uuid.UUID `json:"id"`
	CustomerID     uuid.UUID `json:"customer_id"`
	Principal      int64     `json:"principal"`     // BIGINT
	InterestRate   float64   `json:"interest_rate"` // NUMERIC(5,2)
	TotalRepayment int64     `json:"total_repayment"`
	Outstanding    int64     `json:"outstanding"`
	TermWeeks      int       `json:"term_weeks"`
	CurrentWeek    int       `json:"current_week"`
	Status         string    `json:"status"` // "active", "closed", "defaulted"
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LoanRequest struct {
	CustomerID string `json:"customer_id"`
	Principal  int64  `json:"principal"`
	TermWeeks  int    `json:"term_weeks"`
}
