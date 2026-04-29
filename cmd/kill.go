package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/spf13/cobra"
)

type killImplParams struct {
	args      []string
	pather    pather.Pather
	container container.ContainerCLI
}

func killImpl(p *killImplParams) error {
	path := p.pather.GetPathFromArgs(p.args)
	absPath, _ := p.pather.GetAbsPath(path)

	logger.Info("Iniciando exclusão dos containers")
	logger.Verbose("Caminho absoluto encontrado: %s", absPath)

	return p.container.KillContainer(absPath)
}

var killCmd = &cobra.Command{
	Use:   "kill [caminho]",
	Short: "Encerra o container do workspace atual",
	Long:  "Força o encerramento e destrói o container alvo e todos os serviços acoplados via composer do Motor de containers, limpando de forma definitiva o estado de execução daquele workspace no Motor de containers do host.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		config := config.NewConfig()
		pather := pather.NewPather(
			pather.WithExecutor(executor),
		)

		container := container.NewContainerCLI(
			container.WithExecutor(executor),
			container.WithConfig(config),
			container.WithPather(pather),
		)

		return killImpl(&killImplParams{
			args:      args,
			pather:    pather,
			container: container,
		})
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
