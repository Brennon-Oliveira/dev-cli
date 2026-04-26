package container

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container/container_utils"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
)

type realContainerCLI struct {
	executor                exec.Executor
	config                  config.Config
	pather                  pather.Pather
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
	return func(c *realContainerCLI) {
		c.executor = e
	}
}

func WithConfig(co config.Config) Option {
	return func(c *realContainerCLI) {
		c.config = co
	}
}

func WithPather(p pather.Pather) Option {
	return func(c *realContainerCLI) {
		c.pather = p
	}
}

func WithParseContainerOutput(f container_utils.ParseContainerOutputFunc) Option {
	return func(c *realContainerCLI) {
		c.parseContainerOutput = f
	}
}

func WithFormatGroupedContainers(f container_utils.FormatGroupedContainersFunc) Option {
	return func(c *realContainerCLI) {
		c.formatGroupedContainers = f
	}
}
