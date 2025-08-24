package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// SetupWelcomeRoutes configures welcome endpoints for API v1
func SetupWelcomeRoutes(v1 *echo.Group) {
	v1.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome To Go RestApi V1")
	})
}
