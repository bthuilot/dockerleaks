package image

// EnvVar represents a set environment variable
// during build time
type EnvVar struct {
	Name     string
	Value    string
	Location string
}

// BuildArg represents a build argument that has
// be parsed out of a RUN step.
type BuildArg struct {
	Name     string
	Value    string
	Location string
}
