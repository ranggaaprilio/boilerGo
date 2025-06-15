# User API Documentation

This document provides detailed information about the User API endpoints in the BoilerGo application.

## Endpoints

### Register User

Registers a new user in the system.

**URL**: `/api/v1/users`

**Method**: `POST`

**Request Body**:

```json
{
  "name": "John Doe"
}
```

**Request Parameters**:

| Parameter | Type   | Required | Description     |
| --------- | ------ | -------- | --------------- |
| name      | string | Yes      | The user's name |

**Response**:

- Success (200 OK)

```json
{
  "code": 200,
  "message": "Success save data",
  "data": {
    "ID": 1,
    "CreatedAt": "2025-06-15T10:00:00Z",
    "UpdatedAt": "2025-06-15T10:00:00Z",
    "DeletedAt": null,
    "Name": "John Doe"
  }
}
```

- Form Binding Error (400 Bad Request)

```json
{
  "code": 400,
  "message": "Failed Form Binding",
  "data": "error message"
}
```

- Validation Error (400 Bad Request)

```json
{
  "code": 400,
  "message": "Name is Required",
  "data": "error message"
}
```

- Server Error (500 Internal Server Error)

```json
{
  "code": 500,
  "message": "Oops sorry, Failed Save data",
  "data": "error message"
}
```

## Implementation Details

### Handler

The `UserHandler` struct handles HTTP requests related to user management.

```go
type UserHandler struct {
    userService user.Service
}
```

#### RegisterUser

`RegisterUser` is a method that handles user registration requests.

```go
func (h *UserHandler) RegisterUser(c echo.Context) error
```

This method:

1. Binds the request body to the `AddUserForm` struct
2. Validates the form data
3. Calls the `RegisterUser` method on the user service
4. Returns an appropriate response
