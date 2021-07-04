package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handlers
func Welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome To Go RestApi V1")
}
