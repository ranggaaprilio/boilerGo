package main

import (
	u "github.com/ranggaaprilio/boilerGo/app/v1/modules/user"
	c "github.com/ranggaaprilio/boilerGo/config"
)

func bootstrap() {
	db := c.CreateCon()
	db.AutoMigrate(&u.User{})
}
