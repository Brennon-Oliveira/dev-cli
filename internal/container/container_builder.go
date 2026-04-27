package container

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container/container_utils"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
)

type realContainerCLI struct {
	executor                         exec.Executor
	config                           config.Config
	pather                           pather.Pather
	parseContainerOutput             container_utils.ParseContainerOutputFunc
	formatGroupedContainers          container_utils.FormatGroupedContainersFunc
	tryPaths                         container_utils.TryPathsFunc
	findMainContainersForPath        container_utils.FindMainContainersForPathFunc
	extractProjectFromContainer      container_utils.ExtractProjectFromContainerFunc
	findComposeContainersForProject  container_utils.FindComposeContainersForProjectFunc
	deduplicateAndFilterContainerIDs container_utils.DeduplicateAndFilterContainerIDsFunc
}

type Option func(*realContainerCLI)

func NewContainerCLI(opts ...Option) *realContainerCLI {
	c := &realContainerCLI{
		parseContainerOutput:             container_utils.ParseContainerOutput,
		formatGroupedContainers:          container_utils.FormatGroupedContainers,
		tryPaths:                         container_utils.TryPaths,
		findMainContainersForPath:        container_utils.FindMainContainersForPath,
		extractProjectFromContainer:      container_utils.ExtractProjectFromContainer,
		findComposeContainersForProject:  container_utils.FindComposeContainersForProject,
		deduplicateAndFilterContainerIDs: container_utils.DeduplicateAndFilterContainerIDs,
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

func WithTryPaths(f container_utils.TryPathsFunc) Option {
	return func(c *realContainerCLI) {
		c.tryPaths = f
	}
}

func WithFindMainContainersForPath(f container_utils.FindMainContainersForPathFunc) Option {
	return func(c *realContainerCLI) {
		c.findMainContainersForPath = f
	}
}

func WithExtractProjectFromContainer(f container_utils.ExtractProjectFromContainerFunc) Option {
	return func(c *realContainerCLI) {
		c.extractProjectFromContainer = f
	}
}

func WithFindComposeContainersForProject(f container_utils.FindComposeContainersForProjectFunc) Option {
	return func(c *realContainerCLI) {
		c.findComposeContainersForProject = f
	}
}

func WithDeduplicateAndFilterContainerIDs(f container_utils.DeduplicateAndFilterContainerIDsFunc) Option {
	return func(c *realContainerCLI) {
		c.deduplicateAndFilterContainerIDs = f
	}
}
