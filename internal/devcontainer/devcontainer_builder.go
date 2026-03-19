package devcontainer

import "github.com/Brennon-Oliveira/dev-cli/internal/exec"

type realDevContainerCLI struct {
	executor exec.Executor
}

type Option func(*realDevContainerCLI)

func NewDevContainerCLI(opts ...Option) *realDevContainerCLI {
	d := &realDevContainerCLI{}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func WithExecutor(e exec.Executor) Option {
	return func(d *realDevContainerCLI) {
		d.executor = e
	}
}
