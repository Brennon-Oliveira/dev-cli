package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill [caminho]",
	Short: "Encerra o container do workspace atual",
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
