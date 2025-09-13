package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mas-dhimas/amartha/internal/loan"
)

type loanHandler struct {
	loanService loan.Service
}

func NewLoanHandler(loanService loan.Service) *loanHandler {
	return &loanHandler{loanService: loanService}
}

// GetLoanOutstanding handles the retrieval of the outstanding amount for a loan.
func (h *loanHandler) GetLoanOutstanding(e echo.Context) error {
	loanID := e.QueryParam("loan_id")
	outstanding, err := h.loanService.GetLoanOutstanding(loanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve articles due to internal error")
	}

	return e.JSON(http.StatusOK, map[string]map[string]any{
		"data": {
			"loan_id":            loanID,
			"outstanding_amount": outstanding,
		},
	})
}

// MakeLoan handles the creation of a new loan.
func (h *loanHandler) MakeLoan(e echo.Context) error {
	var req loan.LoanRequest

	if err := e.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload or malformed JSON")
	}

	errUUID := uuid.Validate(req.CustomerID)
	if errUUID != nil || req.Principal == 0 || req.TermWeeks == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields: customer_id, principal, and term_weeks are mandatory")
	}

	id, err := h.loanService.MakeLoan(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create loan due to internal error")
	}

	// Implementation for making a loan would go here
	return e.JSON(http.StatusOK, map[string]map[string]any{
		"data": {
			"id":      id,
			"message": "Loan created successfully",
		},
	})
}
