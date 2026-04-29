package cmd

import (
	"os"

	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/spf13/cobra"
)

var follow bool

type logsImplParams struct {
	args      []string
	pather    pather.Pather
	container container.ContainerCLI
}

func logsImpl(p *logsImplParams) error {
	path := p.pather.GetPathFromArgs(p.args)
	absPath, _ := p.pather.GetAbsPath(path)

	logger.Info("Buscando logs do container")
	logger.Verbose("Caminho absoluto encontrado: %s", absPath)

	return p.container.ShowLogs(absPath, follow)
}

var logsCmd = &cobra.Command{
	Use:   "logs [caminho]",
	Short: "Exibe os logs do container do workspace",
	Long:  "Exibe o stream de saída padrão (stdout/stderr) do container ativo. Utilizado para diagnóstico e debug de falhas de provisionamento, scripts de entrypoint ou da aplicação interna rodando em background.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executorIO := exec.NewExecutor(
			exec.WithStdin(os.Stdin),
			exec.WithStdout(os.Stdout),
		)
		executor := exec.NewExecutor()
		config := config.NewConfig()
		pather := pather.NewPather(
			pather.WithExecutor(executor),
		)

		container := container.NewContainerCLI(
			container.WithExecutor(executorIO),
			container.WithConfig(config),
			container.WithPather(pather),
		)

		return logsImpl(&logsImplParams{
			args:      args,
			pather:    pather,
			container: container,
		})
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Acompanha os logs em tempo real")
	rootCmd.AddCommand(logsCmd)
}
