package logging

import (
	"fmt"
	"os"
)

func Fatal(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func Msg(format string, args ...any) {
	fmt.Printf(format, args...)
}
