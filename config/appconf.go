package config

import (
	"fmt"

	"github.com/ranggaaprilio/boilerGo/exception"
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server ServerConfigurations
}
type ServerConfigurations struct {
	Name string
	Port string
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
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return configuration
}
