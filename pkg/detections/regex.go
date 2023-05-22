package detections

import (
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/sirupsen/logrus"
	"regexp"
)

// Regexp is a Detector implementation for detecting
// secrets using regular expression.
type Regexp struct {
	// patterns are the list of regular expression
	// to check for
	patterns []pattern
}

// pattern is a parsed and validated Pattern,
// comprised of a compiled regular expression
// and the name of the secret it matches
type pattern struct {
	// RegExp is the compiled regular expression to search for
	RegExp *regexp.Regexp
	// Name is a human-readable name of the secret the expression
	// searches for (i.e. AWS Secret Key, OAuth token, etc.)
	Name string
}

// Pattern reprsents a user defined pattern for the Regexp Detector to
// search for.
type Pattern struct {
	// Expression is a regular expression for matching a secret.
	// must be compatible with RE2 Syntax
	// TODO(add link to RE2 syntax)
	Expression string
	// Name is a human-readable name of the secret the expression
	// searches for (i.e. AWS Secret Key, OAuth token, etc.)
	Name string
}

// NewRegexDetector will construct a new Detector that will search all
// environment variables, build arguments and contents of files on the file
// system for strings that matches any of the given Pattern
func NewRegexDetector(patterns []Pattern) (Detector, error) {
	var parsedPatterns []pattern
	for _, p := range patterns {
		parsed, err := regexp.Compile(p.Expression)
		if err != nil {
			// TODO(should i return an error, or just continue?)
			logrus.Warnf("skipping regex '%s', invalid\n", p.Name)
			continue
		}
		parsedPatterns = append(parsedPatterns, pattern{
			RegExp: parsed,
			Name:   p.Name,
		})
	}
	return Regexp{
		patterns: parsedPatterns,
	}, nil
}

// findMatch will see if there is any match of the given string with
// the supplied regular expression patterns. It will return a bool
// indicating there is a match, and will return the name of the regular
// expression matched
func (r Regexp) findMatch(s string) (name string, match bool) {
	for _, p := range r.patterns {
		if p.RegExp.MatchString(s) {
			return p.Name, true
		}
	}
	return "", false
}

// EvalEnvVars will evaluate the environment variables to see if any of them have
// a value that matches one of the configured Pattern
func (r Regexp) EvalEnvVars(envVars []image.EnvVar) (detections []Detection) {
	for _, v := range envVars {
		if name, match := r.findMatch(v.Value); match {
			logrus.Infof("found match %s for env var %s", name, v.Name)
			detections = append(detections, Detection{
				Name:     name,
				Type:     RegexDetection,
				Source:   EnvVarSecret,
				Value:    v.Value,
				Location: v.Location,
			})
		}
	}
	return
}

// EvalBuildArgs will evaluate the build arguments to see if any of them have
// a value that matches one of the configured Pattern
func (r Regexp) EvalBuildArgs(buildArgs []image.BuildArg) (detections []Detection) {
	for _, v := range buildArgs {
		if name, match := r.findMatch(v.Value); match {
			detections = append(detections, Detection{
				Name:     name,
				Type:     RegexDetection,
				Source:   BuildArgSecret,
				Value:    v.Value,
				Location: v.Location,
			})
		}
	}

	return
}
