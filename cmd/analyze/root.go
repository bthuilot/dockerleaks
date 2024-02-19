package analyze

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/bthuilot/dockerleaks/internal/config"
	"github.com/bthuilot/dockerleaks/pkg/analysis"
	"github.com/bthuilot/dockerleaks/pkg/image"
	"github.com/bthuilot/dockerleaks/pkg/logging"
	"github.com/bthuilot/dockerleaks/pkg/secrets"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type contextKey string

const (
	imageContextKey    contextKey = "dockerleaks-docker-image"
	detectorContextKey contextKey = "dockerleaks-secret-detector"
	findingsContextKey contextKey = "dockerleaks-findings"
)

const errorMsgFmt = `!! ERROR: %s !!
This is most like an error with dockerleaks and should be reported
https://github.com/bthuilot/dockerleaks/issues/new/choose`

var Command = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze an image for secrets",
	Long:  `Analyze an image for secrets, either statically or dynamically.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var (
			cfg  config.File
			spnr *spinner.Spinner
			ctx  = context.Background()
		)

		// Parse image name from CLI argss
		imageName, _ := cmd.Flags().GetString("image")

		// Parse the configuration file and user supplied rules
		spnr = logging.StartSpinner("parsing configuration")
		err := viper.Unmarshal(&cfg)
		logrus.Infof("parsing regular expression detection configuration")
		rules, invalidRules := secrets.ParseRules(cfg.Rules)

		logging.FinishSpinnerWithError(spnr, err)
		if len(invalidRules) > 0 && viper.GetBool("ignore-invalid") {
			for _, iR := range invalidRules {
				logrus.Debugf("invalid pattern '%s'", iR.Pattern)
			}
			logging.Msg("%d invalid rules found, ignoring due to flag `ignore-invalid`", len(invalidRules))
		}

		detector := secrets.NewDetector(secrets.Opts{UseDefaultRules: !cfg.ExcludeDefaultRules}, rules...)
		ctx = context.WithValue(ctx, detectorContextKey, detector)
		// Connect to docker daemon and pull image if necessary
		spnr = logging.StartSpinner("connecting to docker daemon")
		i, err := image.NewImage(imageName)
		logging.FinishSpinnerWithError(spnr, err)

		if pull, _ := cmd.Flags().GetBool("pull"); pull {
			spnr = logging.StartSpinner("pulling image from remote")
			err = i.Pull()
			logging.FinishSpinnerWithError(spnr, err)
		}
		ctx = context.WithValue(ctx, imageContextKey, i)
		cmd.SetContext(ctx)
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Retrieve the context from the command
		ctx := cmd.Context()

		// Retrieve your data from the context
		findings, ok := ctx.Value(findingsContextKey).([]analysis.Finding)
		if !ok {
			logging.Fatal(errorMsgFmt, "error parsing findings from context")
		}

		var formatter analysis.Formatter
		switch format, _ := cmd.Flags().GetString("output"); format {
		case "json":
			formatter = analysis.JSONFormatter
		default:
			formatter = analysis.DefaultFormatter
		}

		if len(findings) == 0 {
			logging.Header("no secret strings found", logging.H1)
		} else {
			logging.Header(fmt.Sprintf("%d secrets found", len(findings)), logging.H1)
			logging.Msg("")
			output, err := formatter(findings)
			if err != nil {
				logrus.Errorf("error formatting findings: %s", err)
				logging.Fatal(errorMsgFmt, "error formatting findings")
			}
			logging.Msg(output)
		}
	},
}

func init() {
	Command.PersistentFlags().StringP("image", "i", "", "the name of the image")
	if err := Command.MarkPersistentFlagRequired("image"); err != nil {
		logging.Fatal(err.Error())
	}

	Command.PersistentFlags().BoolP("pull", "p", false, "image should be pulled from remote")

	Command.PersistentFlags().StringP("output", "o", "text", "output format (text, json)")

	Command.AddCommand(static, dynamic)
}

// parseContext will parse the context and return the parsed [image.Image] and [secrets.Detector]
// set by the [Command] PersistentPreRun hook. If the context is not set for both, the program will exit.
func parseContext(ctx context.Context) (image.Image, secrets.Detector) {
	img, ok := ctx.Value(imageContextKey).(image.Image)
	if !ok {
		logging.Fatal(errorMsgFmt, "error parsing image context")
	}

	detector, ok := ctx.Value(detectorContextKey).(secrets.Detector)
	if !ok {
		logging.Fatal(errorMsgFmt, "error parsing detector context")
	}

	return img, detector
}
