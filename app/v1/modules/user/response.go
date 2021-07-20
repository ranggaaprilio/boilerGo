package user

type resAdduser struct {
	Name string `param:"name" query:"name" form:"name" json:"name"`
}

func NewResAdduser() resAdduser {
	return resAdduser{}
}
