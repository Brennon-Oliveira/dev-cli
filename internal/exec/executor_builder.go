package exec

import (
	"io"
	"os"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
)

type realExecutor struct {
	stdout io.Writer
	stdin  io.Reader
}

type Option func(*realExecutor)

func NewExecutor(opts ...Option) Executor {
	e := &realExecutor{
		stdout: logger.GetWriter(),
		stdin:  os.Stdin,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func WithStdout(w io.Writer) Option {
	return func(e *realExecutor) {
		e.stdout = w
	}
}

func WithStdin(r io.Reader) Option {
	return func(e *realExecutor) {
		e.stdin = r
	}
}
