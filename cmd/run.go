package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [caminho]",
	Short: "Sobe o container e abre o VS Code",
	Long:  "Executa a rotina completa de inicialização: provisiona o container (equivalente ao 'up') e imediatamente anexa o VS Code ao ambiente remoto. Resolve dinamicamente o workspaceFolder e contorna falhas de URI em ambientes como WSL.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		absPath, _ := container.GetAbsPath(path)

		if err := container.RunUpSync(absPath); err != nil {
			return err
		}

		uri := container.GetContainerURI(absPath)
		return container.ExecDetached("code", "--folder-uri", uri)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
