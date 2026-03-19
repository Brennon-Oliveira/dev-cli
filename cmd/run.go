package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/Brennon-Oliveira/dev-cli/internal/vscode"
	"github.com/spf13/cobra"
)

type runImplParams struct {
	args         []string
	pather       pather.Pather
	devcontainer devcontainer.DevContainerCLI
	vscode       vscode.VSCode
}

func runImpl(params *runImplParams) error {
	args := params.args
	pather := params.pather
	devcontainer := params.devcontainer
	vscode := params.vscode
	logger.Info("Iniciando projeto")
	path := pather.GetPathFromArgs(args)
	absPath, _ := pather.GetAbsPath(path)

	logger.Verbose("Rodando projeto na pasta %s", absPath)

	if err := devcontainer.Up(absPath); err != nil {
		return err
	}

	workspaceURI, err := vscode.GetContainerWorkspaceURI(absPath)

	if err != nil {
		return err
	}

	return vscode.OpenWorkspaceByURI(workspaceURI)
}

var runCmd = &cobra.Command{
	Use:   "run [caminho]",
	Short: "Sobe o container e abre o VS Code",
	Long:  "Executa a rotina completa de inicialização: provisiona o container (equivalente ao 'up') e imediatamente anexa o VS Code ao ambiente remoto. Resolve dinamicamente o workspaceFolder e contorna falhas de URI em ambientes como WSL.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		pather := pather.NewPather(
			pather.WithExecutor(executor),
		)
		devcontainerCLI := devcontainer.NewDevContainerCLI(
			devcontainer.WithExecutor(executor),
		)
		vscode := vscode.NewVSCode(
			vscode.WithPather(pather),
			vscode.WithExecutor(executor),
			vscode.WithDevcontainerCLI(devcontainerCLI),
		)

		return runImpl(&runImplParams{
			args:         args,
			pather:       pather,
			devcontainer: devcontainerCLI,
			vscode:       vscode,
		})
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
