package devcontainer

import (
	"fmt"

	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
)

func (d *DevContainerCLIImpl) Up(workspaceFolder string) error {
	logs.Info("Subindo container de desenvolvimento")
	logs.Verbose("executando: devcontainer up --workspace-folder %s", workspaceFolder)
	return d.executor.Run("devcontainer", "up", "--workspace-folder", workspaceFolder)
}

func (d *DevContainerCLIImpl) Exec(workspaceFolder string, command []string) error {
	args := append([]string{"exec", "--workspace-folder", workspaceFolder}, command...)
	logs.Verbose("executando: devcontainer %s", argsToString(args))
	return d.executor.RunInteractive("devcontainer", args...)
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
