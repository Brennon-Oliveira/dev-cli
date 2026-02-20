package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove containers e redes parados no Docker",
	Long:  "Executa uma varredura de manutenção no host, liberando recursos do Docker ao remover agressivamente containers parados, volumes anônimos não referenciados e redes órfãs geradas pelos ciclos de Dev Containers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return container.CleanResources()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
