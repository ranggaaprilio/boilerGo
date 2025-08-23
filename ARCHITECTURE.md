# BoilerGo Architecture Documentation

## Overview

This document describes the improved architecture of the BoilerGo application, which follows Go best practices and provides a clean, maintainable, and scalable structure.

## Project Structure

```
boilerGo/
├── app/                          # Application domain logic (existing)
│   └── v1/
│       ├── handler/
│       └── modules/
├── config/                       # Configuration management
│   ├── appconf.go               # Application configuration
│   ├── db.go                    # Database configuration
│   └── testing/
├── internal/                     # Internal packages (not exported)
│   ├── app/                     # Application initialization
│   │   └── app.go
│   ├── health/                  # Health check system
│   │   └── health.go
│   ├── logger/                  # Structured logging
│   │   └── logger.go
│   └── server/                  # HTTP server setup
│       ├── middlewares/         # Custom middlewares
│       │   └── middlewares.go
│       ├── routes/              # Route definitions
│       │   ├── routes.go
│       │   ├── swagger.go
│       │   └── v1.go
│       └── server.go
├── docs/                        # Documentation
├── middleware/                  # External middlewares (existing)
├── public/                      # Static files
├── routes/                      # Legacy routes (to be deprecated)
├── scripts/                     # Utility scripts
├── test/                        # Tests
├── tmp/                         # Temporary files
├── bootstrap.go                 # Application bootstrap
├── main.go                      # Application entry point
└── [other config files...]
```

## Architecture Principles

### 1. Separation of Concerns
- Each package has a single, well-defined responsibility
- Configuration, logging, health checks, and server setup are separated
- Business logic is isolated from infrastructure concerns

### 2. Dependency Injection
- Dependencies are injected rather than created internally
- Makes testing easier and code more flexible
- Clear dependency graph

### 3. Layered Architecture
```
main.go
├── internal/app (Application Layer)
│   ├── internal/server (Presentation Layer)
│   ├── internal/health (Infrastructure Layer)
│   └── internal/logger (Infrastructure Layer)
├── app/v1 (Domain Layer)
└── config (Infrastructure Layer)
```

### 4. Clean Interfaces
- Small, focused interfaces
- Easy to mock for testing
- Clear contracts between components

## Core Components

### 1. Application (`internal/app/`)

**Purpose**: Main application lifecycle management

**Responsibilities**:
- Initialize all components
- Handle graceful shutdown
- Coordinate application startup
- Manage application-wide configuration

**Key Features**:
- Graceful shutdown with timeout
- Structured logging of startup/shutdown events
- Configuration validation and logging

### 2. Server (`internal/server/`)

**Purpose**: HTTP server setup and configuration

**Responsibilities**:
- Echo server initialization
- Middleware configuration
- Route registration
- Request/response handling setup

**Key Features**:
- Modular middleware setup
- Custom validator integration
- Health check endpoints
- Static file serving

### 3. Health Check System (`internal/health/`)

**Purpose**: Application health monitoring

**Responsibilities**:
- Database connectivity checks
- Memory usage monitoring
- Overall system health assessment
- Kubernetes-compatible endpoints

**Endpoints**:
- `/health` - Comprehensive health check
- `/health/live` - Liveness probe
- `/health/ready` - Readiness probe
- `/healthcheck` - Legacy stats endpoint

### 4. Structured Logging (`internal/logger/`)

**Purpose**: Centralized logging with structure

**Responsibilities**:
- Environment-appropriate log formatting
- Contextual logging
- Performance logging
- Security event logging

**Features**:
- JSON logging for production
- Text logging for development
- Request tracing
- Performance monitoring
- Security event logging

### 5. Configuration (`config/`)

**Purpose**: Application configuration management

**Responsibilities**:
- Environment variable binding
- Configuration validation
- Default value management
- Database connection management

**Features**:
- Viper-based configuration
- Environment variable override
- Configuration validation
- Multiple configuration sources

### 6. Bootstrap (`bootstrap.go`)

**Purpose**: Application initialization and database setup

**Responsibilities**:
- Database migrations
- Initial data seeding
- System preparation
- Dependency verification

## Request Flow

```
1. HTTP Request → Echo Server
2. Middlewares (CORS, Logging, Recovery, etc.)
3. Route Handler → Domain Logic (app/v1/)
4. Repository → Database
5. Response ← Handler ← Domain Logic
6. Response ← Middlewares ← Echo Server
```

## Configuration Management

### Configuration Sources (Priority Order)
1. Environment Variables
2. Configuration File (`config.yml`)
3. Default Values

### Environment Variables
```bash
# Server Configuration
SERVER_NAME=BoilerGo
SERVER_PORT=8080
ENVIRONMENT=development

# Database Configuration
DB_USER=username
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=database

# Application Configuration
LOG_LEVEL=info
DEBUG=false
SECRET_KEY=your-secret-key
```

