package logging

import (
	"fmt"
	"os"
)

func Fatal(format string, a ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func Msg(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}
