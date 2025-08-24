package config

import (
	"fmt"
	"strconv"

	"github.com/ranggaaprilio/boilerGo/exception"
	appLogger "github.com/ranggaaprilio/boilerGo/internal/logger"
	"github.com/spf13/viper"
)

// Configurations represents the main configuration structure
type Configurations struct {
	Server   ServerConfigurations `mapstructure:"server" validate:"required"`
	Database DbConfigurations     `mapstructure:"database" validate:"required"`
	App      AppConfigurations    `mapstructure:"app"`
}

// ServerConfigurations holds server-related settings
type ServerConfigurations struct {
	Name         string `mapstructure:"name" validate:"required"`
	Port         string `mapstructure:"port" validate:"required"`
	ReadTimeout  int    `mapstructure:"read_timeout" default:"30"`
	WriteTimeout int    `mapstructure:"write_timeout" default:"30"`
	Environment  string `mapstructure:"environment" default:"development"`
}

// DbConfigurations holds database-related settings
type DbConfigurations struct {
	DbUsername string `mapstructure:"dbusername" validate:"required"`
	DbPassword string `mapstructure:"dbpassword" validate:"required"`
	DbHost     string `mapstructure:"dbhost" validate:"required"`
	DbPort     string `mapstructure:"dbport" validate:"required"`
	DbName     string `mapstructure:"dbname" validate:"required"`
	DbSSL      string `mapstructure:"dbssl" default:"disable"`
}

// AppConfigurations holds general application settings
type AppConfigurations struct {
	LogLevel    string `mapstructure:"log_level" default:"info"`
	Debug       bool   `mapstructure:"debug" default:"false"`
	SecretKey   string `mapstructure:"secret_key"`
	ServiceName string `mapstructure:"service_name" default:"BoilerGo"`
}

// ConfigLoader handles configuration loading and validation
type ConfigLoader struct {
	logger *appLogger.LogrusLogger
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		logger: appLogger.SimpleLogger("config"),
	}
}

// Loadconf loads and validates the application configuration
func Loadconf() Configurations {
	loader := NewConfigLoader()
	return loader.Load()
}

// Load loads the configuration from file and environment variables
func (cl *ConfigLoader) Load() Configurations {
	cl.logger.Info("Loading application configuration...")

	// Initialize viper
	cl.setupViper()

	// Set up environment variable mappings
	cl.setupEnvironmentBindings()

	// Set default values
	cl.setDefaults()

	var configuration Configurations

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cl.logger.Info("Config file not found, using environment variables and defaults")
		} else {
			cl.logger.Error("Error reading config file", "error", err)
			exception.PanicIfNeeded(err)
		}
	} else {
		cl.logger.Info("Using config file", "file", viper.ConfigFileUsed())
	}

	// Unmarshal configuration
	if err := viper.Unmarshal(&configuration); err != nil {
		cl.logger.Error("Error unmarshaling configuration", "error", err)
		exception.PanicIfNeeded(err)
	}

	// Validate configuration
	if err := cl.validateConfiguration(&configuration); err != nil {
		cl.logger.Error("Configuration validation failed", "error", err)
		exception.PanicIfNeeded(err)
	}

	cl.logger.Info("Configuration loaded and validated successfully")
	return configuration
}

// setupViper configures viper settings
func (cl *ConfigLoader) setupViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/boilergo/")
	viper.AutomaticEnv()
}

// setupEnvironmentBindings maps environment variables to config keys
func (cl *ConfigLoader) setupEnvironmentBindings() {
	envMappings := map[string]string{
		"server.name":         "SERVER_NAME",
		"server.port":         "SERVER_PORT",
		"server.environment":  "ENVIRONMENT",
		"database.dbusername": "DB_USER",
		"database.dbpassword": "DB_PASSWORD",
		"database.dbhost":     "DB_HOST",
		"database.dbport":     "DB_PORT",
		"database.dbname":     "DB_NAME",
		"database.dbssl":      "DB_SSL",
		"app.log_level":       "LOG_LEVEL",
		"app.debug":           "DEBUG",
		"app.secret_key":      "SECRET_KEY",
		"app.service_name":    "SERVICE_NAME",
	}

	for configKey, envVar := range envMappings {
		if err := viper.BindEnv(configKey, envVar); err != nil {
			cl.logger.Warn("Failed to bind environment variable", "env_var", envVar, "error", err)
		}
	}
}

// setDefaults sets default configuration values
func (cl *ConfigLoader) setDefaults() {
	viper.SetDefault("server.name", "BoilerGo")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("database.dbssl", "disable")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.debug", false)
	viper.SetDefault("app.service_name", "BoilerGo")
}

// validateConfiguration performs basic validation on the loaded configuration
func (cl *ConfigLoader) validateConfiguration(config *Configurations) error {
	// Validate server configuration
	if config.Server.Name == "" {
		return fmt.Errorf("server name cannot be empty")
	}

	if config.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	// Validate port is a valid number
	if _, err := strconv.Atoi(config.Server.Port); err != nil {
		return fmt.Errorf("server port must be a valid number: %v", err)
	}

	// Validate database configuration
	if config.Database.DbHost == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	if config.Database.DbPort == "" {
		return fmt.Errorf("database port cannot be empty")
	}

	if config.Database.DbUsername == "" {
		return fmt.Errorf("database username cannot be empty")
	}

	if config.Database.DbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	// Validate database port is a valid number
	if _, err := strconv.Atoi(config.Database.DbPort); err != nil {
		return fmt.Errorf("database port must be a valid number: %v", err)
	}

	return nil
}

// IsProduction returns true if the application is running in production mode
func (c *Configurations) IsProduction() bool {
	return c.Server.Environment == "production"
}

// IsDevelopment returns true if the application is running in development mode
func (c *Configurations) IsDevelopment() bool {
	return c.Server.Environment == "development"
}
