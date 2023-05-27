package config

import (
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
	Patterns []Pattern

	// DisableDefaults will disable the default regular expressions
	// from begin include in the regular expression detector. See full
	// list for the detections.DefaultPatterns
	DisableDefaults bool
}

// Pattern reprsents a user defined pattern for the Regexp Detector to
// search for.
type Pattern struct {
	// Expression is a regular expression for matching a secret.
	// must be compatible with [re2 syntax]
	//
	// [re2 syntax]: https://github.com/google/re2/wiki/Syntax
	Expression string
	// Name is a human-readable name of the secret the expression
	// searches for (i.e. AWS Secret Key, OAuth token, etc.)
	Name string
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
