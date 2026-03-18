package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
	"github.com/spf13/cobra"
)

var execPath string

var execCmd = &cobra.Command{
	Use:                "exec [comando] [args...]",
	Short:              "Executa um comando específico dentro do container",
	Long:               "Repassa instruções e argumentos para execução direta no contexto isolado do container ativo. O repasse de flags e parâmetros ocorre de forma transparente ao processo interno.",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		absPath, err := paths.GetAbsPath(execPath)
		if err != nil {
			return err
		}

		executor := exec.NewExecutor()
		devCli := devcontainer.NewDevContainerCLI(executor)

		return devCli.Exec(absPath, args)
	},
}

func init() {
	execCmd.Flags().StringVarP(&execPath, "path", "p", "", "Caminho do projeto (padrão '.')")
	execCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(execCmd)
}
