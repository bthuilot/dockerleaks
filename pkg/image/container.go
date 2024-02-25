package image

import (
	"errors"
	"fmt"
	"github.com/bthuilot/dockerleaks/pkg/image/container"
	"github.com/docker/distribution/uuid"
	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"
)

func (i image) CreateContainer() (container.Container, error) {
	// create container
	name := fmt.Sprintf("dockerleaks-scan-%s", uuid.Generate())
	resp, err := i.cli.ContainerCreate(i.ctx, &containerTypes.Config{
		Image: i.ref.String(),
	}, nil, nil, nil, name)
	if err != nil {
		logrus.Errorf("failure creating container: %s", err)
		return nil, errors.New("unable to create container")
	}

	// start container
	if err = i.cli.ContainerStart(i.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		logrus.Errorf("failure starting container: %s", err)
		return nil, errors.New("unable to start container")
	}

	return container.New(
		resp.ID,
		name,
		i.cli,
		i.ctx,
	)
}

func (i image) DestroyContainer(c container.Container) error {
	if err := c.Stop(); err != nil {
		logrus.Errorf("failure stopping container: %s", err)
		return err
	}

	if err := i.cli.ContainerRemove(i.ctx, c.ID(), types.ContainerRemoveOptions{}); err != nil {
		logrus.Errorf("failure removing container: %s", err)
		return errors.New("unable to remove container")
	}

	return nil
}
