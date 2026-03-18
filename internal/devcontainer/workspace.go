package devcontainer

import (
	"encoding/json"

	"github.com/Brennon-Oliveira/dev-cli/internal/constants"
	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
)

func (d *DevContainerCLIImpl) ReadConfiguration(workspaceFolder string) (*WorkspaceConfig, error) {
	logs.Verbose("executando: devcontainer read-configuration --workspace-folder %s", workspaceFolder)

	out, err := d.executor.Output("devcontainer", "read-configuration", "--workspace-folder", workspaceFolder)
	if err != nil {
		logs.Verbose("falha ao ler configuração, usando workspace padrão: %s", constants.DefaultWorkspaceFolder)
		return &WorkspaceConfig{WorkspaceFolder: constants.DefaultWorkspaceFolder}, nil
	}

	var config struct {
		Workspace struct {
			WorkspaceFolder string `json:"workspaceFolder"`
		} `json:"workspace"`
	}

	if err := json.Unmarshal([]byte(out), &config); err != nil {
		logs.Verbose("falha ao parsear configuração, usando workspace padrão: %s", constants.DefaultWorkspaceFolder)
		return &WorkspaceConfig{WorkspaceFolder: constants.DefaultWorkspaceFolder}, nil
	}

	if config.Workspace.WorkspaceFolder == "" {
		logs.Verbose("workspaceFolder vazio, usando padrão: %s", constants.DefaultWorkspaceFolder)
		return &WorkspaceConfig{WorkspaceFolder: constants.DefaultWorkspaceFolder}, nil
	}

	logs.Verbose("workspaceFolder detectado: %s", config.Workspace.WorkspaceFolder)
	return &WorkspaceConfig{WorkspaceFolder: config.Workspace.WorkspaceFolder}, nil
}
