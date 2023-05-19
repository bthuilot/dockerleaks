package detections

import (
	"fmt"
	"github.com/bthuilot/dockerleaks/pkg/image"
)

type Detector interface {
	Run(image image.Image) []Detection
}

type SecretSource = string

const (
	EnvVarSecret   SecretSource = "environment variable"
	BuildArgSecret              = "build argument"
	FileSecret                  = "file content"
	FileSystem                  = "generic secret"
)

type DetectionType = string

const (
	RegexDetection   DetectionType = "regular expression"
	EntropyDetection               = "entropy"
	FileDetection                  = "file"
)

type Detection struct {
	Type     DetectionType `json:"type"`
	Name     string        `json:"name"`
	Location string        `json:"location"`
	Value    string        `json:"value"`
	Source   SecretSource  `json:"source"`
}

/*
  regex expression detection: 'AWS Access Key' via environment variable
  location: ENV SECRET_KEY
  value: AWS
*/

func (d Detection) String() string {
	return fmt.Sprintf(
		"%s detection: '%s' via %s\nlocation: %s\nvalue: %s",
		d.Type, d.Name, d.Source, d.Location, d.Value,
	)
}
