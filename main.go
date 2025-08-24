// Package main for the BoilerGo API
//
// @title BoilerGo API
// @version 1.0
// @description This is the API documentation for the BoilerGo application.
//
// @contact.name API Support
// @contact.url https://github.com/ranggaaprilio/boilerGo
// @contact.email support@example.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api
// @schemes http https
package main

import (
	_ "github.com/ranggaaprilio/boilerGo/docs" // Import swagger docs
	"github.com/ranggaaprilio/boilerGo/exception"
	"github.com/ranggaaprilio/boilerGo/internal/app"
	appLogger "github.com/ranggaaprilio/boilerGo/internal/logger"
)

func main() {
	defer exception.Catch()

	// Initialize simple logger for main
	mainLogger := appLogger.SimpleLogger("main")

	// Initialize application
	application := app.New()

	// Run bootstrap process
	if err := Bootstrap(); err != nil {
		mainLogger.Fatal("Bootstrap failed", "error", err)
	}

	// Log configuration for debugging
	application.LogConfig()

	// Start application server
	if err := application.Start(); err != nil {
		mainLogger.Fatal("Failed to start application", "error", err)
	}
}
