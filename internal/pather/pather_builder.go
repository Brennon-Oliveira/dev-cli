package pather

import "github.com/Brennon-Oliveira/dev-cli/internal/exec"

type realPather struct {
	executor  exec.Executor
	lookupEnv LookupEnv
}

type LookupEnv func(key string) (string, bool)

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

func WithLookupEnv(l LookupEnv) Option {
	return func(p *realPather) {
		p.lookupEnv = l
	}
}
