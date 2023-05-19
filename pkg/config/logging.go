package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"strings"
)

func initLogger() {
	switch strings.ToUpper(viper.GetString("log_level")) {
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "OFF":
		fallthrough
	default:
		logrus.SetOutput(io.Discard)
	}
}

func ShouldUseSpinner() bool {
	return !viper.IsSet("log_level") || viper.GetString("log_level") == "OFF"
}

func ShouldUseColor() bool {
	return !viper.GetBool("disable_color")
}
