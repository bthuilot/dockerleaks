package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/bthuilot/dockerleaks/internal/config"
	"github.com/bthuilot/dockerleaks/pkg/detections"
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/bthuilot/dockerleaks/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"regexp"
)

// TODO(make sub-command?)

func constructRegexDetector(cfg config.Config) (detections.Detector, error) {
	var patterns []detections.Pattern
	for _, p := range cfg.Regexp.Patterns {
		regex, err := regexp.Compile(p.Expression)
		if err != nil {
			logrus.Errorf(
				"unable to parse regular expression %s `%s`: %s",
				p.Name, p.Expression, err,
			)
			continue
		}
		patterns = append(patterns, detections.Pattern{
			RegExp: regex,
			Name:   p.Name,
		})
	}

	if !cfg.Regexp.DisableDefaults {
		patterns = append(patterns, detections.DefaultPatterns...)
	}

	return detections.NewRegexDetector(patterns)
}

func runDetect(cmd *cobra.Command, args []string) {
	var (
		cfg      config.Config
		spnr     *spinner.Spinner
		detected []detections.Detection
	)
	if err := viper.Unmarshal(&cfg); err != nil {
		logging.Fatal(err)
	}

	imageName, err := cmd.Flags().GetString("image")
	if err != nil {
		_ = cmd.Help()
		logging.Fatal("you must supply the image to pull")
	}

	spnr = logging.StartSpinner("parsing configurations")
	logrus.Infof("parsing regular expression detection configuration")
	regexDetector, err := constructRegexDetector(cfg)
	if err != nil {
		logging.FinishSpinnerWithError(spnr, err)
	}

	logrus.Infof("parsing entropy detection configuration")
	// TODO(entropy detector)
	logrus.Infof("parsing file detection configuration")
	// TODO(file detector)

	spnr = logging.StartSpinner("connecting to docker daemon")
	i, err := image.NewImage(imageName)
	logging.FinishSpinnerWithError(spnr, err)

	if viper.GetBool("pull_image") {
		spnr = logging.StartSpinner("pulling image from remote")
		err = i.Pull()
		logging.FinishSpinnerWithError(spnr, err)
	}

	spnr = logging.StartSpinner("parsing environment variables")
	envVars, err := i.ParseEnvVars()
	logging.FinishSpinnerWithError(spnr, err)

	spnr = logging.StartSpinner("parsing build arguments")
	buildArgs, err := i.ParseBuildArguments()
	logging.FinishSpinnerWithError(spnr, err)

	for _, d := range []detections.Detector{
		regexDetector,
	} {
		spnr = logging.StartSpinner(
			fmt.Sprintf("running detection '%s' on environment variables", d),
		)
		detected = append(detected, d.EvalEnvVars(envVars)...)
		logging.FinishSpinnerWithError(spnr, err)

		spnr = logging.StartSpinner(
			fmt.Sprintf("running detection '%s' on build arguments", d),
		)
		detected = append(detected, d.EvalBuildArgs(buildArgs)...)
		logging.FinishSpinnerWithError(spnr, err)
	}

	if len(detected) == 0 {
		logging.Header("no secrets found", logging.H1)
	} else {
		logging.Header("secrets found", logging.H1)
		for _, d := range detected {
			logging.Msg("%s\n\n", d)
		}
	}
}
