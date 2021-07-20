package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	// "github.com/dgrijalva/jwt-go"
	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ranggaaprilio/boilerGo/app/v1/handler"
	"github.com/ranggaaprilio/boilerGo/app/v1/modules/user"
	"github.com/ranggaaprilio/boilerGo/config"
	"github.com/ranggaaprilio/boilerGo/exception"
)

//ServerHeader Config
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Aprilio/1.0")
		return next(c)
	}

}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

//Init Routing Initialize
func Init() *echo.Echo {
	db := config.CreateCon()

	e := echo.New()
	/** custom Header **/
	e.Use(ServerHeader)

	/** middeleware **/
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} ,time=${time_rfc3339}\n",
	}))
	e.Use(middleware.Recover())
	/** middeleware **/

	//routing
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "GO Version 1.16 ")
	})

	/**v1 Group==============================================================**/
	v1 := e.Group("/api/v1")

	//**Welcome v1 api**//
	WelcomeRouter(v1)
	//**end Welcome v1 api**//

	//**User v1 api**//
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	UserRouter(v1, userHandler)
	//**User v1 api**//

	/**end v1 Group============================================================**/

	//Mapping To Json file
	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		exception.PanicIfNeeded(err)
	}
	ioutil.WriteFile("routes/routes.json", data, 0644)

	return e
}
