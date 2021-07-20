package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ranggaaprilio/boilerGo/app/v1/handler"
)

func WelcomeRouter(v1 *echo.Group) {
	v1.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome To Go RestApi V1")
	})
}

func UserRouter(v1 *echo.Group, userHandler *handler.UserHandler) {

	v1.POST("/adduser", userHandler.RegisterUser)
}
