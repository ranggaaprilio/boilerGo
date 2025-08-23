package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/app/v1/handler"
)

// SetupWelcomeRoutes configures welcome endpoints for API v1
func SetupWelcomeRoutes(v1 *echo.Group) {
	v1.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome To Go RestApi V1")
	})
}

// SetupUserRoutes configures user-related endpoints for API v1
func SetupUserRoutes(v1 *echo.Group, userHandler *handler.UserHandler) {
	// User routes group
	users := v1.Group("/users")

	// User endpoints
	users.POST("", userHandler.RegisterUser)
	users.GET("/:id", userHandler.GetUser)
}
