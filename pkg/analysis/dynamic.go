package analysis

import (
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/bthuilot/dockerleaks/pkg/secrets"
)

func Dynamic(img image.Image, detector secrets.Detector) ([]Finding, error) {
	return nil, nil
}
