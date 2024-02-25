package secrets

import (
	"github.com/sirupsen/logrus"
	"io"
)

// FileMatch represents a match of string that is detected to be a secret value
type FileMatch struct {
	// Rule is the rule that matches this string
	Rule DynamicRule
	// Secret is the actual value of the string
	Secret Secret
	// Path is the path of the file that was searched
	Path string
	//// StartPos is the starting position of the match
	//StartPos int
	//// EndPos is the ending position of the match
	//EndPos int
}

func findDynamicRuleMatches(path string, body io.Reader, rules []DynamicRule) (matches []FileMatch, err error) {
	// TODO: fix this to only read once for all rules, but dont read at all if not needed
	var content []byte
	if content, err = io.ReadAll(body); err != nil {
		return nil, err
	}

	for _, r := range rules {
		logrus.Debugf("checking rule %s with %s", r.Name, r.FilePattern)
		if r.FilePattern != nil && !r.FilePattern.MatchString(path) {
			continue
		}

		// If the pattern is nil, we only need to check if the file matches
		if r.Pattern == nil {
			matches = append(matches, FileMatch{
				Rule: r,
				Path: path,
			})
			continue
		}

		// If the pattern is not nil, we need to check the content of the file
		logrus.Debugf("checking content of file %s", path)
		logrus.Print(string(content))
		for _, s := range r.Pattern.FindSubmatch(content) {
			entropy := CalculateShannonEntropy(string(s))
			if entropy < r.MinEntropy {
				continue
			}
			matches = append(matches, FileMatch{
				Rule: r,
				Secret: Secret{
					Value:   string(s),
					Entropy: entropy,
				},
				Path: path,
			})
		}
	}
	return

}
