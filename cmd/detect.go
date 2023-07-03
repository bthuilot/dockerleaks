package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/bthuilot/dockerleaks/internal/config"
	"github.com/bthuilot/dockerleaks/pkg/common"
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/bthuilot/dockerleaks/pkg/layers"
	"github.com/bthuilot/dockerleaks/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO(make sub-command?)

func runDetect(cmd *cobra.Command, args []string) {
	var (
		cfg                   config.File
		spnr                  *spinner.Spinner
		secretStringsDetected []common.SecretString
	)

	// Parse image name from CLI argss
	imageName, err := cmd.Flags().GetString("image")
	if err != nil {
		_ = cmd.Help()
		logging.Fatal("you must supply the image to pull")
	}

	// Parse the configuration file and user supplied rules
	spnr = logging.StartSpinner("parsing configuration")
	if err = viper.Unmarshal(&cfg); err != nil {
		logging.Fatal(err)
	}

	logrus.Infof("parsing regular expression detection configuration")
	rules, invalidRules := config.ParseRules(cfg.Rules)
	logging.FinishSpinnerWithError(spnr, err)
	if len(invalidRules) > 0 && viper.GetBool("ignore-invalid") {
		for _, iR := range invalidRules {
			logrus.Debugf("invalid pattern '%s'", iR.Pattern)
		}
		logging.Msg("%d invalid rules found, ignoring due to flag `ignore-invalid`", len(invalidRules))
	}

	// Connect to docker daemon and pull image if necessary
	spnr = logging.StartSpinner("connecting to docker daemon")
	i, err := image.NewImage(imageName)
	logging.FinishSpinnerWithError(spnr, err)

	// TODO(refactor to be part of constructing)
	if viper.GetBool("pull_image") {
		spnr = logging.StartSpinner("pulling image from remote")
		err = i.Pull()
		logging.FinishSpinnerWithError(spnr, err)
	}

	// If checking layers are enabled, run layer detector
	if !cfg.Layer.Disable {
		spnr = logging.StartSpinner("checking layers for secrets")
		detector := layers.NewDetector(i).WithRules(rules...)
		if !cfg.ExcludeDefaultRules {
			detector = detector.UseDefaultRules()
		}
		secretStringsDetected, err = detector.Detect()
		logging.FinishSpinnerWithError(spnr, err)
	}

	// If checking filesystem are enabled, run filesystem detector
	if !cfg.Filesystem.Disable {
		spnr = logging.StartSpinner("checking filesystem for secrets")
		err = fmt.Errorf("not impleted yet")
		logging.FinishSpinnerWithError(spnr, err)
	}

	if len(secretStringsDetected) == 0 {
		logging.Header("no secret strings found", logging.H1)
	} else {
		logging.Header("secrets found", logging.H1)
		for _, s := range secretStringsDetected {
			logging.Msg("%s\n\n", s)
		}
	}
}
