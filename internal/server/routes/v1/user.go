package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/app/v1/handler"
)

// SetupUserRoutes configures user-related endpoints for API v1
func SetupUserRoutes(v1 *echo.Group, userHandler *handler.UserHandler) {
	// User routes group
	users := v1.Group("/users")

	// User endpoints
	users.POST("", userHandler.RegisterUser)
	users.GET("/:id", userHandler.GetUser)
}
