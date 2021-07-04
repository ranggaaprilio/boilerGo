package main

import (
	c "github.com/ranggaaprilio/boilerGo/config"
	"github.com/ranggaaprilio/boilerGo/routes"
)

func main() {
	conf := c.Loadconf()
	e := routes.Init()

	e.Logger.Fatal(e.Start(":" + conf.Server.Port))
}
