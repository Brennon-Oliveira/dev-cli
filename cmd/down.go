package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down [caminho]",
	Short: "Para graciosamente o container do workspace atual (docker stop)",
	Long:  "Executa a parada graciosa (docker stop) do container principal e de todos os serviços secundários (bancos de dados, caches, etc.) vinculados à mesma stack do Docker Compose, mantendo os containers intactos para reinício rápido.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		absPath, _ := container.GetAbsPath(path)
		return container.DownContainer(absPath)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
