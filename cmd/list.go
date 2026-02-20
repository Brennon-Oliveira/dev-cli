package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista dev containers ativos",
	RunE: func(cmd *cobra.Command, args []string) error {
		return container.ListContainers()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
