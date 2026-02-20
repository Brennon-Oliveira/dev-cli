package cmd

import (
	"fmt"

	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [caminho]",
	Short: "Apenas sobe o devcontainer",
	Long:  "Provisiona e inicia o dev container associado ao diretório atual em segundo plano (background). Executa o build da imagem e aplica as configurações do devcontainer.json sem instanciar a interface gráfica do VS Code.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		absPath, _ := container.GetAbsPath(path)
		fmt.Printf("Subindo container em: %s\n", absPath)
		return container.RunUpSync(absPath)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
