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
	if strings.HasSuffix(containerPath, "workspaces") {
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

func (c *realDevContainerCLI) RunInteractive(path string, command string) error {
	tool := "devcontainer"

	commandToExecute := strings.Split(command, " ")

	baseCommand := []string{"exec", "--workspace-folder", path}

	finalCommand := append(baseCommand, commandToExecute...)

	logger.Info("Executando comando interativo: %s %s \"%s\"", tool, strings.Join(baseCommand, " "), command)

	out, err := c.executor.Output(tool, finalCommand...)

	if err != nil {
		if out != nil && string(out) != "" {
			logger.Error("O comando foi executado, e resultou no erro:\n```\n%s```", string(out))
			return nil
		}

		logger.Error("Houve um erro ao executar o comando interativo.")
		return err
	}

	logger.Info("O comando foi executado, e resultou na saída:\n```\n%s```", string(out))

	return nil
}

func (c *realDevContainerCLI) OpenShell(path string) error {
	tool := "devcontainer"

	shells := []string{"/bin/zsh", "/bin/bash", "/bin/sh"}
	var preferredShell string

	for _, shell := range shells {
		checkCmd := []string{"exec", "--workspace-folder", path, "test", "-x", shell}
		if c.executor.Run(tool, checkCmd...) == nil {
			preferredShell = shell
			break
		}
	}

	if preferredShell == "" {
		preferredShell = "/bin/sh"
	}

	shellArgs := []string{"exec", "--workspace-folder", path, preferredShell}

	logger.Info("Abrindo shell interativo: %s", preferredShell)
	logger.Verbose("Abrindo shell interativo com o comando: %s %s", tool, strings.Join(shellArgs, " "))

	err := c.executor.RunInteractive(tool, shellArgs...)

	if err != nil {
		logger.Error("Houve um erro ao abrir o shell interativo.")
		return err
	}

	logger.Info("Sessão encerrada")

	return nil
}
