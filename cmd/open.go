package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/Brennon-Oliveira/dev-cli/internal/vscode"
	"github.com/spf13/cobra"
)

type openImplParams struct {
	args   []string
	pather pather.Pather
	vscode vscode.VSCode
}

func openImpl(p *openImplParams) error {
	args := p.args
	logger.Info("Iniciando projeto")
	path := p.pather.GetPathFromArgs(args)
	absPath, _ := p.pather.GetAbsPath(path)

	workspaceURI, err := p.vscode.GetContainerWorkspaceURI(absPath)

	if err != nil {
		return err
	}

	logger.Info("Abrindo editor")
	logger.Verbose("Abrindo VS Code com URI: %s", workspaceURI)

	return p.vscode.OpenWorkspaceByURI(workspaceURI)
}

var openCmd = &cobra.Command{
	Use:   "open [caminho]",
	Short: "Abre o VS Code no container",
	Long:  "Abre o VS Code conectado a um dev container já em execução. Utiliza resolução dinâmica de URIs para forçar a montagem exata da raiz do projeto, independente da profundidade do diretório definido nas configurações.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		pather := pather.NewPather(
			pather.WithExecutor(executor),
		)

		devcontainer := devcontainer.NewDevContainerCLI(
			devcontainer.WithExecutor(executor),
		)
		vscode := vscode.NewVSCode(
			vscode.WithExecutor(executor),
			vscode.WithPather(pather),
			vscode.WithDevcontainerCLI(devcontainer),
		)

		return openImpl(&openImplParams{
			args:   args,
			pather: pather,
			vscode: vscode,
		})
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
