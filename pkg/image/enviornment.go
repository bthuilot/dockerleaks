package image

import (
	"github.com/bthuilot/dockerleaks/internal/util"
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
	return uniqueEnvVars(vars), nil
}

func uniqueEnvVars(envVars []EnvVar) (unique []EnvVar) {
	uniqueVars := make(map[string][]string)
	for _, v := range envVars {
		if _, ok := uniqueVars[v.Name]; !ok {
			uniqueVars[v.Name] = []string{v.Value}
			unique = append(unique, v)
		} else if !util.Any(uniqueVars[v.Name], func(s string) bool {
			return s == v.Value
		}) {
			uniqueVars[v.Name] = append(uniqueVars[v.Name], v.Value)
			unique = append(unique, v)
		}
	}
	return unique
}
