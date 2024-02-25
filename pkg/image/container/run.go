package container

import (
	"errors"
	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"
	"io"
)

func (c container) Stop() error {
	if !c.running {
		return nil
	}

	if err := c.cli.ContainerStop(c.ctx, c.id, containerTypes.StopOptions{}); err != nil {
		logrus.Errorf("failure stopping container: %s", err)
		return errors.New("unable to stop container")
	}
	return nil
}

func (c container) Start() error {
	if c.running {
		return errors.New("container already running")
	}

	if err := c.cli.ContainerStop(c.ctx, c.id, containerTypes.StopOptions{}); err != nil {
		logrus.Errorf("failure stopping container: %s", err)
		return errors.New("unable to stop container")
	}
	return nil
}

func (c container) RunCommand(cmd string) (string, error) {
	if !c.running {
		return "", errors.New("container not running")
	}
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
