# Code Architecture Documentation

This document outlines the architecture and design patterns used in the BoilerGo application.

## Overview

The application follows a clean architecture pattern with clear separation of concerns. The codebase is organized into the following layers:

- **Handler Layer**: Handles HTTP requests and responses
- **Service Layer**: Contains business logic
- **Repository Layer**: Manages data access
- **Entity**: Defines data structures

## Logging Architecture

The application uses a dual logging system with both slog and logrus:

- **Slog Logger**: Used for main application logging with structured JSON/text output
- **Logrus Logger**: Used for component-specific logging (bootstrap, config, database)
- **Component Loggers**: Simple logrus loggers for specific components without config dependencies

### Logger Components

The logging system includes:

```go
// Simple logrus logger for components
func SimpleLogger(component string) *LogrusLogger {
    logger := logrus.New()
    logger.SetLevel(logrus.InfoLevel)
    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
        ForceColors:   true,
    })
    logger.SetOutput(os.Stdout)

    return &LogrusLogger{
        Logger:    logger.WithField("component", component).Logger,
        component: component,
    }
}
```

### GORM Integration

Database operations use logrus through a custom writer:

```go
type logrusGormWriter struct {
    logger *LogrusLogger
}

func (w *logrusGormWriter) Write(p []byte) (n int, err error) {
    w.logger.Debug(string(p))
    return len(p), nil
}

func (w *logrusGormWriter) Printf(format string, args ...interface{}) {
    w.logger.Debug(fmt.Sprintf(format, args...))
}
```

## User Module

### Entity

The `User` entity represents a user in the system:

