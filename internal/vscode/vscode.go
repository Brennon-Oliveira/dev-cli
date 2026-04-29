package vscode

type VSCode interface {
	GetContainerWorkspaceURI(absPath string) (string, error)
	OpenWorkspaceByURI(workspaceURI string) error
}
