package payment

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID         uuid.UUID  `json:"id"`
	LoanID     uuid.UUID  `json:"loan_id"`
	Week       int        `json:"week"`
	DueAmount  int64      `json:"due_amount"`
	PaidAmount int64      `json:"paid_amount"`
	Paid       bool       `json:"paid"`
	PaidAt     *time.Time `json:"paid_at"` // nullable
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type PaymentRequest struct {
	LoanID string `json:"loan_id"`
	Amount int64  `json:"amount"`
}
