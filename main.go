// Package main for the BoilerGo API
//
// @title BoilerGo API
// @version 1.0
// @description This is the API documentation for the BoilerGo application.
//
// @contact.name API Support
// @contact.url https://github.com/ranggaaprilio/boilerGo
// @contact.email support@example.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api
// @schemes http https
package main

import (
	c "github.com/ranggaaprilio/boilerGo/config"
	_ "github.com/ranggaaprilio/boilerGo/docs" // Import swagger docs
	"github.com/ranggaaprilio/boilerGo/exception"
	"github.com/ranggaaprilio/boilerGo/routes"
)

func main() {
	defer exception.Catch()

	conf := c.Loadconf()
	c.DbInit()
	// bootstrap()
	e := routes.Init()

	e.Logger.Fatal(e.Start(":" + conf.Server.Port))
}
