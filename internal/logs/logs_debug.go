//go:build dev

package logs

import "fmt"

func Debug(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(output, "\x1b[33m[DEBUG] %s\x1b[0m\n", msg)
}
