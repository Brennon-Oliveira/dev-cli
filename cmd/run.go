package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
	"github.com/Brennon-Oliveira/dev-cli/internal/vscode"
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
		absPath, err := paths.GetAbsPath(path)
		if err != nil {
			return err
		}

		executor := exec.NewExecutor()
		cfg := config.Load()
		devCli := devcontainer.NewDevContainerCLI(executor)

		if err := devCli.Up(absPath); err != nil {
			return err
		}

		workspaceConfig, err := devCli.ReadConfiguration(absPath)
		if err != nil {
			return err
		}

		uri := vscode.GetContainerURI(absPath, workspaceConfig.WorkspaceFolder)
		logs.Info("Abrindo VS Code...")
		logs.Verbose("executando: code --folder-uri %s", uri)

		return executor.RunDetached("code", "--folder-uri", uri)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
