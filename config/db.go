package config

import (
	"github.com/ranggaaprilio/boilerGo/exception"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

//Init return config DB
func DbInit() {
	conf := Loadconf()
	connectionString := conf.Database.DbUsername + ":" + conf.Database.DbPassword + "@tcp(" + conf.Database.DbHost + ":" + conf.Database.DbPort + ")/" + conf.Database.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		exception.PanicIfNeeded(err)
	}

}

//CreateCon return var db
func CreateCon() *gorm.DB {
	return db
}
