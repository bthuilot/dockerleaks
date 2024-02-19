package config

import (
	"github.com/spf13/viper"
)

const (
	ViperLogLevelKey     = "logLevel"
	ViperUnmaskKey       = "unmaskValues"
	ViperExcludeKey      = "excludeDefaultRules"
	ViperDisableColorKey = "disableColor"
	ViperPullImageKey    = "pullImage"
)

// File is the user configuration file for the application
type File struct {
	// Rules is the list of user defined rules for matching secret strings
	Rules []UserRule
	// ExcludeDefaultRules will disable the default Patterns for detecting
	// secret strings2. See the variable [common.DefaultRules] for the full
	// list of defaults
	ExcludeDefaultRules bool
}

// UserRule represents a user defined string pattern/entropy
// for the layer and filesystem detectors to search
type UserRule struct {
	// Pattern is a regular expression for matching a secret.
	// must be compatible with [re2 syntax]
	//
	// [re2 syntax]: https://github.com/google/re2/wiki/Syntax
	Pattern string
	// Name is a human-readable name of the secret the expression
	// searches for (i.e. AWS SecretString Key, OAuth token, etc.)
	Name string
	// MinEntropy is the minimum entropy the string should have
	MinEntropy float64
}

// Init will initialize the configuration of the application
func Init() error {
	if err := initViper(); err != nil {
		return err
	}
	initLogger(viper.GetString(ViperLogLevelKey))
	return nil
}
