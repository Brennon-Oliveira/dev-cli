package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/spf13/cobra"
)

var execPath string

type execImplParams struct {
	args      []string
	pather    pather.Pather
	container container.ContainerCLI
}

func execImpl(p *execImplParams) error {
	absPath, _ := p.pather.GetAbsPath(execPath)

	return p.container.RunInteractive(absPath, p.args[0])
}

var execCmd = &cobra.Command{
	Use:                "exec \"[comando]\"",
	Short:              "Executa um comando específico dentro do container",
	Long:               "Repassa instruções e argumentos para execução direta no contexto isolado do container ativo. O repasse de flags e parâmetros ocorre de forma transparente ao processo interno.",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		pather := pather.NewPather(
			pather.WithExecutor(executor),
		)
		config := config.NewConfig()

		container := container.NewContainerCLI(
			container.WithExecutor(executor),
			container.WithPather(pather),
			container.WithConfig(config),
		)

		return execImpl(&execImplParams{
			args:      args,
			pather:    pather,
			container: container,
		})
	},
}

func init() {
	execCmd.Flags().StringVarP(&execPath, "path", "p", "", "Caminho do projeto (padrão '.')")
	execCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(execCmd)
}
