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
	patterns []Pattern
}

// String returns the string formatted name of the type
// of detections performed
func (r Regexp) String() string {
	return "Regular Expressions"
}

// Pattern is a regular expression to search for in
// docker image environment variables, files and build arguments
type Pattern struct {
	// RegExp is the compiled regular expression to search for
	RegExp *regexp.Regexp
	// Name is a human-readable name of the secret the expression
	// searches for (i.e. AWS Secret Key, OAuth token, etc.)
	Name string
}

// NewRegexDetector will construct a new Detector that will search all
// environment variables, build arguments and contents of files on the file
// system for strings that matches any of the given Pattern
func NewRegexDetector(patterns []Pattern) (Detector, error) {
	return Regexp{
		patterns: patterns,
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
			logrus.Debugf(
				"env var %s (value: %s) was found to match %s",
				v.Name, v.Value, name,
			)
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
			logrus.Debugf(
				"build argument %s (value: %s) was found to match %s",
				v.Name, v.Value, name,
			)
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

// DefaultPatterns is a list of patterns
// to include in the detector by default
var DefaultPatterns = []Pattern{
	{
		RegExp: regexp.MustCompile(`[1-9][0-9]+-[0-9a-zA-Z]{40}`),
		Name:   "Twitter",
	}, {
		RegExp: regexp.MustCompile(`/(^|[^@\w])@(\w{1,15})\b/`),
		Name:   "Twitter",
	}, {
		RegExp: regexp.MustCompile(`EAACEdEose0cBA[0-9A-Za-z]+`),
		Name:   "Facebook",
	}, {
		RegExp: regexp.MustCompile(`[A-Za-z0-9]{125}`),
		Name:   "Facebook",
	}, {
		RegExp: regexp.MustCompile(`[0-9a-fA-F]{7}\.[0-9a-fA-F]{32}`),
		Name:   "Instagram",
	}, {
		RegExp: regexp.MustCompile(`AIza[0-9A-Za-z-_]{35}`),
		Name:   "Google",
	}, {
		RegExp: regexp.MustCompile(`[0-9a-zA-Z\-_]{24}`),
		Name:   "Google",
	}, {
		RegExp: regexp.MustCompile(`4/[0-9A-Za-z\-_]+`),
		Name:   "Google",
	}, {
		RegExp: regexp.MustCompile(`1/[0-9A-Za-z\-_]{43}|1/[0-9A-Za-z\-_]{64}`),
		Name:   "Google",
	}, {
		RegExp: regexp.MustCompile(`ya29\.[0-9A-Za-z\-_]+`),
		Name:   "Google",
	}, {
		RegExp: regexp.MustCompile(`^ghp_[a-zA-Z0-9]{36}$`),
		Name:   "GitHub",
	}, {
		RegExp: regexp.MustCompile(`^github_pat_[a-zA-Z0-9]{22}_[a-zA-Z0-9]{59}$`),
		Name:   "GitHub",
	}, {
		RegExp: regexp.MustCompile(`^gho_[a-zA-Z0-9]{36}$`),
		Name:   "GitHub",
	}, {
		RegExp: regexp.MustCompile(`^ghu_[a-zA-Z0-9]{36}$`),
		Name:   "GitHub",
	}, {
		RegExp: regexp.MustCompile(`^ghs_[a-zA-Z0-9]{36}$`),
		Name:   "GitHub",
	}, {
		RegExp: regexp.MustCompile(`^ghr_[a-zA-Z0-9]{36}$`),
		Name:   "GitHub",
	}, {
		RegExp: regexp.MustCompile(`([s,p]k.eyJ1Ijoi[\w\.-]+)`),
		Name:   "Mapbox",
	}, {
		RegExp: regexp.MustCompile(`([s,p]k.eyJ1Ijoi[\w\.-]+)`),
		Name:   "Mapbox",
	}, {
		RegExp: regexp.MustCompile(`R_[0-9a-f]{32}`),
		Name:   "Foursquare",
	}, {
		RegExp: regexp.MustCompile(`sk_live_[0-9a-z]{32}`),
		Name:   "Picatic",
	}, {
		RegExp: regexp.MustCompile(`sk_live_[0-9a-zA-Z]{24}`),
		Name:   "Stripe",
	}, {
		RegExp: regexp.MustCompile(`sk_live_[0-9a-zA-Z]{24}`),
		Name:   "Stripe",
	}, {
		RegExp: regexp.MustCompile(`sqOatp-[0-9A-Za-z\-_]{22}`),
		Name:   "Square",
	}, {
		RegExp: regexp.MustCompile(`q0csp-[0-9A-Za-z\-_]{43}`),
		Name:   "Square",
	}, {
		RegExp: regexp.MustCompile(`access_token\,production\$[0-9a-z]{161}[0-9a,]{32}`),
		Name:   "Paypal / Braintree",
	}, {
		RegExp: regexp.MustCompile(`amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-10-9a-f1{4}-[0-9a,]{4}-[0-9a-f]{12}`),
		Name:   "Amazon Marketing Services",
	}, {
		RegExp: regexp.MustCompile(`55[0-9a-fA-F]{32}`),
		Name:   "Twilio",
	}, {
		RegExp: regexp.MustCompile(`key-[0-9a-zA-Z]{32}`),
		Name:   "MailGun",
	}, {
		RegExp: regexp.MustCompile(`[ 0-9a-f ]{ 32 }-us[0-9]{1,2}`),
		Name:   "MailChimp",
	}, {
		RegExp: regexp.MustCompile(`xoxb-[0-9]{11}-[0-9]{11}-[0-9a-zA-Z]{24}`),
		Name:   "Slack",
	}, {
		RegExp: regexp.MustCompile(`xoxp-[0-9]{11}-[0-9]{11}-[0-9a-zA-Z]{24}`),
		Name:   "Slack",
	}, {
		RegExp: regexp.MustCompile(`xoxe.xoxp-1-[0-9a-zA-Z]{166}`),
		Name:   "Slack",
	}, {
		RegExp: regexp.MustCompile(`xoxe-1-[0-9a-zA-Z]{147}`),
		Name:   "Slack",
	}, {
		RegExp: regexp.MustCompile(`T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}`),
		Name:   "Slack",
	}, {
		RegExp: regexp.MustCompile(`A[KS]IA[0-9A-Z]{16}`),
		Name:   "Amazon Web Services",
	}, {
		RegExp: regexp.MustCompile(`[0-9a-zA-Z/+]{40}`),
		Name:   "Amazon Web Services",
	}, {
		RegExp: regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`),
		Name:   "Google Cloud Platform",
	}, {
		RegExp: regexp.MustCompile(`[A-Za-z0-9_]{21}--[A-Za-z0-9_]{8}`),
		Name:   "Google Cloud Platform",
	}, {
		RegExp: regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`),
		Name:   "Heroku",
	}, {
		RegExp: regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`),
		Name:   "Heroku",
	}}
