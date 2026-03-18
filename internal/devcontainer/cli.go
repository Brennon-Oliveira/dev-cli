package devcontainer

import "github.com/Brennon-Oliveira/dev-cli/internal/exec"

type DevContainerCLI interface {
	Up(workspaceFolder string) error
	Exec(workspaceFolder string, command []string) error
	ReadConfiguration(workspaceFolder string) (*WorkspaceConfig, error)
}

type WorkspaceConfig struct {
	WorkspaceFolder string
}

type DevContainerCLIImpl struct {
	executor exec.Executor
}

func NewDevContainerCLI(executor exec.Executor) *DevContainerCLIImpl {
	return &DevContainerCLIImpl{executor: executor}
}

type MockDevContainerCLI struct {
	UpErr            error
	ExecErr          error
	ReadConfigResult *WorkspaceConfig
	ReadConfigErr    error
	Calls            []string
}

func NewMockDevContainerCLI() *MockDevContainerCLI {
	return &MockDevContainerCLI{
		ReadConfigResult: &WorkspaceConfig{WorkspaceFolder: "/workspaces"},
	}
}

func (m *MockDevContainerCLI) Up(workspaceFolder string) error {
	m.Calls = append(m.Calls, "up "+workspaceFolder)
	return m.UpErr
}

func (m *MockDevContainerCLI) Exec(workspaceFolder string, command []string) error {
	m.Calls = append(m.Calls, "exec "+workspaceFolder)
	return m.ExecErr
}

func (m *MockDevContainerCLI) ReadConfiguration(workspaceFolder string) (*WorkspaceConfig, error) {
	m.Calls = append(m.Calls, "read-configuration "+workspaceFolder)
	return m.ReadConfigResult, m.ReadConfigErr
}
