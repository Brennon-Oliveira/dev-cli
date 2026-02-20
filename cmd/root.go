package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "CLI para gerenciar Dev Containers",
}

func Execute() error {
	return rootCmd.Execute()
}
