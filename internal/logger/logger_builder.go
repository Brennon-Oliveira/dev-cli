package logger

import (
	"io"
	"os"
)

type realLogger struct {
	verbose bool
	w       *LoggerWriter
}

type Option func(*realLogger)

func NewLogger(opts ...Option) *realLogger {
	w := NewLoggerWriter(os.Stdout, false)
	l := &realLogger{
		w: w,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func InitLogger(opts ...Option) {
	logger = NewLogger(opts...)
}

func WithWriter(w io.Writer) Option {
	return func(l *realLogger) {
		l.w = NewLoggerWriter(w, l.verbose)
	}
}

func WithVerbose(v bool) Option {
	return func(l *realLogger) {
		l.verbose = v
		l.w.SetAllowWriting(v)
	}
}
