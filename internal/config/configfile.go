package config

import (
	"github.com/spf13/viper"
)

const (
	ViperLogLevelKey     = "logLevel"
	ViperUnmaskKey       = "unmaskValues"
	ViperDisableColorKey = "disableColor"
	ViperConfigFileKey   = "configFile"
)

// File is the user configuration file for the application
type File struct {
	// StaticRules is the list of user defined rules for matching secret strings
	// during a static image analysis
	StaticRules []UserStaticRule
	// DynamicRules is the list of user defined rules for matching secret strings
	// during a dynamic container analysis
	DynamicRules []UserDynamicRule
	// IgnoreInvalidRules will ignore any invalid rules in the configuration
	// file if set to true
	IgnoreInvalidRules bool
	// ExcludeDefaultStaticRules will disable the default Patterns for detecting
	// secret strings during a static scan. See the variable [secrets.DefaultStaticRules] for the full
	// list of defaults
	ExcludeDefaultStaticRules bool
	// ExcludeDefaultDynamicRules will disable the default rules for detecting
	// secret strings or files during a dynamic scan. See the variable [secrets.DefaultDynamicRules] for the full
	// list of defaults
	ExcludeDefaultDynamicRules bool
}

// UserStaticRule represents a user defined string pattern/entropy
// for the layer and filesystem detectors to search
type UserStaticRule struct {
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

type UserDynamicRule struct {
	// Name is a human-readable name of the secret the expression
	// searches for (i.e. .env files, tfstate , etc.)
	Name string
	// FilePattern is a regular expression for matching files to search
	// a nil value means that the rule will match all files
	FilePattern string
	// Pattern is a regular expression for matching text in the file
	// a nil value means that the rule will return true if only the file is matched
	// (matching all the file)
	Pattern string
	// MinEntropy is the minimum entropy the string should have
	MinEntropy float64
}

// Init will initialize the configuration of the application
func Init() error {
	if err := initViper(); err != nil {
		return err
	}

	logLevel := "info"
	if viper.IsSet(ViperLogLevelKey) {
		logLevel = viper.GetString(ViperLogLevelKey)
	}
	return initLogger(logLevel)
}
