package config

import (
	"github.com/spf13/viper"
)

func LoadConfiguration(fileName string) error {
	viper.SetConfigName(fileName)
	viper.AddConfigPath("./app/config")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

// - FIXME: Need to find better way for config file when we are in debugging mode
func LoadConfigurationForDebugging() error {
	viper.AddConfigPath("./app/config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}
