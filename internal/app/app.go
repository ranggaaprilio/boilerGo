package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/config"
	"github.com/ranggaaprilio/boilerGo/internal/logger"
	"github.com/ranggaaprilio/boilerGo/internal/server"
)

// App represents the main application structure
type App struct {
	server *echo.Echo
	config config.Configurations
	logger *logger.Logger
}

// New creates a new application instance
func New() *App {
	// Load configuration
	conf := config.Loadconf()

	// Initialize structured logger
	appLogger := logger.New(conf)

	// Log configuration load
	appLogger.LogConfigLoad("config.yml", true)

	// Initialize database
	config.DbInit()

	// Initialize server
	srv := server.New(conf, appLogger)

	return &App{
		server: srv,
		config: conf,
		logger: appLogger,
	}
}

// Start starts the application server
func (a *App) Start() error {
	// Log startup information
	a.logger.LogStartup(a.config.Server.Port, a.config)

	// Start server in a goroutine
	go func() {
		address := ":" + a.config.Server.Port
		a.logger.Info("Starting server",
			"service", a.config.Server.Name,
			"port", a.config.Server.Port,
			"address", address)

		if err := a.server.Start(address); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	a.logger.LogShutdown("signal received")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server shutdown error", "error", err)
		return err
	}

	a.logger.Info("Server stopped gracefully")
	return nil
}

// GetServer returns the echo server instance
func (a *App) GetServer() *echo.Echo {
	return a.server
}

// LogConfig logs the current configuration for debugging
func (a *App) LogConfig() {
	a.logger.Info("Application Configuration",
		"server_name", a.config.Server.Name,
		"server_port", a.config.Server.Port,
		"server_environment", a.config.Server.Environment,
		"db_host", a.config.Database.DbHost,
		"db_port", a.config.Database.DbPort,
		"db_user", a.config.Database.DbUsername,
		"db_name", a.config.Database.DbName,
		"log_level", a.config.App.LogLevel,
		"debug_mode", a.config.App.Debug,
	)
}
