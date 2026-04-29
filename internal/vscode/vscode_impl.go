package vscode

import (
	"encoding/hex"
	"fmt"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
)

func (r *realVSCode) GetContainerWorkspaceURI(absPath string) (string, error) {

	logger.Info("Iniciando tratamento de caminho do Workspace")

	hostPath, err := r.pather.GetRealPath(absPath)
	if err != nil {
		logger.Error("Houve um erro ao resolver o caminho real do host")
		return "", err
	}

	hexPath := hex.EncodeToString([]byte(hostPath))
	containerPath, err := r.devcontainerCLI.GetWorkspaceFolder(absPath)

	if err != nil {
		logger.Error("Houve um erro ao obter o caminho do workspace dentro do container")
		return "", err
	}

	finalURI := fmt.Sprintf("vscode-remote://dev-container+%s%s", hexPath, containerPath)

	logger.Success("URI do workspace criada com sucesso")
	logger.Verbose(finalURI)

	return finalURI, nil
}

func (r *realVSCode) OpenWorkspaceByURI(workspaceURI string) error {
	return r.executor.RunDetached("code", "--folder-uri", workspaceURI)
}