### Configuration Structure
```go
type Configurations struct {
    Server   ServerConfigurations
    Database DbConfigurations
    App      AppConfigurations
}
```

## Health Checks

### Health Check Types

1. **Liveness Probe** (`/health/live`)
   - Checks if the application is running
   - Returns 200 if alive, 503 if not

2. **Readiness Probe** (`/health/ready`)
   - Checks if the application is ready to serve requests
   - Validates database connectivity

3. **Comprehensive Health** (`/health`)
   - Runs all registered health checkers
   - Returns detailed status information
   - Includes timing and error details

### Health Check Response Format
```json
{
  "status": "healthy|unhealthy|degraded",
  "timestamp": "2023-01-01T00:00:00Z",
  "uptime": "1h30m45s",
  "version": "1.0.0",
  "service": "BoilerGo",
  "checks": [
    {
      "name": "database",
      "status": "healthy",
      "message": "Database connection is healthy",
      "last_checked": "2023-01-01T00:00:00Z",
      "duration": "5ms"
    }
  ]
}
```

## Logging Strategy

### Log Levels
- **Debug**: Detailed information for debugging
- **Info**: General information about application flow
- **Warn**: Warning messages for potentially harmful situations
- **Error**: Error events that might still allow the application to continue

### Log Formats
- **Development**: Human-readable text format
- **Production**: Structured JSON format for log aggregation

### Contextual Logging
```go
logger := appLogger.WithComponent("user-service")
logger = logger.WithRequestID(requestID)
logger.Info("User created", "user_id", userID)
```

## Middleware Stack

1. **ServerHeader** - Custom server identification
2. **Stats** - Request statistics collection
3. **Gzip** - Response compression
4. **RequestID** - Unique request identification
5. **Logger** - Request/response logging
6. **Recover** - Panic recovery
7. **CORS** - Cross-origin resource sharing

## Error Handling

### Exception Handling
- Centralized error handling through `exception` package
- Panic recovery middleware
- Structured error responses

### Error Response Format
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {},
    "timestamp": "2023-01-01T00:00:00Z"
  }
}
```

## Database Management

### Connection Management
- Connection pooling with configurable limits
- Automatic retry logic for connection failures
- Health monitoring of database connections

### Migration System
- Automatic migrations on application startup
- Version-controlled schema changes
- Rollback capabilities

## Security Considerations

### Security Features
- CORS configuration
- Request validation
- Secure headers
- Input sanitization

### Logging Security Events
```go
logger.LogSecurity("failed_login", userID, clientIP, details)
```

## Performance Monitoring

### Metrics Collection
- Request duration tracking
- Database query performance
- Memory usage monitoring
- Health check timing

### Performance Logging
```go
logger.LogPerformance("database_query", duration, "query", query)
```

## Deployment Considerations

### Docker Support
- Multi-stage Docker builds
- Health check integration
- Environment variable configuration

### Kubernetes Deployment
- Liveness and readiness probes
- ConfigMap integration
- Service discovery support

## Migration from Legacy Structure

### What Changed
1. **routes/routes.go** → Split into multiple focused files
2. **main.go** → Simplified with proper application structure
3. **bootstrap.go** → Enhanced with proper error handling
4. **Configuration** → Improved validation and environment support
5. **Logging** → Structured logging with context

### Migration Steps
1. Update import paths to use new internal packages
2. Replace direct logger usage with structured logger
3. Update health check endpoints in monitoring systems
4. Verify environment variable configurations
5. Test graceful shutdown behavior

## Best Practices

### Code Organization
- Use internal packages for non-exported functionality
- Keep domain logic separate from infrastructure
- Use dependency injection for testability

### Configuration
- Always provide default values
- Validate configuration on startup
- Use environment variables for deployment-specific values

### Logging
- Use structured logging with context
- Log performance metrics
- Include request tracing information

### Error Handling
- Handle errors at the appropriate level
- Provide meaningful error messages
- Log errors with sufficient context

### Testing
- Use interfaces for easier mocking
- Test configuration validation
- Test health check endpoints
- Test graceful shutdown

## Future Enhancements

### Potential Improvements
1. **Metrics Collection** - Prometheus integration
2. **Distributed Tracing** - OpenTelemetry support
3. **Rate Limiting** - Request rate limiting middleware
4. **Circuit Breaker** - Database connection resilience
5. **Caching** - Redis integration for performance
6. **Authentication** - JWT middleware
7. **API Versioning** - Enhanced version management

### Scalability Considerations
- Horizontal scaling support
- Database connection pooling optimization
- Caching strategy implementation
- Load balancing configuration

## Conclusion

The new architecture provides:
- **Better Maintainability** - Clear separation of concerns
- **Improved Testability** - Dependency injection and interfaces
- **Enhanced Observability** - Structured logging and health checks
- **Production Readiness** - Graceful shutdown and error handling
- **Scalability** - Modular design for easy extension

This structure follows Go best practices and industry standards, making the codebase more maintainable, testable, and production-ready.