package devcontainer

type DevContainerCLI interface {
	Up(workspace string) error
	GetWorkspaceFolder(absPath string) (string, error)
	ReadConfiguration(absPath string) (*DevContainerConfiguration, error)
	RunInteractive(path string, command string) error
	OpenShell(path string) error
}

type DevContainerConfiguration_Workspace struct {
	WorkspaceFolder string `json:"workspaceFolder"`
}

type DevContainerConfiguration struct {
	Workspace DevContainerConfiguration_Workspace `json:"workspace"`
}
