package loan

import (
	"github.com/google/uuid"
	"github.com/mas-dhimas/amartha/internal/payment"
)

type Service interface {
	GetLoanOutstanding(loanID string) (int64, error)
	MakeLoan(loan LoanRequest) error
}

type loanService struct {
	repo Repository
}

// NewLoanService creates a new Loan service.
func NewLoanService(repo Repository) Service {
	return &loanService{repo: repo}
}

func (s *loanService) GetLoanOutstanding(loanID string) (int64, error) {
	return s.repo.GetLoanOutstanding(loanID)
}

func (s *loanService) MakeLoan(req LoanRequest) error {
	var payments []payment.Payment
	loan := Loan{
		CustomerID:     uuid.MustParse(req.CustomerID),
		Principal:      req.Principal,
		TermWeeks:      req.TermWeeks,
		CurrentWeek:    1,
		InterestRate:   10.00, // Fixed interest rate of 10%
		Status:         "active",
		TotalRepayment: int64(float64(req.Principal) * 1.10), // Principal + 10% interest
		Outstanding:    int64(float64(req.Principal) * 1.10),
	}

	weeklyPayment := loan.TotalRepayment / int64(loan.TermWeeks)
	for week := 1; week <= loan.TermWeeks; week++ {
		payments = append(payments, payment.Payment{
			Week:      week,
			DueAmount: weeklyPayment,
		})
	}

	err := s.repo.InsertLoan(loan, payments)
	if err != nil {
		return err
	}

	return nil
}
