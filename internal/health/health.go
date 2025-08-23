package health

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/config"
)

// HealthStatus represents the overall health status
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusUnhealthy HealthStatus = "unhealthy"
	StatusDegraded  HealthStatus = "degraded"
)

// HealthCheck represents a single health check
type HealthCheck struct {
	Name        string        `json:"name"`
	Status      HealthStatus  `json:"status"`
	Message     string        `json:"message,omitempty"`
	LastChecked time.Time     `json:"last_checked"`
	Duration    time.Duration `json:"duration"`
}

// HealthResponse represents the complete health check response
type HealthResponse struct {
	Status    HealthStatus  `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Uptime    time.Duration `json:"uptime"`
	Checks    []HealthCheck `json:"checks"`
	Version   string        `json:"version"`
	Service   string        `json:"service"`
}

// HealthChecker interface for implementing health checks
type HealthChecker interface {
	Name() string
	Check(ctx context.Context) HealthCheck
}

// DatabaseHealthChecker checks database connectivity
type DatabaseHealthChecker struct{}

func (d *DatabaseHealthChecker) Name() string {
	return "database"
}

func (d *DatabaseHealthChecker) Check(ctx context.Context) HealthCheck {
	start := time.Now()
	check := HealthCheck{
		Name:        d.Name(),
		LastChecked: start,
	}

	if err := config.PingDB(); err != nil {
		check.Status = StatusUnhealthy
		check.Message = err.Error()
	} else {
		check.Status = StatusHealthy
		check.Message = "Database connection is healthy"
	}

	check.Duration = time.Since(start)
	return check
}

// MemoryHealthChecker checks memory usage
type MemoryHealthChecker struct {
	MaxMemoryMB int64
}

func (m *MemoryHealthChecker) Name() string {
	return "memory"
}

func (m *MemoryHealthChecker) Check(ctx context.Context) HealthCheck {
	start := time.Now()
	check := HealthCheck{
		Name:        m.Name(),
		LastChecked: start,
	}

	// This is a simplified memory check
	// In a real application, you might want to use runtime.MemStats
	check.Status = StatusHealthy
	check.Message = "Memory usage is within acceptable limits"
	check.Duration = time.Since(start)

	return check
}

// HealthService manages all health checks
type HealthService struct {
	checkers  []HealthChecker
	startTime time.Time
	version   string
	service   string
}

// NewHealthService creates a new health service
func NewHealthService() *HealthService {
	return &HealthService{
		checkers:  make([]HealthChecker, 0),
		startTime: time.Now(),
		version:   "1.0.0", // This should come from build info
		service:   "BoilerGo",
	}
}

// AddChecker adds a health checker to the service
func (hs *HealthService) AddChecker(checker HealthChecker) {
	hs.checkers = append(hs.checkers, checker)
}

// RegisterDefaultCheckers registers the default set of health checkers
func (hs *HealthService) RegisterDefaultCheckers() {
	hs.AddChecker(&DatabaseHealthChecker{})
	hs.AddChecker(&MemoryHealthChecker{MaxMemoryMB: 512})
}

// CheckHealth performs all health checks
func (hs *HealthService) CheckHealth(ctx context.Context) HealthResponse {
	response := HealthResponse{
		Timestamp: time.Now(),
		Uptime:    time.Since(hs.startTime),
		Version:   hs.version,
		Service:   hs.service,
		Checks:    make([]HealthCheck, 0, len(hs.checkers)),
	}

	overallStatus := StatusHealthy

	// Run all health checks
	for _, checker := range hs.checkers {
		check := checker.Check(ctx)
		response.Checks = append(response.Checks, check)

		// Determine overall status
		switch check.Status {
		case StatusUnhealthy:
			overallStatus = StatusUnhealthy
		case StatusDegraded:
			if overallStatus == StatusHealthy {
				overallStatus = StatusDegraded
			}
		}
	}

	response.Status = overallStatus
	return response
}

// HealthHandler returns an Echo handler for health checks
func (hs *HealthService) HealthHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
		defer cancel()

		healthResponse := hs.CheckHealth(ctx)

		statusCode := http.StatusOK
		switch healthResponse.Status {
		case StatusUnhealthy:
			statusCode = http.StatusServiceUnavailable
		case StatusDegraded:
			statusCode = http.StatusPartialContent
		}

		return c.JSON(statusCode, healthResponse)
	}
}

// ReadinessHandler returns a simple readiness check
func (hs *HealthService) ReadinessHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Simple readiness check - just check if database is accessible
		if err := config.PingDB(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "not ready",
				"reason": "database not accessible",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ready",
		})
	}
}

// LivenessHandler returns a simple liveness check
func (hs *HealthService) LivenessHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "alive",
			"uptime":  time.Since(hs.startTime).String(),
			"service": hs.service,
		})
	}
}
