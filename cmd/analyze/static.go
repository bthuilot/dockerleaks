package analyze

import (
	"context"
	"github.com/bthuilot/dockerleaks/pkg/analysis"
	"github.com/bthuilot/dockerleaks/pkg/logging"
	"github.com/spf13/cobra"
)

var static = &cobra.Command{
	Use:   "static",
	Short: "Static analyze an image for secrets",
	Long:  `Analyze a built docker image by inspect contents of layer commands`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		img, detector := parseContext(ctx)
		spnr := logging.StartSpinner("beginning static analysis...")
		findings, err := analysis.Static(img, detector)
		ctx = context.WithValue(ctx, findingsContextKey, findings)
		cmd.SetContext(ctx)
		logging.FinishSpinnerWithError(spnr, err)
	},
}
