package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// ServerHeader sets custom server header
func ServerHeader() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "BoilerGo/1.0")
			return next(c)
		}
	}
}

// Stats represents server statistics
type Stats struct {
	Uptime       time.Time      `json:"uptime"`
	RequestCount uint64         `json:"requestCount"`
	Statuses     map[string]int `json:"statuses"`
	mutex        sync.RWMutex
}

// NewStats creates a new Stats instance
func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: make(map[string]int),
	}
}

// Process is the middleware function for collecting stats
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := c.Response().Status
		s.Statuses[http.StatusText(status)]++
		return nil
	}
}

// Handle returns the stats as JSON
func (s *Stats) Handle(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := map[string]interface{}{
		"uptime":       time.Since(s.Uptime).String(),
		"requestCount": s.RequestCount,
		"statuses":     s.Statuses,
	}

	return c.JSON(200, stats)
}
