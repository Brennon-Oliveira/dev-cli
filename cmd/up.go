package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
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

		absPath, err := paths.GetAbsPath(path)
		if err != nil {
			return err
		}

		logs.Info("Subindo container em: %s", absPath)

		executor := exec.NewExecutor()
		cfg := config.Load()
		devCli := devcontainer.NewDevContainerCLI(executor)

		logs.Debug("usando tool: %s", cfg.Core.Tool)
		return devCli.Up(absPath)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
