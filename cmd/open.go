package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
	"github.com/Brennon-Oliveira/dev-cli/internal/vscode"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [caminho]",
	Short: "Abre o VS Code no container",
	Long:  "Abre o VS Code conectado a um dev container já em execução. Utiliza resolução dinâmica de URIs para forçar a montagem exata da raiz do projeto, independente da profundidade do diretório definido nas configurações.",
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
		devCli := devcontainer.NewDevContainerCLI(executor)

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
	rootCmd.AddCommand(openCmd)
}
