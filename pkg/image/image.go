package image

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bthuilot/dockerleaks/pkg/image/container"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"io"
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

	// CreateContainer will start a container from this image
	CreateContainer() (container.Container, error)

	// DestroyContainer will remove a container from the docker daemon
	DestroyContainer(container.Container) error
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

	if _, err = cli.Ping(ctx); err != nil {
		return nil, errors.New("docker daemon is not running")
	}

	return image{
		ref: ref,
		cli: cli,
		ctx: ctx,
	}, err
}

// Pull will pull down the image from remote.
func (i image) Pull() error {
	reader, err := i.cli.ImagePull(i.ctx, i.ref.String(), types.ImagePullOptions{})
	if err != nil {
		logrus.Errorf("failure pulling docker image '%s': %s", i.ref.String(), err)
		return errors.New("unable to pull docker image")
	}

	decoder := json.NewDecoder(reader)
	defer reader.Close()
	for {
		var message interface{}
		if err = decoder.Decode(&message); errors.Is(err, io.EOF) {
			logrus.Debug("end of image pull")
			break
		} else if err != nil {
			logrus.Errorf("error while decoding status from pull: %s", err)
			return errors.New("unable to pull image")
		}
		// TODO(refine to display status only when complete a layer using info)
		logrus.Debug(message)
	}
	return nil
}

func (i image) Delete() error {
	_, err := i.cli.ImageRemove(i.ctx, i.ref.String(), types.ImageRemoveOptions{})
	if err != nil {
		logrus.Errorf("failure removing docker image '%s': %s", i.ref.String(), err)
		return errors.New("unable to remove docker image")
	}
	return nil
}
