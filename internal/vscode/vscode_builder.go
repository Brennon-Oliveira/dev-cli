package vscode

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
)

type realVSCode struct {
	pather          pather.Pather
	executor        exec.Executor
	devcontainerCLI devcontainer.DevContainerCLI
}

type Option func(*realVSCode)

func NewVSCode(opts ...Option) VSCode {
	vs := &realVSCode{}

	for _, opt := range opts {
		opt(vs)
	}

	return vs
}

func WithPather(p pather.Pather) Option {
	return func(vs *realVSCode) {
		vs.pather = p
	}
}

func WithDevcontainerCLI(d devcontainer.DevContainerCLI) Option {
	return func(vs *realVSCode) {
		vs.devcontainerCLI = d
	}
}

func WithExecutor(e exec.Executor) Option {
	return func(vs *realVSCode) {
		vs.executor = e
	}
}
