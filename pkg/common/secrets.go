package common

import "fmt"

type SecretString struct {
	Name string
	// Location is the line in the docker file or
	// filesystem path where the secret was found
	Location string `json:"location"`
	// Value is the actual value of the secret
	Value string `json:"value"`
	// Source is the SecretSource of where the secret originated from.
	Source SecretSource `json:"source"`
}

// SecretSource represents the source of where a secret was detected
type SecretSource = string

const (
	// EnvVar is a SecretSource for secret detected in environment variables
	EnvVar SecretSource = "environment variable"
	// BuildArgument is a SecretSource for secret detected in build argument
	BuildArgument SecretSource = "build argument"
	// File is a SecretSource for secret detected in file content
	File SecretSource = "file content"
)

/* EXAMPLE OUTPUT:
regex expression detection: 'AWS Access Key' via environment variable
location: ENV SECRET_KEY
value: AWS
*/

// String is the human-readable output of the detection
// TODO(support JSON and YAML outputs)
func (d SecretString) String() string {
	return fmt.Sprintf(
		"detection: '%s' via %s\nlocation: %s\nvalue: %s",
		d.Name, d.Source, d.Location, d.Value,
	)
}
