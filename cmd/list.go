package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista dev containers ativos",
	Long:  "Consulta o daemon do Motor de containers e retorna uma listagem contendo exclusivamente os processos mapeados como Dev Containers, filtrando ativamente através das labels de controle da extensão.",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		cfg := config.Load()
		client := container.NewDockerClient(cfg.Core.Tool, cfg.Core.UseSudo, executor)

		return client.ListContainers()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
