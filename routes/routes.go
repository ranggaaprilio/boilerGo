package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"

	// "github.com/dgrijalva/jwt-go"
	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ranggaaprilio/boilerGo/app/v1/handler"
	"github.com/ranggaaprilio/boilerGo/app/v1/modules/user"
	"github.com/ranggaaprilio/boilerGo/config"
	"github.com/ranggaaprilio/boilerGo/exception"
	customMiddleware "github.com/ranggaaprilio/boilerGo/middleware"
)

// ServerHeader Config
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

// Init Routing Initialize
func Init() *echo.Echo {
	db := config.CreateCon()

	e := echo.New()
	/** custom Middleware **/
	e.Use(ServerHeader)
	s := customMiddleware.NewStats()
	e.Use(s.Process)

	/** middeleware **/
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	//make Uniq ID for Every single request
	e.Use(middleware.RequestID())

	//Make custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	//Loging
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} ,time=${time_rfc3339}\n",
	}))

	//Recover
	e.Use(middleware.Recover())

	//CORS HANDLER
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	//static file
	e.Static("/", "public")
	/** middeleware **/

	//routing
	e.GET("/healthcheck", s.Handle)
	e.GET("/", func(c echo.Context) error {
		//get go version
		versionSTATUS := "The application was built with the Go version: " + runtime.Version() + "\n"
		return c.String(http.StatusOK, versionSTATUS)
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
