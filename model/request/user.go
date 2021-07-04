package request

type AddUser struct {
	Name string `param:"name" query:"name" form:"name" json:"name"`
}

func NewAdduser() AddUser {
	return AddUser{}
}
