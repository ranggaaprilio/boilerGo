// Package helper provides utility functions and structures for the application
package helper

// BadRequestResponse represents a standardized error response for bad requests
type BadRequestResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
	Data    string `json:"data,omitempty"`
}

// InternalServerErrorResponse represents a standardized error response for internal server errors
type InternalServerErrorResponse struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
	Data    string `json:"data,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data,omitempty"`
}

// NotFoundResponse represents a standardized error response for not found errors
type NotFoundResponse struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"Not Found"`
	Data    string `json:"data,omitempty"`
}
