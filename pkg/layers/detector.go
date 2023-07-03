package layers

import (
	"fmt"
	"github.com/bthuilot/dockerleaks/pkg/common"
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/sirupsen/logrus"
)

// Detector represents a struct to detect secret strings
// located in the layers of a docker image.
type Detector struct {
	patterns []common.SecretStringRule
	img      image.Image
}

// NewDetector will construct a new image layer secrets detector
func NewDetector(img image.Image) Detector {
	return Detector{img: img}
}

// WithRules will add the list of rules to this detection list.
// NOTE: calling this method multiple times with the duplicate rules will r
// result in duplicate findings being returned
func (d Detector) WithRules(patterns ...common.SecretStringRule) Detector {
	d.patterns = append(d.patterns, patterns...)
	return d
}

// UseDefaultRules will add the default rules to this list of rules to detect.
// NOTE: calling this method multiple times will result in duplicate rules being added
func (d Detector) UseDefaultRules() Detector {
	d.patterns = append(d.patterns, common.DefaultRules...)
	return d
}

// Detect will look for any strings that match any common.SecretStringRule in
// environment variables and build arguments
func (d Detector) Detect() (detections []common.SecretString, err error) {
	var (
		buildArgs []image.BuildArg
		envVars   []image.EnvVar
	)
	// Check Build Args
	if buildArgs, err = d.img.ParseBuildArguments(); err != nil {
		logrus.Error("unable to parse build arguments: %s", err)
		return nil, fmt.Errorf("an unknown error occurred while parsing build arguments")
	}
	for _, v := range buildArgs {
		if matches := common.FindRuleMatches(v.Value, d.patterns); len(matches) != 0 {
			logrus.Debugf(
				"build argument %s (value: %s) has %d match(es): %#v",
				v.Name, v.Value, len(matches), matches,
			)
			for _, match := range matches {
				detections = append(detections, common.SecretString{
					Name:     match.Rule.Name,
					Source:   common.BuildArgument,
					Value:    match.Value,
					Location: v.Location,
				})
			}
		}
	}

	// Check Build Args
	if envVars, err = d.img.ParseEnvVars(); err != nil {
		logrus.Error("unable to parse build arguments: %s", err)
		return nil, fmt.Errorf("an unknown error occurred while parsing environment variables")
	}
	for _, v := range envVars {
		if matches := common.FindRuleMatches(v.Value, d.patterns); len(matches) != 0 {
			logrus.Debugf(
				"environment variable %s (value: %s) has %d match(es): %#v",
				v.Name, v.Value, len(matches), matches,
			)
			for _, match := range matches {
				detections = append(detections, common.SecretString{
					Name:     match.Rule.Name,
					Source:   common.EnvVar,
					Value:    match.Value,
					Location: v.Location,
				})
			}
		}
	}
	return
}
