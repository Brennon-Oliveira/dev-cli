package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/spf13/cobra"
)

type upImplParams struct {
	args         []string
	pather       pather.Pather
	devcontainer devcontainer.DevContainerCLI
}

func upImpl(p *upImplParams) error {
	logger.Info("Iniciando projeto")
	path := p.pather.GetPathFromArgs(p.args)
	absPath, _ := p.pather.GetAbsPath(path)

	logger.Verbose("Rodando projeto na pasta %s", absPath)

	return p.devcontainer.Up(absPath)
}

var upCmd = &cobra.Command{
	Use:   "up [caminho]",
	Short: "Apenas sobe o devcontainer",
	Long:  "Provisiona e inicia o dev container associado ao diretório atual em segundo plano (background). Executa o build da imagem e aplica as configurações do devcontainer.json sem instanciar a interface gráfica do VS Code.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		pather := pather.NewPather(
			pather.WithExecutor(executor),
		)
		devcontainerCLI := devcontainer.NewDevContainerCLI(
			devcontainer.WithExecutor(executor),
		)

		return upImpl(&upImplParams{
			args:         args,
			pather:       pather,
			devcontainer: devcontainerCLI,
		})
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
