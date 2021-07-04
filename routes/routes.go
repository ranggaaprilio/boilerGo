package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	// "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ranggaaprilio/boilerGo/exception"
)

//ServerHeader Config
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Aprilio/1.0")
		return next(c)
	}

}

//Init Routing Initialize
func Init() *echo.Echo {
	e := echo.New()
	/** custom Header **/
	e.Use(ServerHeader)

	/** middeleware **/
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} ,time=${time_rfc3339}\n",
	}))
	e.Use(middleware.Recover())
	/** middeleware **/

	//routing
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "GO Version 1.16 ")
	})

	/**v1 Group**/
	v1 := e.Group("/api/v1")
	WelcomeRouter(v1)
	UserRouter(v1)
	/**v1 Group**/

	//Mapping To Json file
	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		exception.PanicIfNeeded(err)
	}
	ioutil.WriteFile("routes/routes.json", data, 0644)

	return e
}
