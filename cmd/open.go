package cmd

import (
	"fmt"

	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [caminho]",
	Short: "Abre o VS Code no container",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		absPath, _ := container.GetAbsPath(path)
		uri := container.GetContainerURI(absPath)
		fmt.Printf("Abrindo VS Code...\n")
		return container.ExecDetached("code", "--folder-uri", uri)
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
