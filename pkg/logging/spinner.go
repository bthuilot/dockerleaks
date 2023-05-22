package logging

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/bthuilot/dockerleaks/pkg/config"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func StartSpinner(msg string) (spnr *spinner.Spinner) {
	spnr = spinner.New(spinner.CharSets[69], time.Second/10)
	spnr.Prefix = fmt.Sprintf("%s ", msg)
	if !config.ShouldUseSpinner() {
		logrus.Info(msg)
		return
	}
	spnr.Start()
	return
}

func FinishSpinner(spnr *spinner.Spinner, finalMsg string) {
	if spnr == nil {
		logrus.Info(finalMsg)
		return
	}
	msg := strings.TrimSpace(spnr.Prefix)
	spnr.FinalMSG = fmt.Sprintf("%s %s\n", msg, finalMsg)
	if !config.ShouldUseSpinner() {
		logrus.Info(spnr.FinalMSG)
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

func FinishSpinnerWithError(spnr *spinner.Spinner, err error) {
	result := successText("complete")
	if err != nil {
		result = errorText("error")
	}
	FinishSpinner(spnr, result)
	if err != nil {
		Fatal(err)
	}
}
