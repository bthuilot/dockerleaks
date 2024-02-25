package container

import (
	"archive/tar"
	"context"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type Container interface {
	// Start will start the container
	Start() error
	// Stop will stop the container
	Stop() error
	// ID will return the container ID
	ID() string

	// RunCommand will run a command in the container and return the output
	RunCommand(cmd string) (string, error)

	// Export will export the container filesystem to a tarball
	Export() (*tar.Reader, error)
}

type container struct {
	id      string
	name    string
	running bool
	cli     *client.Client
	ctx     context.Context
}

func New(id string, name string, cli *client.Client, ctx context.Context) (Container, error) {
	logrus.Debugf("creating new container: %s", id)
	return container{
		id:      id,
		name:    name,
		cli:     cli,
		ctx:     ctx,
		running: false,
	}, nil
}

func (c container) ID() string {
	return c.id
}

func (c container) Export() (*tar.Reader, error) {
	// export container filesystem
	resp, err := c.cli.ContainerExport(c.ctx, c.id)
	if err != nil {
		return nil, err
	}
	return tar.NewReader(resp), nil
}
