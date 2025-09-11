package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes registers the API routes with the provided router.
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.Add("GET", "/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm alive")
	})
}
