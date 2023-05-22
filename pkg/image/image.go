package image

import (
	"context"
	"errors"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// Image represents a built docker image
type Image interface {
	// ParseEnvVars will search an images configuration for all
	// set environment variables and return a list of EnvVar representing
	// the set environment variables
	ParseEnvVars() ([]EnvVar, error)
	// ParseBuildArguments will parse an images history for all set build args
	// during run commands and return a list of BuildArg representing the
	// discovered build arguments
	ParseBuildArguments() ([]BuildArg, error)
	// Pull will pull down an image from remote
	Pull() error
	//ParseFS() (fs.FS, error)
}

// image is the concrete implementation of the Image interface
type image struct {
	// ref is the parsed image reference supplied by the user
	ref reference.Reference
	// cli is the moby client.Client for interacting with docker
	cli *client.Client
	// ctx is the global context for all client interactions
	ctx context.Context
}

// NewImage connects to the docker daemon and constructs a new Image from its reference
func NewImage(name string) (Image, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("failure constructing new client: %s", err)
		return nil, errors.New("could not construct client")
	}

	ctx := context.Background()

	ref, err := reference.ParseAnyReference(name)
	if err != nil {
		logrus.Errorf("failure parsing docker name: %s", err)
		return nil, errors.New("invalid docker name")
	}

	return image{
		ref: ref,
		cli: cli,
		ctx: ctx,
	}, err
}
