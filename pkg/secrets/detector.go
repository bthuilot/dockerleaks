package secrets

import (
	"github.com/sirupsen/logrus"
	"io"
)

// Detector is the interface for a secrets detector.
// It is used to search for secrets in text.
type StaticDetector interface {
	// SearchText searches for secrets in the given text.
	// It returns a slice of matches and an error.
	SearchText(text string) ([]TextMatch, error)
}

type DynamicDetector interface {
	// SearchFile searches for secrets in the given file.
	SearchFile(path string, body io.Reader) ([]FileMatch, error)
}

type Detector interface {
	StaticDetector
	DynamicDetector
}

// Opts is used to configure a Detector.
type Opts struct {
	// UseDefaultStaticRules will include the default rules in the Detector.
	UseDefaultStaticRules bool

	// UseDefaultDynamicRules will include the default rules in the Detector.
	UseDefaultDynamicRules bool
}

// NewDetector creates a new Detector with the given rules,
// and configured with the given Opts.
func NewDetector(opts Opts, staticRules []StaticRule, dynamicRules []DynamicRule) Detector {
	var (
		baseStaticRules  []StaticRule
		baseDynamicRules []DynamicRule
	)
	if opts.UseDefaultStaticRules {
		logrus.Debugf("using default static rules")
		baseStaticRules = DefaultStaticRules
	}
	if opts.UseDefaultDynamicRules {
		logrus.Debugf("using default dynamic rules")
		baseDynamicRules = DefaultDynamicRules
	}
	return detector{
		staticRules:  append(baseStaticRules, staticRules...),
		dynamicRules: append(baseDynamicRules, dynamicRules...),
	}
}

type detector struct {
	staticRules  []StaticRule
	dynamicRules []DynamicRule
}

func (d detector) SearchText(text string) (matches []TextMatch, err error) {
	return findStaticRuleMatches(text, d.staticRules)
}

func (d detector) SearchFile(path string, body io.Reader) (matches []FileMatch, err error) {
	return findDynamicRuleMatches(path, body, d.dynamicRules)
}
