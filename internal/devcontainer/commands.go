package devcontainer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
)

var (
	ErrDevContainerNotFound = errors.New("devcontainer CLI não encontrado. Instale com: npm install -g @devcontainers/cli")
	ErrDockerNotFound       = errors.New("docker não encontrado. Verifique se o Docker está instalado e em execução")
	ErrDockerPermission     = errors.New("permissão negada ao acessar Docker")
)

type devContainerError struct {
	Outcome    string `json:"outcome"`
	Message    string `json:"message"`
	Desc       string `json:"description"`
	exitStatus int
}

func (e *devContainerError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("devcontainer falhou com exit status %d", e.exitStatus)
}

type devContainerOutput struct {
	Outcome string `json:"outcome"`
	Message string `json:"message"`
}

func (d *DevContainerCLIImpl) Up(workspaceFolder string) error {
	logs.Info("Subindo container de desenvolvimento")
	logs.Verbose("executando: devcontainer up --workspace-folder %s", workspaceFolder)

	out, err := d.executor.CombinedOutput("devcontainer", "up", "--workspace-folder", workspaceFolder)
	if err != nil {
		return wrapDevContainerError(err, out)
	}

	return nil
}

func (d *DevContainerCLIImpl) Exec(workspaceFolder string, command []string) error {
	args := append([]string{"exec", "--workspace-folder", workspaceFolder}, command...)
	logs.Verbose("executando: devcontainer %s", argsToString(args))
	err := d.executor.RunInteractive("devcontainer", args...)
	return wrapDevContainerError(err, "")
}

func wrapDevContainerError(err error, output string) error {
	if err == nil {
		return nil
	}

	var execErr *exec.Error
	if errors.As(err, &execErr) {
		if strings.Contains(execErr.Error(), "executable file not found") {
			if execErr.Name == "devcontainer" {
				return ErrDevContainerNotFound
			}
			if execErr.Name == "docker" {
				return ErrDockerNotFound
			}
		}
	}

	if output != "" {
		jsonErr := extractJSONError(output)
		if jsonErr != nil && jsonErr.Outcome == "error" {
			if isDockerPermissionError(jsonErr.Message) {
				return fmt.Errorf("%w\n\nO devcontainer CLI não consegue acessar o Docker.\nSolução: adicione seu usuário ao grupo docker com:\n  sudo usermod -aG docker $USER\n\nApós executar, faça logout/login ou reinicie o terminal", ErrDockerPermission)
			}
			return fmt.Errorf("%s", jsonErr.Message)
		}
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return &devContainerError{exitStatus: exitErr.ExitCode(), Message: "falha ao executar devcontainer. Verifique se o Docker está em execução"}
	}

	return err
}

func extractJSONError(output string) *devContainerOutput {
	lines := strings.Split(output, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		var result devContainerOutput
		if err := json.Unmarshal([]byte(line), &result); err == nil {
			return &result
		}
	}
	return nil
}

func isDockerPermissionError(msg string) bool {
	lowerMsg := strings.ToLower(msg)
	fmt.Println(lowerMsg)
	return strings.Contains(lowerMsg, "permission denied") ||
		strings.Contains(lowerMsg, "cannot connect to the docker daemon") ||
		strings.Contains(lowerMsg, "got permission denied") ||
		strings.Contains(lowerMsg, "is the docker daemon running")
}

func argsToString(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		if containsSpace(arg) {
			result += fmt.Sprintf("'%s'", arg)
		} else {
			result += arg
		}
	}
	return result
}

func containsSpace(s string) bool {
	for _, r := range s {
		if r == ' ' {
			return true
		}
	}
	return false
}
