package image

// EnvVar represents a set environment variable
// during build time
type EnvVar struct {
	// Name is the name of the environment variable
	Name string
	// Value is the full value of the variable
	Value string
	// Location is the line it was defined
	Location string
}

// BuildArg represents a build argument that has
// be parsed out of a RUN step.
type BuildArg struct {
	// Name is the name of the build argument
	Name string
	// Value is the value of the build argument
	Value string
	// Location is the RUN line in which the build argument was found
	Location string
}
