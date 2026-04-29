package container

type ContainerCLI interface {
	ListContainersOfActiveDevcontainers() error
	CleanResources() error
	DownContainer(path string) error
	GetAllRelatedContainers(path string) ([]string, error)
	KillContainer(path string) error
	ShowLogs(path string, follow bool) error
}
