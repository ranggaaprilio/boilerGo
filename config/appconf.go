package config

import (
	"github.com/ranggaaprilio/boilerGo/exception"
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DbConfigurations
}
type ServerConfigurations struct {
	Name string
	Port string
}

type DbConfigurations struct {
	DbUsername string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

func Loadconf() Configurations {
	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var configuration Configurations

	if err := viper.ReadInConfig(); err != nil {
		exception.PanicIfNeeded(err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		exception.PanicIfNeeded(err)
	}

	return configuration
}
