package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type realLogger struct {
	verbose bool
	output  io.Writer
}

func NewLogger() *realLogger {
	return &realLogger{
		output: os.Stdout,
	}
}

func InitLogger(customLogger Logger) {
	if customLogger != nil {
		logger = customLogger
		return
	}
	logger = NewLogger()
}

func (l *realLogger) SetVerbose(verbose bool) {
	l.verbose = verbose
}

func SetVerbose(verbose bool) {
	logger.SetVerbose(verbose)
}

func (l *realLogger) SetOutput(writer io.Writer) {
	l.output = writer
}

func SetOutput(writer io.Writer) {
	logger.SetOutput(writer)
}

func (l *realLogger) Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "\x1b[36m\x1b[1m➜\x1b[0m %s\n", msg)
}

func Info(format string, args ...any) {
	logger.Info(format, args)
}

func (l *realLogger) Verbose(format string, args ...any) {
	if !l.verbose {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "\x1b[90m  │ %s\x1b[0m\n", msg)
}

func Verbose(format string, args ...any) {
	logger.Verbose(format, args...)
}

func (l *realLogger) VerboseFromOutput(output *io.Reader) {
	if !l.verbose {
		return
	}
	bytes := bytes.Buffer{}
	bytes.WriteString("\x1b[90m  │ ")
	io.Copy(&bytes, *output)
	bytes.WriteString("\x1b[0m\n")
	l.output.Write(bytes.Bytes())
}

func VerboseFromOutput(output *io.Reader) {
	logger.VerboseFromOutput(output)
}

func (l *realLogger) Success(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "\x1b[32m✓\x1b[0m %s\n", msg)
}

func Success(format string, args ...any) {
	logger.Success(format, args...)
}

func (l *realLogger) Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "\x1b[33m⚠ %s\x1b[0m\n", msg)
}

func Warn(format string, args ...any) {
	logger.Warn(format, args...)
}

func (l *realLogger) Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "\x1b[31m✗ %s\x1b[0m\n", msg)
}

func Error(format string, args ...any) {
	logger.Error(format, args...)
}
