package container

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container/container_utils"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
)

type realContainerCLI struct {
	executor                exec.Executor
	config                  config.Config
	parseContainerOutput    container_utils.ParseContainerOutputFunc
	formatGroupedContainers container_utils.FormatGroupedContainersFunc
}

type Option func(*realContainerCLI)

func NewContainerCLI(opts ...Option) *realContainerCLI {
	c := &realContainerCLI{
		parseContainerOutput:    container_utils.ParseContainerOutput,
		formatGroupedContainers: container_utils.FormatGroupedContainers,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithExecutor(e exec.Executor) Option {
	return func(d *realContainerCLI) {
		d.executor = e
	}
}

func WithConfig(c config.Config) Option {
	return func(d *realContainerCLI) {
		d.config = c
	}
}

func WithParseContainerOutput(f container_utils.ParseContainerOutputFunc) Option {
	return func(d *realContainerCLI) {
		d.parseContainerOutput = f
	}
}

func WithFormatGroupedContainers(f container_utils.FormatGroupedContainersFunc) Option {
	return func(d *realContainerCLI) {
		d.formatGroupedContainers = f
	}
}