```go
type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(250)" `
}
```

The `gorm.Model` embedding provides the standard ID, CreatedAt, UpdatedAt, and DeletedAt fields.

### Request/DTO

The `AddUserForm` defines the structure for user creation requests:

```go
type AddUserForm struct {
	Name string `param:"name" query:"name" form:"name" json:"name" validate:"required"`
}
```

This struct is used for binding and validating request data.

### Repository Layer

The `Repository` interface defines methods to interact with the user data:

```go
type Repository interface {
	Save(user User) (User, error)
}
```

Implementation:

```go
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
```

### Service Layer

The `Service` interface defines the business operations for users:

```go
type Service interface {
	RegisterUser(input *AddUserForm) (User, error)
}
```

Implementation:

```go
type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input *AddUserForm) (User, error) {
	user := User{}
	user.Name = input.Name

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}
```

### Handler Layer

The `UserHandler` struct handles HTTP requests:

```go
type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{userService}
}
```

The `RegisterUser` method processes user registration:

```go
func (h *UserHandler) RegisterUser(c echo.Context) error {
	req := new(user.AddUserForm)
	var res helper.WebResponse

	// Bind request data
	if err = c.Bind(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "Failed Form Binding"
		res.Data = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	// Validate request
	if err = c.Validate(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "Name is Required"
		res.Data = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	// Process through service layer
	newUser, err := h.userService.RegisterUser(req)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "Oops sorry, Failed Save data"
		res.Data = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	// Return success response
	res.Code = http.StatusOK
	res.Message = "Success save data"
	res.Data = newUser
	return c.JSON(http.StatusOK, res)
}
```

## Configuration System

The configuration system uses logrus for logging during the loading process:

```go
type ConfigLoader struct {
    logger *LogrusLogger
}

func NewConfigLoader() *ConfigLoader {
    return &ConfigLoader{
        logger: SimpleLogger("config"),
    }
}
```

Configuration loading includes structured logging:

```go
func (cl *ConfigLoader) Load() Configurations {
    cl.logger.Info("Loading application configuration...")
    
    // Configuration loading logic...
    
    cl.logger.Info("Configuration loaded and validated successfully")
    return configuration
}
```

## Database Layer

The database layer integrates logrus logging for all database operations:

```go
func DbInit() {
    conf := Loadconf()
    dbConfig := DefaultDatabaseConfig()

    // Initialize logger for database operations
    dbLogger := SimpleLogger("database")

    // Log connection attempt
    dbLogger.Info("Attempting to connect to database",
        "host", conf.Database.DbHost,
        "port", conf.Database.DbPort,
        "user", conf.Database.DbUsername)

    // GORM configuration with logrus writer
    logrusWriter := &logrusGormWriter{logger: dbLogger}
    gormConfig := &gorm.Config{
        Logger: logger.New(
            logrusWriter,
            logger.Config{
                SlowThreshold:             200 * time.Millisecond,
                LogLevel:                  logger.Warn,
                IgnoreRecordNotFoundError: true,
                Colorful:                  false,
            },
        ),
    }
}
```

## Bootstrap Process

The bootstrap process uses logrus for initialization logging:

```go
func Bootstrap() error {
    // Initialize simple logger for bootstrap
    bootstrapLogger := SimpleLogger("bootstrap")

    bootstrapLogger.Info("Starting bootstrap process...")

    // Database migrations with logging
    bootstrapLogger.Info("Running database migrations...")
    if err := db.AutoMigrate(&user.User{}); err != nil {
        bootstrapLogger.Error("Failed to migrate User model", "error", err)
        return err
    }

    bootstrapLogger.Info("Bootstrap process completed successfully")
    return nil
}
```

## Helper Functions

The application uses standardized response formats:

```go
type WebResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
```

This ensures consistent API responses throughout the application.

### Logging in Handlers

Handlers can utilize both the main application logger and component-specific loggers:

```go
func (h *UserHandler) RegisterUser(c echo.Context) error {
    // Main application logger for request processing
    // Component logger for specific operations
    
    // Business logic with appropriate logging
    newUser, err := h.userService.RegisterUser(req)
    if err != nil {
        // Error logging with context
        return c.JSON(http.StatusBadRequest, res)
    }
    
    return c.JSON(http.StatusOK, res)
}
```

## Server Architecture

The server layer provides HTTP handling using Echo framework with comprehensive middleware stack:

```go
// Server represents the HTTP server
type Server struct {
    echo   *echo.Echo
    config config.Configurations
    logger *logger.SlogLogger
}

// New creates a new server instance
func New(conf config.Configurations, appLogger *logger.SlogLogger) *echo.Echo {
    e := echo.New()

    // Setup custom validator
    e.Validator = &CustomValidator{validator: validator.New()}

    // Setup health checks
    healthService := health.NewHealthService()
    healthService.RegisterDefaultCheckers()

    // Setup middlewares
    setupMiddlewares(e, appLogger, healthService)

    // Setup routes
    routes.SetupRoutes(e)

    return e
}
```

### Middleware Stack

The application includes comprehensive middleware configuration:

```go
func setupMiddlewares(e *echo.Echo, appLogger *logger.SlogLogger, healthService *health.HealthService) {
    // Custom server header
    e.Use(middlewares.ServerHeader())

    // Request/Response statistics
    stats := middlewares.NewStats()
    e.Use(stats.Process)

    // Health check endpoints
    e.GET("/health", healthService.HealthHandler())
    e.GET("/health/live", healthService.LivenessHandler())
    e.GET("/health/ready", healthService.ReadinessHandler())

    // Compression middleware
    e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
        Level: 5,
    }))

    // Request ID for distributed tracing
    e.Use(middleware.RequestID())

    // HTTP request logging
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
    }))

    // Panic recovery
    e.Use(middleware.Recover())

    // CORS configuration
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
}
```

## Health Check System

The application includes comprehensive health monitoring:

```go
type HealthService struct {
    checkers map[string]HealthChecker
    mu       sync.RWMutex
}

// Health check endpoints
// GET /health - Overall system health
// GET /health/live - Liveness probe (for Kubernetes)
// GET /health/ready - Readiness probe (for Kubernetes)
```

Health checks include:
- Database connectivity
- External service dependencies
- System resource availability
- Application readiness status

## Application Layer

The main application layer orchestrates all components:

```go
// App represents the main application structure
type App struct {
    server *echo.Echo
    config config.Configurations
    logger *logger.SlogLogger
}

