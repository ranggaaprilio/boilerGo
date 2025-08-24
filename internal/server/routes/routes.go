package routes

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/app/v1/handler"
	"github.com/ranggaaprilio/boilerGo/app/v1/modules/user"
	"github.com/ranggaaprilio/boilerGo/config"
	"github.com/ranggaaprilio/boilerGo/exception"
	"github.com/ranggaaprilio/boilerGo/internal/server/routes/v1"
)

// SetupRoutes configures all application routes
func SetupRoutes(e *echo.Echo) {
	// Setup Swagger documentation
	SetupSwagger(e)

	// Setup API versioned routes
	setupV1Routes(e)

	// Export routes to JSON file for documentation
	exportRoutes(e)
}

// setupV1Routes configures version 1 API routes
func setupV1Routes(e *echo.Echo) {
	// Create v1 group
	v1 := e.Group("/api/v1")

	// Setup welcome routes
	routes.SetupWelcomeRoutes(v1)

	// Setup user routes
	setupUserRoutes(v1)
}

// setupUserRoutes configures user-related routes
func setupUserRoutes(v1 *echo.Group) {
	// Initialize user dependencies
	db := config.CreateCon()
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	// Setup user routes
	routes.SetupUserRoutes(v1, userHandler)
}

// exportRoutes saves all routes to a JSON file for documentation
func exportRoutes(e *echo.Echo) {
	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		exception.PanicIfNeeded(err)
		return
	}

	if err := ioutil.WriteFile("routes/routes.json", data, 0644); err != nil {
		// Log error but don't panic for file write issues
		// Could use proper logging here
		return
	}
}
