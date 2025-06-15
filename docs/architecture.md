# Code Architecture Documentation

This document outlines the architecture and design patterns used in the BoilerGo application.

## Overview

The application follows a clean architecture pattern with clear separation of concerns. The codebase is organized into the following layers:

- **Handler Layer**: Handles HTTP requests and responses
- **Service Layer**: Contains business logic
- **Repository Layer**: Manages data access
- **Entity**: Defines data structures

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
