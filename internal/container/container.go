package container

type ContainerCLI interface {
	ListContainersOfActiveDevcontainers() error
	CleanResources() error
	DownContainer(path string) error
}
