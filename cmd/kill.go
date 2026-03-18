package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill [caminho]",
	Short: "Encerra o container do workspace atual",
	Long:  "Força o encerramento e destrói o container alvo e todos os serviços acoplados via composer do Motor de containers, limpando de forma definitiva o estado de execução daquele workspace no Motor de containers do host.",
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
		client := container.NewDockerClient(cfg.Core.Tool, cfg.Core.UseSudo, executor)

		ids, err := client.GetAllRelatedContainers(absPath)
		if err != nil {
			return err
		}

		return client.RemoveContainers(ids)
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
