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
	Secret string       `json:"secret,omitempty"`
	Rule   secrets.Rule `json:"rule"`
	Source Source       `json:"source"`
	Path   string       `json:"path,omitempty"`
}

func (f Finding) String() string {
	var lines []string
	if f.Secret != "" {
		lines = append(lines, fmt.Sprintf("Secret: %s", f.Secret))
	}
	lines = append(lines, fmt.Sprintf("Rule: %s", f.Rule))
	lines = append(lines, fmt.Sprintf("Source: %s", f.Source))
	if f.Path != "" {
		lines = append(lines, fmt.Sprintf("Path: %s", f.Path))
	}
	return strings.Join(lines, "\n")
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
