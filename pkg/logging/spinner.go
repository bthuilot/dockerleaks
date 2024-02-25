package logging

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/bthuilot/dockerleaks/internal/config"
	"github.com/fatih/color"
	"os"
	"strings"
	"time"
)

func StartSpinner(msg string) *spinner.Spinner {
	prefix := fmt.Sprintf("%s %s ", grayText("=>"), msg)
	if !config.ShouldUseSpinner() {
		Msg(prefix + "\n")
		return nil
	}
	spnr := spinner.New(spinner.CharSets[69], time.Second/10,
		spinner.WithWriterFile(os.Stderr),
	)
	spnr.Prefix = prefix
	spnr.Start()
	return spnr
}

func FinishSpinner(spnr *spinner.Spinner, finalMsg string) {
	if spnr == nil {
		Msg(finalMsg)
		return
	}
	msg := strings.TrimSpace(spnr.Prefix)
	spnr.FinalMSG = fmt.Sprintf("%s %s", msg, finalMsg)
	if !config.ShouldUseSpinner() {
		Msg(finalMsg)
	}
	spnr.Stop()
}

func errorText(text string) string {
	if config.ShouldUseColor() {
		return color.New(color.FgRed).SprintFunc()(text)
	}
	return text
}

func successText(text string) string {
	if config.ShouldUseColor() {
		return color.New(color.FgGreen).SprintFunc()(text)
	}
	return text
}

func grayText(text string) string {
	if config.ShouldUseColor() {
		return color.New(color.FgHiBlack).SprintFunc()(text)
	}
	return text
}

func FinishSpinnerWithError(spnr *spinner.Spinner, err error) {
	result := successText("complete\n")
	if err != nil {
		result = errorText("error\n")
	}
	FinishSpinner(spnr, result)
	if err != nil {
		Fatal("%s: %s", errorText("ERROR"), err.Error())
	}
}
