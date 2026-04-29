package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/spf13/cobra"
)

type cleanImplParams struct {
	container container.ContainerCLI
}

func cleanImpl(p *cleanImplParams) error {
	return p.container.CleanResources()
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove containers e redes parados no Docker",
	Long:  "Executa uma varredura de manutenção no host, liberando recursos do Motor de containers ao remover agressivamente containers parados, volumes anônimos não referenciados e redes órfãs geradas pelos ciclos de Dev Containers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		config := config.NewConfig()

		container := container.NewContainerCLI(
			container.WithExecutor(executor),
			container.WithConfig(config),
		)

		return cleanImpl(&cleanImplParams{
			container: container,
		})
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
