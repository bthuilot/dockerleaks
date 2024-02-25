package analysis

import (
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/bthuilot/dockerleaks/pkg/secrets"
	"github.com/sirupsen/logrus"
	"io"
)

func Dynamic(img image.Image, detector secrets.Detector) ([]Finding, error) {
	// create a container from the image
	logrus.Infof("creating container from image")
	container, err := img.CreateContainer()
	if err != nil {
		return nil, err
	}
	defer func() {
		if destroyErr := img.DestroyContainer(container); destroyErr != nil {
			logrus.Errorf("failure destroying container: %s", destroyErr)
		}
	}()

	// export the container filesystem
	logrus.Infof("exporting container filesystem")
	fs, err := container.Export()
	if err != nil {
		return nil, err
	}

	var findings []Finding
	for {
		hdr, err := fs.Next()
		if err == io.EOF {
			logrus.Debugf("end of filesystem")
			break
		}
		if err != nil {
			logrus.Errorf("error reading filesystem: %s", err)
			return nil, err
		}
		logrus.Debugf("checking file %s", hdr.Name)
		matches, err := detector.SearchFile(hdr.Name, io.LimitReader(fs, hdr.Size))
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			findings = append(findings, Finding{
				Secret: m.Secret.String(),
				Rule:   m.Rule,
				Source: File,
				Path:   hdr.Name,
			})
		}
	}
	return findings, nil
}
