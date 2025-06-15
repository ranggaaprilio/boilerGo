# OpenAPI Documentation

This directory contains OpenAPI (Swagger) documentation for the BoilerGo API.

## Files

- `swagger.go`: Main Swagger configuration file
- `docs.go`: Generated Swagger specification

## Usage

The Swagger UI is available at:

```
http://localhost:8080/swagger/index.html
```

## Generating Documentation

To regenerate the Swagger documentation after making changes to your API:

```bash
./scripts/generate-swagger.sh
```

This will scan your code for Swagger annotations and update the documentation files.

## Annotations

To document your API endpoints, use Swagger annotations in your handler functions. For example:

```go
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.AddUserForm true "User Data"
// @Success 200 {object} helper.WebResponse
// @Failure 400 {object} helper.WebResponse
// @Failure 500 {object} helper.WebResponse
// @Router /v1/users [post]
func (h *UserHandler) RegisterUser(c echo.Context) error {
    // Implementation
}
```
