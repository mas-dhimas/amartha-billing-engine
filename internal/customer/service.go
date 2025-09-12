package customer

type Service interface {
	CheckIsCustomerDelinquent(customerID string) (bool, error)
}

type customerService struct {
	repo Repository
}

// NewCustomerService creates a new Customer service.
func NewCustomerService(repo Repository) Service {
	return &customerService{repo: repo}
}

func (s *customerService) CheckIsCustomerDelinquent(customerID string) (bool, error) {
	return s.repo.CheckIsCustomerDelinquent(customerID)
}
