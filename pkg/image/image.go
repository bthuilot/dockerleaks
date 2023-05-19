package image

import (
	"context"
	"errors"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type Image interface {
	GetEnvVars() ([]EnvVar, bool)
	GetBuildArgs() ([]BuildArg, bool)
	//ParseFS() (fs.FS, error)
}

type image struct {
	// Global
	ref string
	cli *client.Client
	ctx context.Context

	// EnvVars
	envVars      []EnvVar
	parseEnvVars bool

	// buildArgs
	buildArgs      []BuildArg
	parseBuildArgs bool
}

func (i image) GetEnvVars() ([]EnvVar, bool) {
	return i.envVars, i.parseEnvVars
}

func (i image) GetBuildArgs() ([]BuildArg, bool) {
	return i.buildArgs, i.parseBuildArgs
}

func NewImage(name string, pull bool) (Image, error) {
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

	if _, err = cli.ImagePull(ctx, ref.String(), types.ImagePullOptions{All: true}); pull && err != nil {
		logrus.Errorf("failure pulling docker image '%s': %s", ref, err)
		return nil, errors.New("unable to pull docker image")
	}

	envVars, err := parseEnvVars(cli, ctx, name)
	if err != nil {
		logrus.Errorf("failure parsing environment variables: %s", err)
		return nil, errors.New("unable to parse environment variables")
	}

	buildArgVars, err := parseBuildArgs(cli, ctx, ref.String())
	if err != nil {
		logrus.Errorf("failure parsing build arguments: %s", err)
		return nil, errors.New("unable to parse build arguments")
	}

	return image{
		ref: ref.String(),
		cli: cli,
		ctx: ctx,

		buildArgs:      buildArgVars,
		parseBuildArgs: true,
		envVars:        envVars,
		parseEnvVars:   true,
	}, err
}
