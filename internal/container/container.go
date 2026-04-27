package container

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
)

type ContainerCLI interface {
	ListContainersOfActiveDevcontainers() error
	CleanResources() error
	DownContainer(path string) error
	GetAllRelatedContainers(path string) ([]string, error)
}

type ResolvePathsFunc func(path string, pth pather.Pather) []string
type FindMainContainersForPathFunc func(tool string, p string, executor exec.Executor) ([]string, error)
type ExtractProjectFromContainerFunc func(tool string, id string, executor exec.Executor) (string, error)
type FindComposeContainersForProjectFunc func(tool string, project string, executor exec.Executor) ([]string, error)
type DeduplicateAndFilterContainerIDsFunc func(idMap map[string]bool) []string
