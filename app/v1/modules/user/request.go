package user

// AddUserForm represents the request data structure for adding a new user
// @Description User registration request form
type AddUserForm struct {
	Name string `param:"name" query:"name" form:"name" json:"name" validate:"required" example:"John Doe"`
}

// NewAdduser creates a new instance of AddUserForm
func NewAdduser() AddUserForm {
	return AddUserForm{}
}