func New() *App {
    // Load configuration
    conf := config.Loadconf()

    // Initialize structured logger
    appLogger := logger.NewSlog(
        conf.App.LogLevel,
        conf.Server.Environment,
        conf.App.ServiceName,
        conf.IsDevelopment(),
    )

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
```

### Graceful Shutdown

The application supports graceful shutdown with proper cleanup:

```go
func (a *App) Start() error {
    // Start server in goroutine
    go func() {
        address := ":" + a.config.Server.Port
        if err := a.server.Start(address); err != nil && err != http.ErrServerClosed {
            a.logger.Fatal("Failed to start server", "error", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    a.logger.LogShutdown("signal received")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := a.server.Shutdown(ctx); err != nil {
        a.logger.Error("Server shutdown error", "error", err)
        return err
    }

    a.logger.Info("Server stopped gracefully")
    return nil
}
```

## Route Organization

Routes are organized by version and functionality:

```
/api/v1/user/register - User registration endpoint
/health             - Health check endpoint
/health/live        - Liveness probe
/health/ready       - Readiness probe
/docs/*             - Swagger documentation
```

The routing system supports:
- Version-based API organization
- Middleware application per route group
- Swagger documentation integration
- RESTful endpoint patterns

## Error Handling & Validation

The application implements comprehensive error handling and validation patterns:

### Custom Validator

```go
type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
    if err := cv.validator.Struct(i); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    return nil
}
```

### Request Validation

Input validation is performed at the handler level:

```go
func (h *UserHandler) RegisterUser(c echo.Context) error {
    req := new(user.AddUserForm)
    var res helper.WebResponse

    // Bind request data
    if err = c.Bind(req); err != nil {
        res.Code = http.StatusBadRequest
        res.Message = "Failed Form Binding"
        res.Data = err.Error()
        return c.JSON(http.StatusBadRequest, res)
    }

    // Validate request structure
    if err = c.Validate(req); err != nil {
        res.Code = http.StatusBadRequest
        res.Message = "Name is Required"
        res.Data = err.Error()
        return c.JSON(http.StatusBadRequest, res)
    }
}
```

### Error Response Standardization

All API responses follow a consistent format:

```go
type WebResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
```

### Exception Handling

The application includes centralized exception handling:

```go
// Global panic recovery
func Catch() {
    if r := recover(); r != nil {
        // Log panic with context
        // Perform cleanup operations
        // Exit gracefully
    }
}

// Usage in main function
func main() {
    defer exception.Catch()
    // Application logic
}
```

### Validation Tags

Request structures use validation tags for automatic validation:

```go
type AddUserForm struct {
    Name string `param:"name" query:"name" form:"name" json:"name" validate:"required"`
}
```

Supported validation rules:
- `required` - Field must be present
- `email` - Valid email format
- `min`, `max` - String/number length constraints
- `uuid` - Valid UUID format
- Custom validation functions

## Swagger Documentation

The application includes comprehensive API documentation:

```go
// @title BoilerGo API
// @version 1.0
// @description This is the API documentation for the BoilerGo application.
// @contact.name API Support
// @contact.url https://github.com/ranggaaprilio/boilerGo
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api
// @schemes http https
```

### API Endpoint Documentation

Each endpoint includes detailed Swagger annotations:

```go
// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param user body AddUserForm true "User registration data"
// @Success 200 {object} WebResponse
// @Failure 400 {object} WebResponse
// @Router /api/v1/user/register [post]
func (h *UserHandler) RegisterUser(c echo.Context) error {
    // Handler implementation
}
```

## Security Considerations

The application implements several security measures:

### CORS Configuration

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"}, // Configure for production
    AllowHeaders: []string{
        echo.HeaderOrigin,
        echo.HeaderContentType,
        echo.HeaderAccept,
        echo.HeaderAuthorization,
    },
}))
```

### Request ID Tracing

Each request receives a unique identifier for tracing:

```go
e.Use(middleware.RequestID())
```

### Panic Recovery

Automatic recovery from panics with logging:

```go
e.Use(middleware.Recover())
```

## Performance Optimizations

### Gzip Compression

Response compression is enabled by default:

```go
e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
    Level: 5,
}))
```

### Database Connection Pooling

Optimized database connections:

```go
if sqlDB, poolErr := db.DB(); poolErr == nil {
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
}
```

### Request Statistics

Built-in request monitoring and statistics:

```go
stats := middlewares.NewStats()
e.Use(stats.Process)
e.GET("/healthcheck", stats.Handle)
```

## Environment Configuration

The application supports multiple deployment environments with flexible configuration:

### Configuration Sources

Configuration is loaded from multiple sources in order of precedence:

1. **Environment Variables** (highest priority)
2. **Configuration File** (`config.yml`)
3. **Default Values** (lowest priority)

```go
// Environment variable mappings
envMappings := map[string]string{
    "server.name":         "SERVER_NAME",
    "server.port":         "SERVER_PORT",
    "server.environment":  "ENVIRONMENT",
    "database.dbusername": "DB_USER",
    "database.dbpassword": "DB_PASSWORD",
    "database.dbhost":     "DB_HOST",
    "database.dbport":     "DB_PORT",
    "database.dbname":     "DB_NAME",
    "app.log_level":       "LOG_LEVEL",
    "app.debug":           "DEBUG",
}
```

### Environment-Specific Behavior

The application adapts its behavior based on the environment:

```go
// Development environment
func (c *Configurations) IsDevelopment() bool {
    return c.Server.Environment == "development"
}

// Production environment
func (c *Configurations) IsProduction() bool {
    return c.Server.Environment == "production"
}
```

**Development Mode:**
- Colored console logging
- Detailed error messages
- Source code line numbers in logs
- Debug-level logging enabled

**Production Mode:**
- JSON structured logging
- Minimal error exposure
- Performance optimizations
- Info-level logging by default

### Configuration Structure

```go
type Configurations struct {
    Server   ServerConfigurations `mapstructure:"server"`
    Database DbConfigurations     `mapstructure:"database"`
    App      AppConfigurations    `mapstructure:"app"`
}

type ServerConfigurations struct {
    Name         string `mapstructure:"name"`
    Port         string `mapstructure:"port"`
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
    Environment  string `mapstructure:"environment"`
}
```

## Deployment Strategies

### Docker Deployment

The application is designed for containerized deployment:

```dockerfile
# Example Dockerfile structure
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config.yml .
CMD ["./main"]
```

### Kubernetes Deployment

Health check endpoints support Kubernetes probes:

```yaml
# Kubernetes deployment example
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: boilergo
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Environment Variables for Production

```bash
# Server configuration
SERVER_NAME=BoilerGo-Production
SERVER_PORT=8080
ENVIRONMENT=production

# Database configuration
DB_HOST=db.example.com
DB_PORT=3306
DB_USER=app_user
DB_PASSWORD=secure_password
DB_NAME=boilergo_prod

# Application configuration
LOG_LEVEL=info
DEBUG=false
SECRET_KEY=your-secret-key
```

## Monitoring and Observability

### Structured Logging

All logs include structured data for monitoring:

```json
{
  "level": "info",
  "time": "2024-01-01T12:00:00Z",
  "service": "BoilerGo",
  "component": "database",
  "environment": "production",
  "message": "Database connection established",
  "host": "db.example.com",
  "port": "3306"
}
```

### Metrics Endpoints

- `/health` - Overall application health
- `/health/live` - Liveness probe for container orchestration
- `/health/ready` - Readiness probe for load balancers
- `/healthcheck` - Request statistics and performance metrics

### Request Tracing

Each request includes tracing information:

```go
// Request ID middleware adds unique identifier
e.Use(middleware.RequestID())

// Accessible in handlers via context
requestID := c.Response().Header().Get(echo.HeaderXRequestID)
```

## Best Practices

### Error Handling
- Always return structured error responses
- Log errors with appropriate context
- Use proper HTTP status codes
- Never expose internal errors to clients

### Performance
- Enable gzip compression for responses
- Use database connection pooling
- Implement proper timeout configurations
- Monitor slow database queries

### Security
- Validate all input data
- Use environment variables for secrets
- Configure CORS appropriately for production
- Implement proper authentication/authorization
- Log security-related events

### Logging
- Use structured logging with consistent fields
- Include request IDs for tracing
- Log at appropriate levels (debug, info, warn, error)
- Avoid logging sensitive information
