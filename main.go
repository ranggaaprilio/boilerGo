package main

import (
	c "github.com/ranggaaprilio/boilerGo/config"
	"github.com/ranggaaprilio/boilerGo/exception"
	"github.com/ranggaaprilio/boilerGo/routes"
)

func main() {
	defer exception.Catch()

	conf := c.Loadconf()
	c.DbInit()
	bootstrap()
	e := routes.Init()

	e.Logger.Fatal(e.Start(":" + conf.Server.Port))
}
