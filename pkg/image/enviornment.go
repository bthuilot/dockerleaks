package image

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"strings"
)

func parseEnvVars(cli *client.Client, ctx context.Context, ref string) ([]EnvVar, error) {
	imageInspect, _, err := cli.ImageInspectWithRaw(ctx, ref)
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
