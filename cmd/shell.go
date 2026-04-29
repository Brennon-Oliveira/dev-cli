package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/spf13/cobra"
)

type shellImplParams struct {
	args         []string
	pather       pather.Pather
	devcontainer devcontainer.DevContainerCLI
}

func shellImpl(p *shellImplParams) error {
	path := p.pather.GetPathFromArgs(p.args)
	absPath, _ := p.pather.GetAbsPath(path)

	logger.Info("Iniciando shell interativo")
	logger.Verbose("Caminho absoluto encontrado: %s", absPath)

	return p.devcontainer.OpenShell(absPath)
}

var shellCmd = &cobra.Command{
	Use:   "shell [caminho]",
	Short: "Abre um shell interativo dentro do container",
	Long:  "Aloca um TTY e injeta uma sessão de terminal interativa no container ativo, detectando e priorizando automaticamente o uso de zsh, bash ou sh, conforme a disponibilidade no sistema de arquivos remoto.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		pather := pather.NewPather(
			pather.WithExecutor(executor),
		)

		devcontainer := devcontainer.NewDevContainerCLI(
			devcontainer.WithExecutor(executor),
		)

		return shellImpl(&shellImplParams{
			args:         args,
			pather:       pather,
			devcontainer: devcontainer,
		})
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
