package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mas-dhimas/amartha/internal/customer"
)

type customerHandler struct {
	customerService customer.Service
}

func NewCustomerHandler(customerService customer.Service) *customerHandler {
	return &customerHandler{customerService: customerService}
}

// CheckIsCustomerDelinquent handles the check for delinquency status of a customer.
func (h *customerHandler) CheckIsCustomerDelinquent(e echo.Context) error {
	customerID := e.QueryParam("customer_id")
	isDelinquent, err := h.customerService.CheckIsCustomerDelinquent(customerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve articles due to internal error")
	}

	return e.JSON(http.StatusOK, map[string]map[string]any{
		"data": {
			"customer_id":   customerID,
			"is_delinquent": isDelinquent,
		},
	})
}
