package logger

import "io"

type Logger interface {
	SetVerbose(v bool)
	SetOutput(w io.Writer)
	Info(format string, args ...any)
	Verbose(format string, args ...any)
	VerboseFromOutput(output *io.Reader)
	Success(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
}

var logger Logger
