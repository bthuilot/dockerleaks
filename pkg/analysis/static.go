package analysis

import (
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/bthuilot/dockerleaks/pkg/secrets"
	"github.com/sirupsen/logrus"
)

func Static(img image.Image, detector secrets.Detector) (findings []Finding, err error) {
	var (
		envVars   []image.EnvVar
		buildArgs []image.BuildArg
		matches   []secrets.TextMatch
	)

	envVars, err = img.ParseEnvVars()
	if err != nil {
		return nil, err
	}

	for _, v := range envVars {
		logrus.Debugf("searching for secrets in env var %s", v.Name)
		matches, err = detector.SearchText(v.Value)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			findings = append(findings, Finding{
				Secret: m.Secret.String(),
				Rule:   m.Rule,
				Source: EnvVar,
			})
		}
	}

	buildArgs, err = img.ParseBuildArguments()
	if err != nil {
		return nil, err
	}

	for _, v := range buildArgs {
		matches, err = detector.SearchText(v.Value)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			findings = append(findings, Finding{
				Secret: m.Secret.String(),
				Rule:   m.Rule,
				Source: BuildArgument,
			})
		}
	}
	return
}
