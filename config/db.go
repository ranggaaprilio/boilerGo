package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ranggaaprilio/boilerGo/exception"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var err error

// Database connection configuration
type DatabaseConfig struct {
	MaxRetries  int
	RetryDelay  time.Duration
	MaxIdleConn int
	MaxOpenConn int
}

// DefaultDatabaseConfig returns default database configuration
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		MaxRetries:  5,
		RetryDelay:  5 * time.Second,
		MaxIdleConn: 10,
		MaxOpenConn: 100,
	}
}

// DbInit initializes the database connection with retry mechanism
func DbInit() {
	conf := Loadconf()
	dbConfig := DefaultDatabaseConfig()

	// Initialize logger for database operations
	dbLogger := log.New(os.Stdout, "[Database] ", log.LstdFlags)

	// Log connection attempt
	dbLogger.Printf("Attempting to connect to database at %s:%s as user %s",
		conf.Database.DbHost, conf.Database.DbPort, conf.Database.DbUsername)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.DbUsername,
		conf.Database.DbPassword,
		conf.Database.DbHost,
		conf.Database.DbPort,
		conf.Database.DbName)

	// Configure GORM with custom logger
	gormConfig := &gorm.Config{
		Logger: logger.New(
			dbLogger,
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	}

	// Retry loop for database connection
	for i := 0; i < dbConfig.MaxRetries; i++ {
		db, err = gorm.Open(mysql.Open(connectionString), gormConfig)
		if err == nil {
			// Configure connection pool
			if sqlDB, poolErr := db.DB(); poolErr == nil {
				sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
				sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)
				sqlDB.SetConnMaxLifetime(time.Hour)
			}

			dbLogger.Println("Successfully connected to database")
			return
		}

		dbLogger.Printf("Failed to connect to database (attempt %d/%d): %v",
			i+1, dbConfig.MaxRetries, err)

		if i < dbConfig.MaxRetries-1 {
			dbLogger.Printf("Retrying in %v...", dbConfig.RetryDelay)
			time.Sleep(dbConfig.RetryDelay)
		}
	}

	// If we get here, we've exhausted all retries
	exception.PanicIfNeeded(fmt.Errorf("failed to connect to database after %d attempts: %w",
		dbConfig.MaxRetries, err))
}

// GetDB returns the database instance (alternative to CreateCon)
func GetDB() *gorm.DB {
	if db == nil {
		exception.PanicIfNeeded(fmt.Errorf("database connection not initialized"))
	}
	return db
}

// CloseDB closes the database connection gracefully
func CloseDB() error {
	if db != nil {
		if sqlDB, err := db.DB(); err == nil {
			return sqlDB.Close()
		}
	}
	return nil
}

// PingDB checks if the database connection is alive
func PingDB() error {
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// CreateCon return var db
func CreateCon() *gorm.DB {
	return db
}
