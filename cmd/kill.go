package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill [caminho]",
	Short: "Encerra o container do workspace atual",
	Long:  "Força o encerramento e destrói (docker rm -f) o container alvo e todos os serviços acoplados via Docker Compose, limpando de forma definitiva o estado de execução daquele workspace no Docker do host.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		absPath, _ := container.GetAbsPath(path)
		return container.KillContainer(absPath)
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
