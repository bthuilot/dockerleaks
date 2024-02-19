package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"strings"
)

// LoggingLevel is the verbosity level of logging
type LoggingLevel = string

const (
	// Debug is the most verbose logging level and will
	// report all logs
	Debug = "DEBUG"
	// Info level reports logrus.Info, logrus.Warn, and logrus.Error
	Info = "INFO"
	// Warn level reports logrus.Warn, and logrus.Error
	Warn = "WARN"
	// Error level reports ony error messages
	Error = "ERROR"
	// Off level turns off logs from logrus. instead the "stylized"
	// logging using terminal spinners will be used
	Off = "OFF"
)

// initLogger will initialize the logging configuration of the application
// and set the logging level to the given LoggingLevel. Defaults to Off
func initLogger(level LoggingLevel) {
	switch strings.ToUpper(level) {
	case Debug:
		logrus.SetLevel(logrus.DebugLevel)
	case Info:
		logrus.SetLevel(logrus.InfoLevel)
	case Warn:
		logrus.SetLevel(logrus.WarnLevel)
	case Error:
		logrus.SetLevel(logrus.ErrorLevel)
	case Off:
		fallthrough
	default:
		logrus.SetOutput(io.Discard)
	}
}

// ShouldUseSpinner will return true if the "stylized"
func ShouldUseSpinner() bool {
	return !viper.IsSet(ViperLogLevelKey) || strings.ToUpper(viper.GetString(ViperLogLevelKey)) == "OFF"
}

// ShouldUseColor will return true, if colored output should be used
func ShouldUseColor() bool {
	return !viper.GetBool("disable_color")
}
