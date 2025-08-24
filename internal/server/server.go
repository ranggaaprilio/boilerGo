package server

import (
	"net/http"
	"runtime"

	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ranggaaprilio/boilerGo/config"
	"github.com/ranggaaprilio/boilerGo/internal/health"
	"github.com/ranggaaprilio/boilerGo/internal/logger"
	"github.com/ranggaaprilio/boilerGo/internal/server/middlewares"
	"github.com/ranggaaprilio/boilerGo/internal/server/routes"
)

// Server represents the HTTP server
type Server struct {
	echo   *echo.Echo
	config config.Configurations
	logger *logger.SlogLogger
}

// CustomValidator wraps the validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// New creates a new server instance
func New(conf config.Configurations) *echo.Echo {
	e := echo.New()

	// Setup custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup health checks
	healthService := health.NewHealthService()
	healthService.RegisterDefaultCheckers()

	// Setup middlewares
	setupMiddlewares(e, healthService)

	// Setup routes
	routes.SetupRoutes(e)

	return e
}

// setupMiddlewares configures all middlewares
func setupMiddlewares(e *echo.Echo, healthService *health.HealthService) {
	// Custom server header
	e.Use(middlewares.ServerHeader())

	// Stats middleware
	stats := middlewares.NewStats()
	e.Use(stats.Process)

	// Health check endpoints
	e.GET("/health", healthService.HealthHandler())
	e.GET("/health/live", healthService.LivenessHandler())
	e.GET("/health/ready", healthService.ReadinessHandler())
	e.GET("/healthcheck", stats.Handle) // Keep legacy endpoint

	// Gzip compression
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// Request ID for tracing
	e.Use(middleware.RequestID())

	// Request logging
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
	}))

	// Panic recovery
	e.Use(middleware.Recover())

	// CORS handling
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	}))

	// Static file serving
	e.Static("/", "public")

	// Root endpoint
	e.GET("/", func(c echo.Context) error {
		version := "The application was built with the Go version: " + runtime.Version()
		return c.String(http.StatusOK, version)
	})
}
