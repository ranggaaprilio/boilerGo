package user

type AddUserForm struct {
	Name string `param:"name" query:"name" form:"name" json:"name" validate:"required"`
}

func NewAdduser() AddUserForm {
	return AddUserForm{}
}
