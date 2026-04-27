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
	resolvePaths                     ResolvePathsFunc
	findMainContainersForPath        FindMainContainersForPathFunc
	extractProjectFromContainer      ExtractProjectFromContainerFunc
	findComposeContainersForProject  FindComposeContainersForProjectFunc
	deduplicateAndFilterContainerIDs DeduplicateAndFilterContainerIDsFunc
}

type Option func(*realContainerCLI)

func NewContainerCLI(opts ...Option) *realContainerCLI {
	c := &realContainerCLI{
		parseContainerOutput:             container_utils.ParseContainerOutput,
		formatGroupedContainers:          container_utils.FormatGroupedContainers,
		resolvePaths:                     container_utils.ResolvePaths,
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

func WithResolvePaths(f ResolvePathsFunc) Option {
	return func(c *realContainerCLI) {
		c.resolvePaths = f
	}
}

func WithFindMainContainersForPath(f FindMainContainersForPathFunc) Option {
	return func(c *realContainerCLI) {
		c.findMainContainersForPath = f
	}
}

func WithExtractProjectFromContainer(f ExtractProjectFromContainerFunc) Option {
	return func(c *realContainerCLI) {
		c.extractProjectFromContainer = f
	}
}

func WithFindComposeContainersForProject(f FindComposeContainersForProjectFunc) Option {
	return func(c *realContainerCLI) {
		c.findComposeContainersForProject = f
	}
}

func WithDeduplicateAndFilterContainerIDs(f DeduplicateAndFilterContainerIDsFunc) Option {
	return func(c *realContainerCLI) {
		c.deduplicateAndFilterContainerIDs = f
	}
}
