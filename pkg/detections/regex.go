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

type Pattern struct {
	// Expression is a regular expression for matching a secret.
	// must be compatible with RE2 Syntax
	// TODO(add link)
	Expression string
	// Name is a human-readable name of the secret the expression
	// searches for (i.e. AWS Secret Key, OAuth token, etc.)
	Name string
}

func NewRegexDetector(patterns []Pattern) (Detector, error) {
	var parsedPatterns []pattern
	for _, p := range patterns {
		parsed, err := regexp.Compile(p.Expression)
		if err != nil {
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

func (r Regexp) findMatch(s string) (name string, match bool) {
	for _, p := range r.patterns {
		if p.RegExp.MatchString(s) {
			return p.Name, true
		}
	}
	return "", false
}

func (r Regexp) Run(img image.Image) []Detection {
	var detections []Detection
	if envVars, use := img.GetEnvVars(); use {
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
	}

	if buildArgs, use := img.GetBuildArgs(); use {
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
	}

	return detections
}
