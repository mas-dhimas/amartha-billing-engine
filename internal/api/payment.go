package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mas-dhimas/amartha/internal/payment"
)

type paymentHandler struct {
	paymentService payment.Service
}

func NewPaymentHandler(paymentService payment.Service) *paymentHandler {
	return &paymentHandler{paymentService: paymentService}
}

// MakePayment handles the payment processing for a loan.
func (h *paymentHandler) MakePayment(e echo.Context) error {
	var req payment.PaymentRequest

	if err := e.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload or malformed JSON")
	}

	errUUID := uuid.Validate(req.LoanID)
	if errUUID != nil || req.Amount == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields: loan_id and amount are mandatory")
	}

	err := h.paymentService.MakePayment(req.LoanID, req.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process payment due to internal error")
	}

	return e.JSON(http.StatusOK, map[string]map[string]any{
		"data": {
			"message": "Loan created successfully",
		},
	})
}
