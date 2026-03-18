package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:   "ports [caminho]",
	Short: "Lista as portas mapeadas do container",
	Long:  "Inspeciona as interfaces de rede do Motor de containers e exibe o mapeamento ativo de portas e protocolos expostos/bindados entre a máquina host e a rede isolada do container do workspace atual.",
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

		return client.ListPorts(absPath)
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
