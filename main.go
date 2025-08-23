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
	"log"
	"os"

	_ "github.com/ranggaaprilio/boilerGo/docs" // Import swagger docs
	"github.com/ranggaaprilio/boilerGo/exception"
	"github.com/ranggaaprilio/boilerGo/internal/app"
)

func main() {
	defer exception.Catch()

	// Initialize logger
	logger := log.New(os.Stdout, "[Main] ", log.LstdFlags)

	// Initialize application
	application := app.New()

	// Run bootstrap process
	if err := Bootstrap(); err != nil {
		logger.Fatal("Bootstrap failed:", err)
	}

	// Log configuration for debugging
	application.LogConfig()

	// Start application server
	if err := application.Start(); err != nil {
		logger.Fatal("Failed to start application:", err)
	}
}
