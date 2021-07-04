package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/model/request"
	"github.com/ranggaaprilio/boilerGo/model/response"
)

func AddUser(c echo.Context) error {
	res := response.NewWebResponse()
	req := request.NewAdduser()

	if err := c.Bind(&req); err != nil {
		res.Code = http.StatusBadRequest
		res.Status = "Invalid Form Binding"
		res.Data = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, req)
}
