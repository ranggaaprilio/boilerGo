package main

import (
	"log"

	"github.com/ranggaaprilio/boilerGo/app/v1/modules/user"
	"github.com/ranggaaprilio/boilerGo/config"
)

// Bootstrap initializes the application's database and performs necessary migrations
func Bootstrap() error {
	logger := log.New(log.Writer(), "[Bootstrap] ", log.LstdFlags)

	logger.Println("Starting bootstrap process...")

	// Get database connection
	db := config.CreateCon()
	if db == nil {
		logger.Fatal("Failed to get database connection")
		return nil
	}

	logger.Println("Running database migrations...")

	// Run migrations for all models
	if err := db.AutoMigrate(&user.User{}); err != nil {
		logger.Printf("Failed to migrate User model: %v", err)
		return err
	}

	logger.Println("Database migrations completed successfully")

	// Add any seed data or additional bootstrap logic here
	if err := seedData(db, logger); err != nil {
		logger.Printf("Failed to seed data: %v", err)
		return err
	}

	logger.Println("Bootstrap process completed successfully")
	return nil
}

// seedData adds initial data to the database if needed
func seedData(db interface{}, logger *log.Logger) error {
	// Add any initial data seeding logic here
	logger.Println("Checking for seed data requirements...")

	// Example: Create default admin user, default settings, etc.
	// This is where you'd add any initial data your application needs

	logger.Println("Seed data check completed")
	return nil
}
