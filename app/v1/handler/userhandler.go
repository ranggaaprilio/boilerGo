/**
 * Package handler provides HTTP request handlers for the API endpoints.
 */
package handler

import (
	// "log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/app/v1/modules/user"
	"github.com/ranggaaprilio/boilerGo/helper"
)

/**
 * UserHandler handles HTTP requests related to user management.
 * It depends on the user service for business logic operations.
 */
type UserHandler struct {
	userService user.Service
}

// UserResponse represents a user for return in API responses
type UserResponse struct {
	ID        int     `json:"ID" example:"1"`
	Name      string  `json:"name" example:"John Doe"`
	CreatedAt string  `json:"CreatedAt" example:"2025-06-15T19:22:47.091+07:00"`
	UpdatedAt string  `json:"UpdatedAt" example:"2025-06-15T19:22:47.091+07:00"`
	DeletedAt *string `json:"DeletedAt,omitempty" example:"2025-06-15T19:22:47.091+07:00"`
}

/**
 * NewUserHandler creates a new instance of UserHandler with the provided user service.
 *
 * @param userService The service that handles user-related business logic
 * @return A pointer to a new UserHandler instance
 */
func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{userService}
}

/**
 * RegisterUser handles the user registration HTTP request.
 * It processes POST requests for creating new users.
 *
 * This method:
 * 1. Binds the request body to an AddUserForm struct
 * 2. Validates the form data
 * 3. Calls the user service to register the user
 * 4. Returns an appropriate response
 *
 * @param c Echo context containing the HTTP request and response
 * @return An error if one occurs during processing
 */

// @Summary Register a new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.AddUserForm true "User Data"
// @Success 200 {object} helper.SuccessResponse{data=UserResponse}
// @Failure 400 {object} helper.BadRequestResponse
// @Failure 500 {object} helper.InternalServerErrorResponse
// @Router /v1/users [post]
func (h *UserHandler) RegisterUser(c echo.Context) error {
	req := new(user.AddUserForm)
	var res helper.SuccessResponse
	var err error
	if err = c.Bind(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "Failed Form Binding"
		res.Data = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	if err = c.Validate(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "Name is Required"
		res.Data = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	// log.Fatal(err)

	newUser, err := h.userService.RegisterUser(req)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "Oops sorry ,Failed Save data"
		res.Data = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Code = http.StatusOK
	res.Message = "Success save data"
	res.Data = newUser
	return c.JSON(http.StatusOK, res)

}

/**
 * GetUser handles the HTTP request for retrieving a user by ID.
 * It processes GET requests for fetching user details.
 *
 * This method:
 * 1. Extracts the user ID from the request URL
 * 2. Calls the user service to fetch the user details
 * 3. Returns an appropriate response
 *
 * @param c Echo context containing the HTTP request and response
 * @return An error if one occurs during processing
 */

// @Summary Get a user by ID
// @Description Retrieves user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} helper.SuccessResponse{data=UserResponse}
// @Failure 400 {object} helper.BadRequestResponse
// @Failure 404 {object} helper.NotFoundResponse
// @Failure 500 {object} helper.InternalServerErrorResponse
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	var res helper.SuccessResponse

	// Example implementation - in a real app, this would call the user service
	// user, err := h.userService.GetUserByID(id)

	// Mock response for demonstration
	uid, err := strconv.Atoi(id)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "Invalid user ID"
		res.Data = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}
	userData := UserResponse{
		ID:        uid,
		Name:      "Example User",
		CreatedAt: "2025-06-15T19:22:47.091+07:00",
		UpdatedAt: "2025-06-15T19:22:47.091+07:00",
		DeletedAt: nil,
	}

	res.Code = http.StatusOK
	res.Message = "User found successfully"
	res.Data = userData
	return c.JSON(http.StatusOK, res)
}
