package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:   "ports [caminho]",
	Short: "Lista as portas mapeadas do container",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		absPath, _ := container.GetAbsPath(path)
		return container.ListPorts(absPath)
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
