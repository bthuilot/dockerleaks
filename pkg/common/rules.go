package common

import "regexp"

// SecretStringRule represents a pattern and entropy rule for matching
// secret string
type SecretStringRule struct {
	// Name is the human-readable name secret that this
	// rule detects
	Name string
	// Pattern is the regular expression to match this secret
	Pattern *regexp.Regexp
	// Entropy is the minimum entropy the string must be
	Entropy float64
}

// SecretStringMatch represents a match of string that is detected to be a secret value
type SecretStringMatch struct {
	// Rule is the rule that matches this string
	Rule SecretStringRule
	// Value is the actual value of the string
	Value string
}

// FindRuleMatches will search a string's content for for any matches to
// the list of SecretStringRules provided
func FindRuleMatches(content string, rules []SecretStringRule) (matches []SecretStringMatch) {
	for _, r := range rules {
		// TODO(entropy)
		for _, m := range r.Pattern.FindStringSubmatch(content) {
			matches = append(matches, SecretStringMatch{
				Rule:  r,
				Value: m,
			})
		}
	}
	return
}

// DefaultRules is the default list of rules
// this list contains rules to match a common
// set of secrets
// TODO(entropy for default detections)
var DefaultRules = []SecretStringRule{
	{
		Pattern: regexp.MustCompile(`[1-9][0-9]+-[0-9a-zA-Z]{40}`),
		Name:    "Twitter",
	}, {
		Pattern: regexp.MustCompile(`/(^|[^@\w])@(\w{1,15})\b/`),
		Name:    "Twitter",
	}, {
		Pattern: regexp.MustCompile(`EAACEdEose0cBA[0-9A-Za-z]+`),
		Name:    "Facebook",
	}, {
		Pattern: regexp.MustCompile(`[A-Za-z0-9]{125}`),
		Name:    "Facebook",
	}, {
		Pattern: regexp.MustCompile(`[0-9a-fA-F]{7}\.[0-9a-fA-F]{32}`),
		Name:    "Instagram",
	}, {
		Pattern: regexp.MustCompile(`AIza[0-9A-Za-z-_]{35}`),
		Name:    "Google",
	}, {
		Pattern: regexp.MustCompile(`[0-9a-zA-Z\-_]{24}`),
		Name:    "Google",
	}, {
		Pattern: regexp.MustCompile(`4/[0-9A-Za-z\-_]+`),
		Name:    "Google",
	}, {
		Pattern: regexp.MustCompile(`1/[0-9A-Za-z\-_]{43}|1/[0-9A-Za-z\-_]{64}`),
		Name:    "Google",
	}, {
		Pattern: regexp.MustCompile(`ya29\.[0-9A-Za-z\-_]+`),
		Name:    "Google",
	}, {
		Pattern: regexp.MustCompile(`^ghp_[a-zA-Z0-9]{36}$`),
		Name:    "GitHub",
	}, {
		Pattern: regexp.MustCompile(`^github_pat_[a-zA-Z0-9]{22}_[a-zA-Z0-9]{59}$`),
		Name:    "GitHub",
	}, {
		Pattern: regexp.MustCompile(`^gho_[a-zA-Z0-9]{36}$`),
		Name:    "GitHub",
	}, {
		Pattern: regexp.MustCompile(`^ghu_[a-zA-Z0-9]{36}$`),
		Name:    "GitHub",
	}, {
		Pattern: regexp.MustCompile(`^ghs_[a-zA-Z0-9]{36}$`),
		Name:    "GitHub",
	}, {
		Pattern: regexp.MustCompile(`^ghr_[a-zA-Z0-9]{36}$`),
		Name:    "GitHub",
	}, {
		Pattern: regexp.MustCompile(`([s,p]k.eyJ1Ijoi[\w\.-]+)`),
		Name:    "Mapbox",
	}, {
		Pattern: regexp.MustCompile(`([s,p]k.eyJ1Ijoi[\w\.-]+)`),
		Name:    "Mapbox",
	}, {
		Pattern: regexp.MustCompile(`R_[0-9a-f]{32}`),
		Name:    "Foursquare",
	}, {
		Pattern: regexp.MustCompile(`sk_live_[0-9a-z]{32}`),
		Name:    "Picatic",
	}, {
		Pattern: regexp.MustCompile(`sk_live_[0-9a-zA-Z]{24}`),
		Name:    "Stripe",
	}, {
		Pattern: regexp.MustCompile(`sk_live_[0-9a-zA-Z]{24}`),
		Name:    "Stripe",
	}, {
		Pattern: regexp.MustCompile(`sqOatp-[0-9A-Za-z\-_]{22}`),
		Name:    "Square",
	}, {
		Pattern: regexp.MustCompile(`q0csp-[0-9A-Za-z\-_]{43}`),
		Name:    "Square",
	}, {
		Pattern: regexp.MustCompile(`access_token\,production\$[0-9a-z]{161}[0-9a,]{32}`),
		Name:    "Paypal / Braintree",
	}, {
		Pattern: regexp.MustCompile(`amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-10-9a-f1{4}-[0-9a,]{4}-[0-9a-f]{12}`),
		Name:    "Amazon Marketing Services",
	}, {
		Pattern: regexp.MustCompile(`55[0-9a-fA-F]{32}`),
		Name:    "Twilio",
	}, {
		Pattern: regexp.MustCompile(`key-[0-9a-zA-Z]{32}`),
		Name:    "MailGun",
	}, {
		Pattern: regexp.MustCompile(`[ 0-9a-f ]{ 32 }-us[0-9]{1,2}`),
		Name:    "MailChimp",
	}, {
		Pattern: regexp.MustCompile(`xoxb-[0-9]{11}-[0-9]{11}-[0-9a-zA-Z]{24}`),
		Name:    "Slack",
	}, {
		Pattern: regexp.MustCompile(`xoxp-[0-9]{11}-[0-9]{11}-[0-9a-zA-Z]{24}`),
		Name:    "Slack",
	}, {
		Pattern: regexp.MustCompile(`xoxe.xoxp-1-[0-9a-zA-Z]{166}`),
		Name:    "Slack",
	}, {
		Pattern: regexp.MustCompile(`xoxe-1-[0-9a-zA-Z]{147}`),
		Name:    "Slack",
	}, {
		Pattern: regexp.MustCompile(`T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}`),
		Name:    "Slack",
	}, {
		Pattern: regexp.MustCompile(`A[KS]IA[0-9A-Z]{16}`),
		Name:    "Amazon Web Services",
	}, {
		Pattern: regexp.MustCompile(`[0-9a-zA-Z/+]{40}`),
		Name:    "Amazon Web Services",
	}, {
		Pattern: regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`),
		Name:    "Google Cloud Platform",
	}, {
		Pattern: regexp.MustCompile(`[A-Za-z0-9_]{21}--[A-Za-z0-9_]{8}`),
		Name:    "Google Cloud Platform",
	}, {
		Pattern: regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`),
		Name:    "Heroku",
	}, {
		Pattern: regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`),
		Name:    "Heroku",
	},
}
