package config

import (
	"github.com/bthuilot/dockerleaks/pkg/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"regexp"
)

// File is the user configuration file for the application
type File struct {
	// Layer is the configuration of the layer detector
	// (i.e. build arguments, environment variables and statically defined creds)
	Layer LayerConfig
	// Filesystem is the configuration of the filesystem detector
	// (i.e. file and folder content)
	Filesystem FilesystemConfig
	// Rules is the list of user defined rules for matching secret strings
	Rules []UserRule
	// ExcludeDefaultRules will disable the default Patterns for detecting
	// secret strings2. See the variable [common.DefaultRules] for the full
	// list of defaults
	ExcludeDefaultRules bool
}

// LayerConfig is the configuration of the layer detector, which detects secret strings
// in build arguments and environment variables
type LayerConfig struct {
	// Disable is a boolean indicating whether to run the detector
	Disable bool
}

type FilesystemConfig struct {
	// Disable is a boolean indicating whether to run the detector
	Disable bool
	// TODO(filesystem config)
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
	// Entropy is the minimum entropy the string should have
	Entropy float64
}

// Init will initialize the configuration of the application
func Init() error {
	if err := initViper(); err != nil {
		return err
	}
	initLogger(viper.GetString("log_level"))
	return nil
}

// ParseRules will parse a list of UserRule patterns into regexp.Regexp and a common.SecretStringRule.
// All rules that result in error are returned in the second variables
func ParseRules(userRules []UserRule) (rules []common.SecretStringRule, errors []UserRule) {
	for _, r := range userRules {
		regex, err := regexp.Compile(r.Pattern)
		if err != nil {
			logrus.Errorf(
				"unable to parse regular expression %s `%s`: %s",
				r.Name, r.Pattern, err,
			)
			errors = append(errors, r)
		}
		rules = append(rules, common.SecretStringRule{
			Pattern: regex,
			Name:    r.Name,
			Entropy: r.Entropy,
		})
	}
	return
}
