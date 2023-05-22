package image

import (
	"github.com/sirupsen/logrus"
	"strings"
)

// ParseEnvVars will parse out environment variables by
// inspecting the image and pull out each environment
// and split on the first '=' into the name and value
func (i image) ParseEnvVars() ([]EnvVar, error) {
	imageInspect, _, err := i.cli.ImageInspectWithRaw(i.ctx, i.ref.String())
	if err != nil {
		return nil, err
	}
	var vars []EnvVar
	for _, env := range imageInspect.Config.Env {
		splitEnv := strings.SplitN(env, "=", 2)
		if len(splitEnv) != 2 {
			logrus.Warnf("skipping invalid env %s", env)
			continue
		}
		vars = append(vars, EnvVar{
			Name:     splitEnv[0],
			Value:    splitEnv[1],
			Location: env,
		})
	}
	return vars, nil
}
