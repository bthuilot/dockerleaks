package cmd

import (
	"github.com/bthuilot/dockerleaks/pkg/config"
	"github.com/bthuilot/dockerleaks/pkg/detections"
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/bthuilot/dockerleaks/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO(make sub-command?)

func runDetect(cmd *cobra.Command, args []string) {
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logging.Fatal(err)
	}

	imageName, err := cmd.Flags().GetString("image")
	if err != nil {
		_ = cmd.Help()
		logging.Fatal("you must supply the image to pull")
	}

	spnr := logging.StartSpinner("parsing config")
	regex, err := detections.NewRegexDetector(cfg.Regexp.Patterns)
	logging.FinishSpinnerWithError(spnr, err)

	spnr = logging.StartSpinner("retrieving image information")
	i, err := image.NewImage(imageName, viper.GetBool("pull_image"))
	logging.FinishSpinnerWithError(spnr, err)

	spnr = logging.StartSpinner("running regex detections")
	found := regex.Run(i)
	logging.FinishSpinnerWithError(spnr, err)

	spnr = logging.StartSpinner("running entropy detections")
	logging.FinishSpinner(spnr, "not implemented, skipping")

	spnr = logging.StartSpinner("running filesystem detections")
	logging.FinishSpinner(spnr, "not implemented, skipping")

	logging.Header("secrets found", logging.H1)
	for _, d := range found {
		logging.Msg("%s\n\n", d)
	}
}
