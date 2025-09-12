package customer

import "database/sql"

type Repository interface {
	CheckIsCustomerDelinquent(customerID string) (bool, error)
	GetCustomerByID(customerID string) (*Customer, error)
}

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository.
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CheckIsCustomerDelinquent(customerID string) (bool, error) {
	var unpaidTotal int64
	query := `SELECT COUNT(*) AS unpaid_count FROM payments p `
	query += `JOIN loans l ON l.id = p.loan_id JOIN customers c ON c.id = l.customer_id `
	query += `WHERE c.id = $1 AND p.week IN (l.current_week - 1, l.current_week - 2) AND p.paid = false`

	err := r.db.QueryRow(query, customerID).Scan(&unpaidTotal)
	if err != nil {
		return false, err
	}

	return unpaidTotal >= 2, nil
}

func (r *postgresRepository) GetCustomerByID(customerID string) (*Customer, error) {
	return nil, nil
}
