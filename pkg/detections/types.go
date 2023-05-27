package detections

import (
	"fmt"
	"github.com/bthuilot/dockerleaks/pkg/image"
)

type Detector interface {
	// EvalBuildArgs will attempt to detect any
	// secrets in the given build arguments of an image.
	// It will return a list of Detection representing the
	// detected secrets found
	EvalBuildArgs([]image.BuildArg) []Detection
	// EvalEnvVars will attempt to detect any
	// secrets in the given environment variables of an image.
	// It will return a list of Detection representing the
	// detected secrets found
	EvalEnvVars([]image.EnvVar) []Detection

	//EvalFileSystem(filesystem fs.FS) []Detection

	// String returns the formatted name of the detector
	String() string
}

// SecretSource represents the source of where the
// secret was found within the image
type SecretSource = string

const (
	// EnvVarSecret is a secret that was found in an
	// environment variable
	EnvVarSecret SecretSource = "environment variable"
	// BuildArgSecret is a secret that was found within
	// a supplied build argument
	BuildArgSecret = "build argument"
	// FileSecret is a secret that was found within the
	// contents of a file
	FileSecret = "file content"
	// FileSystem is a secret that is just an entire file,
	// identified by its path or name. (e.g. terraform.tfstate)
	FileSystem = "file path"
)

// DetectionType is the method by which the
// secret was found
type DetectionType = string

const (
	// RegexDetection is a detection that identifies secrets using a list
	// of regular expressions
	RegexDetection DetectionType = "regular expression"
	// EntropyDetection is a detection that identifies secrets by
	// calculating the entropy of a string, and checking if that entropy is greater
	// than a given threshold
	EntropyDetection = "entropy"
	// FileDetection is a detection that identified a secret by the path or
	// name of a file. (e.g. terraform.tfstate)
	FileDetection = "file"
)

// Detection represents a detected secret
type Detection struct {
	// Type is the DetectionType of this secret, or rather
	// how this secret was detected
	Type DetectionType `json:"type"`
	// Name is the name of the secret, i.e. what does
	// the secret itself belong to/represent
	// i.e. AWS Access Token, GitLab API Key
	Name string `json:"name"`
	// Location is the line in the docker file or
	// filesystem path where the secret was found
	Location string `json:"location"`
	// Value is the actual value of the secret
	Value string `json:"value"`
	// Source is the SecretSource of where the secret originated from.
	Source SecretSource `json:"source"`
}

/* EXAMPLE OUTPUT:
regex expression detection: 'AWS Access Key' via environment variable
location: ENV SECRET_KEY
value: AWS
*/

// String is the human-readable output of the detection
// TODO(support JSON and YAML outputs)
func (d Detection) String() string {
	return fmt.Sprintf(
		"%s detection: '%s' via %s\nlocation: %s\nvalue: %s",
		d.Type, d.Name, d.Source, d.Location, d.Value,
	)
}
