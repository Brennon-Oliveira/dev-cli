package container

import "github.com/Brennon-Oliveira/dev-cli/internal/exec"

type ContainerClient interface {
	ListContainers() error
	GetContainerID(path string) (string, error)
	GetAllRelatedContainers(path string) ([]string, error)
	StopContainers(ids []string) error
	RemoveContainers(ids []string) error
	ShowLogs(path string, follow bool) error
	ListPorts(path string) error
	CleanResources() error
}

type DockerClient struct {
	tool     string
	executor exec.Executor
}

func NewDockerClient(tool string, executor exec.Executor) *DockerClient {
	return &DockerClient{
		tool:     tool,
		executor: executor,
	}
}

type MockContainerClient struct {
	ListContainersErr       error
	GetContainerIDResult    string
	GetContainerIDErr       error
	GetAllRelatedResult     []string
	GetAllRelatedErr        error
	StopContainersErr       error
	RemoveContainersErr     error
	ShowLogsErr             error
	ListPortsErr            error
	CleanResourcesErr       error
	Calls                   []string
}

func NewMockContainerClient() *MockContainerClient {
	return &MockContainerClient{}
}

func (m *MockContainerClient) ListContainers() error {
	m.Calls = append(m.Calls, "ListContainers")
	return m.ListContainersErr
}

func (m *MockContainerClient) GetContainerID(path string) (string, error) {
	m.Calls = append(m.Calls, "GetContainerID "+path)
	return m.GetContainerIDResult, m.GetContainerIDErr
}

func (m *MockContainerClient) GetAllRelatedContainers(path string) ([]string, error) {
	m.Calls = append(m.Calls, "GetAllRelatedContainers "+path)
	return m.GetAllRelatedResult, m.GetAllRelatedErr
}

func (m *MockContainerClient) StopContainers(ids []string) error {
	m.Calls = append(m.Calls, "StopContainers")
	return m.StopContainersErr
}

func (m *MockContainerClient) RemoveContainers(ids []string) error {
	m.Calls = append(m.Calls, "RemoveContainers")
	return m.RemoveContainersErr
}

func (m *MockContainerClient) ShowLogs(path string, follow bool) error {
	m.Calls = append(m.Calls, "ShowLogs")
	return m.ShowLogsErr
}

func (m *MockContainerClient) ListPorts(path string) error {
	m.Calls = append(m.Calls, "ListPorts")
	return m.ListPortsErr
}

func (m *MockContainerClient) CleanResources() error {
	m.Calls = append(m.Calls, "CleanResources")
	return m.CleanResourcesErr
}
