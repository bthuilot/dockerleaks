package image

import (
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
	"io"
)

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
