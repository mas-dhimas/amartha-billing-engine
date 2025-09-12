package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	customerHandler *customerHandler
	loanHandler     *loanHandler
	paymentHandler  *paymentHandler
}

func NewHandler(cust customerHandler, loan loanHandler, payment paymentHandler) *Handler {
	return &Handler{
		customerHandler: &cust,
		loanHandler:     &loan,
		paymentHandler:  &payment,
	}
}

// RegisterRoutes registers the API routes with the provided router.
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.Add("GET", "/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm alive")
	})

	v1 := e.Group("/api/v1")

	customer := v1.Group("/customer")
	customer.GET("/status", h.customerHandler.CheckIsCustomerDelinquent)

	loan := v1.Group("/loan")
	loan.GET("/outstanding", h.loanHandler.GetLoanOutstanding)
	loan.POST("", h.loanHandler.MakeLoan)

	payment := v1.Group("/payment")
	payment.POST("", h.paymentHandler.MakePayment)
}
