package devcontainer

type DevContainerCLI interface {
	Up(workspace string) error
	GetWorkspaceFolder(absPath string) (string, error)
	ReadConfiguration(absPath string) (*DevContainerConfiguration, error)
}

type DevContainerConfiguration struct {
	Workspace struct {
		WorkspaceFolder string `json:"workspaceFolder"`
	} `json:"workspace"`
}
