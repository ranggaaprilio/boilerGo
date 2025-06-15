# Swagger Documentation for BoilerGo API

This folder contains the Swagger/OpenAPI documentation for the BoilerGo API.

## How to View Documentation

The Swagger UI is available at:

```
http://localhost:8080/swagger/index.html
```

## Files

- `docs.go`: Generated Swagger specification file
- `swagger.json`: Generated Swagger JSON specification
- `swagger.yaml`: Generated Swagger YAML specification

## Updating Documentation

To update the Swagger documentation after making code changes:

```bash
# Using the provided script
./scripts/generate-swagger.sh

# Or using make
make swagger
```

## Adding Documentation to Endpoints

Use Swagger annotations in your handler functions to document API endpoints. Example:

```go
// @Summary Get a user by ID
// @Description Retrieves user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} helper.WebResponse
// @Failure 404 {object} helper.WebResponse
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
    // Implementation...
}
```

## Common Annotation Types

- `@Summary`: Brief description of what the endpoint does
- `@Description`: More detailed explanation of the endpoint
- `@Tags`: Grouping for the endpoint in Swagger UI
- `@Accept`: Content types the API can consume (e.g., json, xml)
- `@Produce`: Content types the API can produce (e.g., json, xml)
- `@Param`: Parameters the endpoint accepts
- `@Success`: Successful response definition
- `@Failure`: Error response definition
- `@Router`: The path and HTTP method for the endpoint

## More Information

For detailed information on Swagger annotations, visit:
[Swaggo Documentation](https://github.com/swaggo/swag)
