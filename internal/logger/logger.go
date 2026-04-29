package logger

import "io"

type Logger interface {
	SetVerbose(v bool)
	SetOutput(w io.Writer)
	Info(format string, args ...any)
	Verbose(format string, args ...any)
	VerboseFromOutput(output io.Reader)
	Debug(format string, args ...any)
	Success(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
	GetWriter() io.Writer
}

var logger Logger
