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

	// Set up mappings for environment variables
	viper.SetEnvPrefix("")
	viper.BindEnv("database.dbusername", "DB_USER")
	viper.BindEnv("database.dbpassword", "DB_PASSWORD")
	viper.BindEnv("database.dbhost", "DB_HOST")
	viper.BindEnv("database.dbport", "DB_PORT")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("server.port", "SERVER_PORT")

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
