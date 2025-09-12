package payment

import "fmt"

type Service interface {
	MakePayment(loanID string, amount int64) error
}

type paymentService struct {
	repo Repository
}

// NewPaymentService creates a new Payment service.
func NewPaymentService(repo Repository) Service {
	return &paymentService{repo: repo}
}

func (s *paymentService) MakePayment(loanID string, amount int64) error {
	outstandingWeeks, outstandingAmount, err := s.repo.GetPaymentOutstanding(loanID)
	if err != nil {
		return err
	}

	if outstandingAmount != 0 && amount != outstandingAmount {
		return fmt.Errorf("payment amount %d is not equal to outstanding amount %d", amount, outstandingAmount)
	}

	return s.repo.InsertPayment(loanID, amount, outstandingWeeks)
}
