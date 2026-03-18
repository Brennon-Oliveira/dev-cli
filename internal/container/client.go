package container

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
)

var ErrPermissionDenied = errors.New("permissão negada ao acessar Docker/Podman")

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
	useSudo  bool
	executor exec.Executor
}

func NewDockerClient(tool string, useSudo bool, executor exec.Executor) *DockerClient {
	return &DockerClient{
		tool:     tool,
		useSudo:  useSudo,
		executor: executor,
	}
}

func (d *DockerClient) buildArgs(args ...string) []string {
	if d.useSudo {
		return append([]string{"sudo", d.tool}, args...)
	}
	return append([]string{d.tool}, args...)
}

func wrapPermissionError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()
	if strings.Contains(errStr, "permission denied") ||
		strings.Contains(errStr, "Permission denied") ||
		strings.Contains(errStr, "Cannot connect to the Docker daemon") {
		return fmt.Errorf("%w\n\nDica: execute 'dev config core.use-sudo true --global' para usar sudo", ErrPermissionDenied)
	}

	return err
}

type MockContainerClient struct {
	ListContainersErr    error
	GetContainerIDResult string
	GetContainerIDErr    error
	GetAllRelatedResult  []string
	GetAllRelatedErr     error
	StopContainersErr    error
	RemoveContainersErr  error
	ShowLogsErr          error
	ListPortsErr         error
	CleanResourcesErr    error
	Calls                []string
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
