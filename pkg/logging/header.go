package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type HeaderSize = uint8

const (
	H1 HeaderSize = iota
	H2
	H3
)

func Header(text string, s HeaderSize) {
	var header string
	switch s {
	case H1:
		header = h1Display(text)
	default:
		logrus.Warn("invalid header size given, displaying default")
		header = fmt.Sprintf("# %s", header)
	}
	Msg(header)
}

func h1Display(text string) string {
	header := fmt.Sprintf("# %s #", text)
	return fmt.Sprintf(
		"\n%s\n%s\n%s\n",
		strings.Repeat("#", len(header)),
		header,
		strings.Repeat("#", len(header)),
	)
}
