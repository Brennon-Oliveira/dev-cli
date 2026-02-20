package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove containers e redes parados no Docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		return container.CleanResources()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
