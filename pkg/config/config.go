package config

import (
	"github.com/bthuilot/dockerleaks/pkg/detections"
	"github.com/spf13/viper"
)

// Config is the user configuration for the application
type Config struct {
	// Regexp is the configuration for the regular expression detector
	Regexp RegexpConfig
	// Entropy is the configuration for the string entropy detector
	Entropy EntropyConfig
}

// RegexpConfig represents the configuration for the
type RegexpConfig struct {
	// Patterns is the list of regular expressions
	// patterns to search for
	Patterns []detections.Pattern
}

// EntropyConfig configuration for the string entropy detector
type EntropyConfig struct {
	// TODO(entropy config)
}

// Init will initialize the configuration of the application
func Init() error {
	if err := initViper(); err != nil {
		return err
	}
	initLogger(viper.GetString("log_level"))
	return nil
}
