package config

import (
	"fmt"
	"time"

	"github.com/ranggaaprilio/boilerGo/exception"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// DbInit initializes the database connection with retry mechanism
func DbInit() {
	conf := Loadconf()

	// Log connection attempt
	fmt.Printf("Attempting to connect to database at %s:%s as user %s\n",
		conf.Database.DbHost, conf.Database.DbPort, conf.Database.DbUsername)

	// Retry parameters
	maxRetries := 5
	retryDelay := 5 * time.Second

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.DbUsername,
		conf.Database.DbPassword,
		conf.Database.DbHost,
		conf.Database.DbPort,
		conf.Database.DbName)

	// Retry loop for database connection
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
		if err == nil {
			fmt.Println("Successfully connected to database")
			return
		}

		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n",
			i+1, maxRetries, err)

		if i < maxRetries-1 {
			fmt.Printf("Retrying in %v...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	// If we get here, we've exhausted all retries
	exception.PanicIfNeeded(fmt.Errorf("failed to connect to database after %d attempts: %w",
		maxRetries, err))
}

// CreateCon return var db
func CreateCon() *gorm.DB {
	return db
}
