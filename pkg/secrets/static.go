package secrets

// TextMatch represents a match of string that is detected to be a secret value
type TextMatch struct {
	// Rule is the rule that matches this string
	Rule StaticRule
	// Secret is the actual value of the string
	Secret Secret
	// FullText is the full text that was searches
	FullText string
	// StartPos is the starting position of the match
	StartPos int
	// EndPos is the ending position of the match
	EndPos int
}

// findRuleMatches will search a string's content for any matches to
// the list of SecretStringRules provided
func findStaticRuleMatches(content string, rules []StaticRule) (matches []TextMatch, err error) {
	for _, r := range rules {
		for _, s := range r.Pattern.FindStringSubmatch(content) {
			entropy := CalculateShannonEntropy(s)
			if entropy < r.MinEntropy {
				continue
			}
			matches = append(matches, TextMatch{
				Rule: r,
				Secret: Secret{
					Value:   s,
					Entropy: entropy,
				},
			})
		}
	}
	return
}
