package analysis

import (
	"encoding/json"
	"fmt"
	"github.com/bthuilot/dockerleaks/pkg/secrets"
	"strings"
)

type Source string

const (
	BuildArgument Source = "build-arg"
	EnvVar        Source = "env-var"
	File          Source = "file"
)

type Finding struct {
	Secret secrets.Secret `json:"secret"`
	Rule   secrets.Rule   `json:"rule"`
	Source Source         `json:"source"`
}

func (f Finding) String() string {
	return fmt.Sprintf(`Secret: %s
Rule: %s
Source: %s`, f.Secret, f.Rule, f.Source)
}

type Formatter func([]Finding) (string, error)

func DefaultFormatter(findings []Finding) (string, error) {
	formatted := make([]string, 0, len(findings))
	for _, f := range findings {
		formatted = append(formatted, f.String())
	}
	return strings.Join(formatted, "\n---\n"), nil
}

func JSONFormatter(findings []Finding) (string, error) {
	raw, err := json.MarshalIndent(findings, "", "  ")
	return string(raw), err
}
