package routes

import (
	"github.com/labstack/echo/v4"
	_ "github.com/ranggaaprilio/boilerGo/docs" // Import swagger docs
	echoSwagger "github.com/swaggo/echo-swagger"
)

// SetupSwagger configures Swagger documentation endpoints
func SetupSwagger(e *echo.Echo) {
	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())
}
