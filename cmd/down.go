package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/spf13/cobra"
)

type downImplParams struct {
	args      []string
	pather    pather.Pather
	container container.ContainerCLI
}

func downImpl(p *downImplParams) error {
	path := p.pather.GetPathFromArgs(p.args)
	absPath, _ := p.pather.GetAbsPath(path)

	logger.Info("Iniciando queda dos containers")
	logger.Verbose("Caminho absoluto encontrado: %s", absPath)

	return p.container.DownContainer(absPath)
}

var downCmd = &cobra.Command{
	Use:   "down [caminho]",
	Short: "Para graciosamente o container do workspace atual",
	Long:  "Executa a parada graciosa do container principal e de todos os serviços secundários (bancos de dados, caches, etc.) vinculados à mesma stack do composer do Motor de containers, mantendo os containers intactos para reinício rápido.",
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

		return downImpl(&downImplParams{
			args:      args,
			pather:    pather,
			container: container,
		})
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
