package secrets

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestFindStaticRuleMatches(t *testing.T) {

	passwordRule := StaticRule{
		Name:       "match password",
		Pattern:    regexp.MustCompile(`^password$`),
		MinEntropy: 0,
	}

	passwordVarRule := StaticRule{
		Name:    "match password var",
		Pattern: regexp.MustCompile(`(?m)^(var)?\s*password\s*(:)?=\s*['"](.*)['"]$`),
	}

	authorizationHeader := StaticRule{
		Name:    "match authorization header",
		Pattern: regexp.MustCompile(`"Authorization": (.*)[,}\n]`),
	}

	tests := []struct {
		name  string
		text  string
		rules []StaticRule
		want  []TextMatch
		err   error
	}{
		{
			name:  "empty string",
			text:  "",
			rules: []StaticRule{},
			want:  []TextMatch{},
			err:   nil,
		}, {
			name: "single match",
			text: "password",
			rules: []StaticRule{
				passwordRule,
			},
			want: []TextMatch{
				{
					Rule: passwordRule,
					Secret: Secret{
						Value:   "password",
						Entropy: CalculateShannonEntropy("password"),
					},
					FullText: "password",
				},
			},
			err: nil,
		},
		{
			name: "multiple matches",
			text: "var password = 'mypassword'\nheaders := map[string]string{\"Authorization\": fmt.Sprintf(\"Bearer %s\", password)}",
			rules: []StaticRule{
				passwordVarRule,
				authorizationHeader,
			},
			want: []TextMatch{
				{
					Rule: passwordVarRule,
					Secret: Secret{
						Value:   "var password = 'mypassword'",
						Entropy: CalculateShannonEntropy("var password = 'mypassword'"),
					},
					FullText: "var password = 'mypassword'\nheaders := map[string]string{\"Authorization\": fmt.Sprintf(\"Bearer %s\", password)}",
				},
				{
					Rule: authorizationHeader,
					Secret: Secret{
						Value:   "\"Authorization\": fmt.Sprintf(\"Bearer %s\", password)}",
						Entropy: CalculateShannonEntropy("\"Authorization\": fmt.Sprintf(\"Bearer %s\", password)}"),
					},
					FullText: "var password = 'mypassword'\nheaders := map[string]string{\"Authorization\": fmt.Sprintf(\"Bearer %s\", password)}",
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, err := findStaticRuleMatches(tt.text, tt.rules)
			for i, match := range matches {
				expected := tt.want[i]
				assert.Equal(t, expected.Rule, match.Rule)
				assert.Equal(t, expected.Secret.Value, match.Secret.Value)
				assert.InDelta(t, expected.Secret.Entropy, match.Secret.Entropy, 0.01)
				assert.Equal(t, expected.FullText, match.FullText)
			}
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
