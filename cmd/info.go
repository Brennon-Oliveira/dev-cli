package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/spf13/cobra"
)

type infoImplParams struct {
	args      []string
	container container.ContainerCLI
}

func infoImpl(params *infoImplParams) error {
	container := params.container

	return container.ListContainersOfActiveDevcontainers()
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Lista dev containers ativos",
	Long:  "Consulta o daemon do Motor de containers e retorna uma listagem contendo exclusivamente os processos mapeados como Dev Containers, filtrando ativamente através das labels de controle da extensão.",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		config := config.NewConfig()

		container := container.NewContainerCLI(
			container.WithExecutor(executor),
			container.WithConfig(config),
		)

		return infoImpl(&infoImplParams{
			args:      args,
			container: container,
		})
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
