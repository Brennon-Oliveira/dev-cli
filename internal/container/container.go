package container

type ContainerCLI interface {
	ListContainersOfActiveDevcontainers() error
}
