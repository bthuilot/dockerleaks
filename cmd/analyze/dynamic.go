package analyze

import (
	"context"
	"github.com/bthuilot/dockerleaks/pkg/analysis"
	"github.com/bthuilot/dockerleaks/pkg/logging"
	"github.com/spf13/cobra"
)

var dynamic = &cobra.Command{
	Use:   "dynamic",
	Short: "Analyze an image for secrets dynamically",
	Long:  `Analyze a built docker image by starting the container and inspecting file contents`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		img, detector := parseContext(ctx) // will exit if error

		spnr := logging.StartSpinner("beginning dynamic analysis...")
		findings, err := analysis.Dynamic(img, detector)
		logging.FinishSpinnerWithError(spnr, err) // Exit if error

		ctx = context.WithValue(ctx, findingsContextKey, findings)
		cmd.SetContext(ctx)
	},
	Args: cobra.MinimumNArgs(1), // Require at least one subcommand
}
