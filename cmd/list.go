package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista dev containers ativos",
	Long:  "Consulta o daemon do Motor de containers e retorna uma listagem contendo exclusivamente os processos mapeados como Dev Containers, filtrando ativamente através das labels de controle da extensão.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return container.ListContainers()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
