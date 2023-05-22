package config

import (
	"github.com/spf13/viper"
)

// initViper will initialize the viper configuration
// and read in the config file
func initViper() error {
	viper.SetConfigName("dockerleaks")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dockerleaks/")
	viper.AddConfigPath("$HOME/.dockerleaks")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
