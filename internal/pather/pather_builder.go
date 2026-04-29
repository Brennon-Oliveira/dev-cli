package pather

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/env"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
)

type realPather struct {
	executor  exec.Executor
	lookupEnv env.LookupEnvFunc
}

type Option func(*realPather)

func NewPather(opts ...Option) *realPather {
	p := &realPather{
		lookupEnv: env.LookupEnv,
	}

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

func WithLookupEnv(l env.LookupEnvFunc) Option {
	return func(p *realPather) {
		p.lookupEnv = l
	}
}
