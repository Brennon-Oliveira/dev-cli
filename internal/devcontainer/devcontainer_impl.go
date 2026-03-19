package devcontainer

import (
	"encoding/json"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
)

func (d *realDevContainerCLI) Up(workspace string) error {
	logger.Info("Subindo dev containers")
	err := d.executor.Run("devcontainer", "up", "--workspace-folder", workspace)
	if err != nil {
		logger.Error("Houve um erro ao subir os devcontainers")
		return err
	}
	logger.Success("Containers subiram com sucesso")
	return nil
}

func (d *realDevContainerCLI) ReadConfiguration(absPath string) (*DevContainerConfiguration, error) {
	devcontainerJsonRaw, err := d.executor.Output("devcontainer", "read-configuration", "--workspace-folder", absPath)
	if err != nil {
		return nil, err
	}

	var config DevContainerConfiguration

	if err := json.Unmarshal(devcontainerJsonRaw, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func formatWorkspaceFolderSuffix(containerPath string) string {
	if strings.HasSuffix(containerPath, "workspaces/") {
		return containerPath + "/"
	}
	if !strings.HasSuffix(containerPath, "workspaces") {
		return containerPath + "//"
	}
	return containerPath
}

func (d *realDevContainerCLI) GetWorkspaceFolder(absPath string) (string, error) {
	config, err := d.ReadConfiguration(absPath)

	if err != nil {
		logger.Warn("Erro ao ler configuração do devcontainer para buscar 'workspaceFolder', tentando caminho padrão")
		return formatWorkspaceFolderSuffix("/workspaces"), err
	}

	if config.Workspace.WorkspaceFolder == "" {
		logger.Warn("Configuração de 'workspaceFolder' não foi encontrada, tentando caminho padrão")
		return formatWorkspaceFolderSuffix("/workspaces"), nil
	}

	workspaceFolder := formatWorkspaceFolderSuffix(config.Workspace.WorkspaceFolder)

	return workspaceFolder, nil
}
