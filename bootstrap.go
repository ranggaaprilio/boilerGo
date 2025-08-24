package main

import (
	"github.com/ranggaaprilio/boilerGo/app/v1/modules/user"
	"github.com/ranggaaprilio/boilerGo/config"
	appLogger "github.com/ranggaaprilio/boilerGo/internal/logger"
)

// Bootstrap initializes the application's database and performs necessary migrations
func Bootstrap() error {
	// Initialize simple logger for bootstrap
	bootstrapLogger := appLogger.SimpleLogger("bootstrap")

	bootstrapLogger.Info("Starting bootstrap process...")

	// Get database connection
	db := config.CreateCon()
	if db == nil {
		bootstrapLogger.Fatal("Failed to get database connection")
		return nil
	}

	bootstrapLogger.Info("Running database migrations...")

	// Run migrations for all models
	if err := db.AutoMigrate(&user.User{}); err != nil {
		bootstrapLogger.Error("Failed to migrate User model", "error", err)
		return err
	}

	bootstrapLogger.Info("Database migrations completed successfully")

	// Add any seed data or additional bootstrap logic here
	if err := seedData(db, bootstrapLogger); err != nil {
		bootstrapLogger.Error("Failed to seed data", "error", err)
		return err
	}

	bootstrapLogger.Info("Bootstrap process completed successfully")
	return nil
}

// seedData adds initial data to the database if needed
func seedData(db interface{}, logger *appLogger.LogrusLogger) error {
	// Add any initial data seeding logic here
	logger.Info("Checking for seed data requirements...")

	// Example: Create default admin user, default settings, etc.
	// This is where you'd add any initial data your application needs

	logger.Info("Seed data check completed")
	return nil
}
