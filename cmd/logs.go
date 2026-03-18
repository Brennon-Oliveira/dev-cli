package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
	"github.com/spf13/cobra"
)

var follow bool

var logsCmd = &cobra.Command{
	Use:   "logs [caminho]",
	Short: "Exibe os logs do container do workspace",
	Long:  "Exibe o stream de saída padrão (stdout/stderr) do container ativo. Utilizado para diagnóstico e debug de falhas de provisionamento, scripts de entrypoint ou da aplicação interna rodando em background.",
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
		client := container.NewDockerClient(cfg.Core.Tool, executor)

		return client.ShowLogs(absPath, follow)
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Acompanha os logs em tempo real")
	rootCmd.AddCommand(logsCmd)
}
