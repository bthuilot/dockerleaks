package config

import (
	"github.com/spf13/viper"
)

func initViper() error {
	viper.SetConfigName("dockerleaks")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dockerleaks/")
	viper.AddConfigPath("$HOME/.dockerleaks")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
