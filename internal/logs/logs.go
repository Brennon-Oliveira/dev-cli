package logs

import (
	"fmt"
	"io"
	"os"
)

var (
	verbose bool
	output  io.Writer = os.Stdout
)

func SetVerbose(v bool) {
	verbose = v
}

func SetOutput(w io.Writer) {
	output = w
}

func Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(output, "\x1b[36m\x1b[1m➜\x1b[0m %s\n", msg)
}

func Verbose(format string, args ...any) {
	if !verbose {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(output, "\x1b[90m  │ %s\x1b[0m\n", msg)
}

func Success(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(output, "\x1b[32m✓\x1b[0m %s\n", msg)
}

func Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(output, "\x1b[33m⚠ %s\x1b[0m\n", msg)
}

func Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(output, "\x1b[31m✗ %s\x1b[0m\n", msg)
}
