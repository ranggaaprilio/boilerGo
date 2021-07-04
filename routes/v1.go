package routes

import (
	"github.com/labstack/echo/v4"
	ver1 "github.com/ranggaaprilio/boilerGo/controller/v1"
)

func UserRouter(v1 *echo.Group) {
	v1.POST("/adduser", ver1.AddUser)
}

func WelcomeRouter(v1 *echo.Group) {
	v1.GET("/", ver1.Welcome)
}
