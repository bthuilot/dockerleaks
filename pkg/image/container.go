package image

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"io"
)

type Container interface {
	// Stop will stop the container
	Stop() error
	// RunCommand will run a command in the container and return the output
	RunCommand(cmd string) (string, error)
}

type c struct {
	id  string
	cli *client.Client
	ctx context.Context
}

func (c c) RunCommand(cmd string) (string, error) {
	// exec command
	exec, err := c.cli.ContainerExecCreate(c.ctx, c.id, types.ExecConfig{
		Cmd: []string{"sh", "-c", cmd},
	})
	if err != nil {
		return "", err
	}
	// pull output
	resp, err := c.cli.ContainerExecAttach(c.ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer resp.Close()
	// read output
	output, err := io.ReadAll(resp.Reader)
	return string(output), err
}

func (c c) Stop() error {
	if err := c.cli.ContainerStop(c.ctx, c.id, container.StopOptions{}); err != nil {
		logrus.Errorf("failure stopping container: %s", err)
		return errors.New("unable to stop container")
	}
	return nil
}

func (i image) StartContainer() (Container, error) {
	// create container
	resp, err := i.cli.ContainerCreate(i.ctx, &container.Config{
		Image: i.ref.String(),
	}, nil, nil, nil, fmt.Sprintf("dockerleaks-scan-%s", i.ref.String()))
	if err != nil {
		logrus.Errorf("failure creating container: %s", err)
		return nil, errors.New("unable to create container")
	}

	// start container
	if err = i.cli.ContainerStart(i.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		logrus.Errorf("failure starting container: %s", err)
		return nil, errors.New("unable to start container")
	}

	return c{
		id: resp.ID,
	}, nil
}
