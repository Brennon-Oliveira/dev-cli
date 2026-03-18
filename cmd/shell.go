package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell [caminho]",
	Short: "Abre um shell interativo dentro do container",
	Long:  "Aloca um TTY e injeta uma sessão de terminal interativa no container ativo, detectando e priorizando automaticamente o uso de zsh, bash ou sh, conforme a disponibilidade no sistema de arquivos remoto.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		absPath, err := paths.GetAbsPath(path)
		if err != nil {
			return err
		}

		executor := exec.NewExecutor()
		devCli := devcontainer.NewDevContainerCLI(executor)

		return devCli.Exec(absPath, []string{"/bin/sh", "-c", "if command -v zsh >/dev/null 2>&1; then zsh; elif command -v bash >/dev/null 2>&1; then bash; else sh; fi"})
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
