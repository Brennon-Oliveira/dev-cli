package container

type ContainerCLI interface {
	ListContainersOfActiveDevcontainers() error
	CleanResources() error
	DownContainer(path string) error
	GetAllRelatedContainers(path string) ([]string, error)
	RunInteractive(path string, command string) error
}
