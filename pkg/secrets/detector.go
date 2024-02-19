package secrets

// Detector is the interface for a secrets detector.
// It is used to search for secrets in text.
type Detector interface {
	// SearchText searches for secrets in the given text.
	// It returns a slice of matches and an error.
	SearchText(text string) ([]Match, error)
}

// Opts is used to configure a Detector.
type Opts struct {
	// UseDefaultRules will include the default rules in the Detector.
	UseDefaultRules bool
}

// NewDetector creates a new Detector with the given rules,
// and configured with the given Opts.
func NewDetector(opts Opts, rules ...Rule) Detector {
	var baseRules []Rule
	if opts.UseDefaultRules {
		baseRules = DefaultRules
	}
	return &detector{
		rules: append(baseRules, rules...),
	}
}

type detector struct {
	rules []Rule
}

func (d *detector) SearchText(text string) (matches []Match, err error) {
	return findRuleMatches(text, d.rules)
}
