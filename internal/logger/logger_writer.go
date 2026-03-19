package logger

import (
	"fmt"
	"io"

	loggerutils "github.com/Brennon-Oliveira/dev-cli/internal/logger/logger_utils"
)

type LoggerWriter struct {
	dest         io.Writer
	allowWriting bool
}

func NewLoggerWriter(dest io.Writer, allowWriting bool) *LoggerWriter {
	return &LoggerWriter{dest: dest, allowWriting: allowWriting}
}

func (sw *LoggerWriter) Write(p []byte) (n int, err error) {

	if !sw.allowWriting {
		return 0, nil
	}

	fmt.Fprintf(sw.dest, "%s  │ ", loggerutils.HighIntensityBlackColor)
	sw.dest.Write(p)
	fmt.Fprintf(sw.dest, "%s", loggerutils.ResetColor)

	return len(p), nil
}

func (sw *LoggerWriter) SetAllowWriting(allowWriting bool) {
	sw.allowWriting = allowWriting
}
