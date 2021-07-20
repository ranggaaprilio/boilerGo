package handler

import (
	// "log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/app/v1/modules/user"
	"github.com/ranggaaprilio/boilerGo/helper"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	req := new(user.AddUserForm)
	var res helper.WebResponse
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
