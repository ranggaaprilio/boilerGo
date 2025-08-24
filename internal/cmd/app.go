package cmd

import (
	"context"
	"fmt"
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
	logger *logger.LogrusLogger
}

// New creates a new application instance
func New() *App {
	// Load configuration
	conf := config.Loadconf()

	// Initialize structured logger
	appLogger := logger.SimpleLogger("app")

	// Initialize database
	config.DbInit()

	// Initialize server
	srv := server.New(conf)

	return &App{
		server: srv,
		config: conf,
		logger: appLogger,
	}
}

// Start starts the application server
func (a *App) Start() error {
	// Log startup information
	a.logger.Info("Starting application")

	// Start server in a goroutine
	go func() {
		address := ":" + a.config.Server.Port
		a.logger.Info("Server is up and running")

		if err := a.server.Start(address); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutdown signal received")

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
	a.logger.Info(fmt.Sprintf(
		"Configuration:\n"+
			"  Service: %s\n"+
			"  Port: %s\n"+
			"  Env: %s\n",
		a.config.Server.Name,
		a.config.Server.Port,
		a.config.Server.Environment,
	))
}
