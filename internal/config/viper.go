package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
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
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// TODO(somehow log this after logger set up)
		} else {
			logrus.Errorf("unable to parse in configuration: %s", err)
			return fmt.Errorf("unable to parse configuration file")
		}
	}

	return nil
}
