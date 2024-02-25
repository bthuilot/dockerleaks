package secrets

import (
	"fmt"
	"github.com/bthuilot/dockerleaks/internal/config"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// Rule represents a rule for matching secret strings
type Rule interface {
	String() string
}

// StaticRule represents a pattern and entropy rule for matching
// secret string in a static context
type StaticRule struct {
	// Name is the human-readable name secret that this
	// rule detects
	Name string `json:"name"`
	// Pattern is the regular expression to match this secret
	Pattern *regexp.Regexp `json:"pattern"`
	// MinEntropy is the minimum entropy the string must be
	MinEntropy float64 `json:"min_entropy,omitempty"`
}

func (r StaticRule) String() string {
	var conditions []string
	if r.Pattern != nil {
		conditions = append(conditions, fmt.Sprintf("regex '%s'", r.Pattern))
	}
	if r.MinEntropy > 0 {
		conditions = append(conditions, fmt.Sprintf("minimum entropy of %f", r.MinEntropy))
	}
	if len(conditions) > 0 {
		return fmt.Sprintf("'%s' via %s", r.Name, strings.Join(conditions, " and "))
	}
	return fmt.Sprintf("'%s'", r.Name)
}

type DynamicRule struct {
	// Name is the human-readable name secret that this
	// rule detects
	Name string `json:"name"`
	// FilePattern is the regular expression to match the files to search
	// a nil value means that the rule will match all files
	FilePattern *regexp.Regexp `json:"file_pattern,omitempty"`
	// Pattern is the regular expression to match text in the file
	// a nil value means that the rule will return true if only the file is matched
	// (matching all the file)
	Pattern *regexp.Regexp `json:"pattern,omitempty"`
	// MinEntropy is the minimum entropy the string must be.
	// This will only be used if Pattern is not nil
	// a value of 0 means that the entropy will not be checked
	MinEntropy float64 `json:"min_entropy,omitempty"`
}

func (r DynamicRule) String() string {
	var conditions []string
	if r.Pattern != nil {
		conditions = append(conditions, fmt.Sprintf("regex '%s'", r.Pattern))
	}
	if r.MinEntropy > 0 {
		conditions = append(conditions, fmt.Sprintf("minimum entropy of %f", r.MinEntropy))
	}
	if r.FilePattern != nil {
		conditions = append(conditions, fmt.Sprintf("file pattern '%s'", r.FilePattern))
	}
	if len(conditions) > 0 {
		return fmt.Sprintf("'%s' via %s", r.Name, strings.Join(conditions, " and "))
	}
	return fmt.Sprintf("'%s'", r.Name)
}

// DefaultStaticRules is the default list of rules
// this list contains rules to match a common
// set of secrets
// TODO(improve this list)
var DefaultStaticRules = []StaticRule{
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

var DefaultDynamicRules = []DynamicRule{
	{
		FilePattern: regexp.MustCompile(`^(.*/)*[-\w._]*\.env(\.[-\w._]*)?$`),
		Name:        ".env file",
	},
	{
		Name:        "Terraform state file",
		FilePattern: regexp.MustCompile(`^(.*/)*terraform.tfstate$`),
	},
	{
		Name:        "AWS credentials file",
		FilePattern: regexp.MustCompile(`^(.*/)*\.aws/credentials$`),
	},
}

// ParseStaticRules will parse a list of UserRule patterns into regexp.Regexp and a common.SecretStringRule.
// All rules that result in error are returned in the second variables
func ParseStaticRules(userRules []config.UserStaticRule) (rules []StaticRule, errors []config.UserStaticRule) {
	for _, r := range userRules {
		regex, err := regexp.Compile(r.Pattern)
		if err != nil {
			logrus.Errorf(
				"unable to parse regular expression %s `%s`: %s",
				r.Name, r.Pattern, err,
			)
			errors = append(errors, r)
		}
		rules = append(rules, StaticRule{
			Pattern:    regex,
			Name:       r.Name,
			MinEntropy: r.MinEntropy,
		})
	}
	return
}

// ParseDynamicRules will parse a list of UserRule patterns into regexp.Regexp and a common.SecretStringRule.
// All rules that result in error are returned in the second variables
func ParseDynamicRules(userRules []config.UserDynamicRule) (rules []DynamicRule, errors []config.UserDynamicRule) {
	for _, r := range userRules {
		var rule DynamicRule
		rule.Name = r.Name
		if r.FilePattern != "" {
			regex, err := regexp.Compile(r.FilePattern)
			if err != nil {
				logrus.Errorf(
					"unable to parse regular expression %s `%s`: %s",
					r.Name, r.FilePattern, err,
				)
				errors = append(errors, r)
				continue
			}
			rule.FilePattern = regex
		}
		if r.Pattern != "" {
			regex, err := regexp.Compile(r.Pattern)
			if err != nil {
				logrus.Errorf(
					"unable to parse regular expression %s `%s`: %s",
					r.Name, r.Pattern, err,
				)
				errors = append(errors, r)
				continue
			}
			rule.Pattern = regex
			rule.MinEntropy = r.MinEntropy
		}
		rules = append(rules, rule)
	}
	return
}
