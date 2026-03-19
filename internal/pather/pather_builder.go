package pather

import "github.com/Brennon-Oliveira/dev-cli/internal/exec"

type realPather struct {
	executor exec.Executor
}

type Option func(*realPather)

func NewPather(opts ...Option) *realPather {
	p := &realPather{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func WithExecutor(e exec.Executor) Option {
	return func(p *realPather) {
		p.executor = e
	}
}
