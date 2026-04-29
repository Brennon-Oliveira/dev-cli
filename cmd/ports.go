package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/spf13/cobra"
)

type portsImplParams struct {
	args      []string
	pather    pather.Pather
	container container.ContainerCLI
}

func portsImpl(p *portsImplParams) error {
	path := p.pather.GetPathFromArgs(p.args)
	absPath, _ := p.pather.GetAbsPath(path)

	logger.Info("Bucando portas")
	logger.Verbose("Caminho absoluto encontrado: %s", absPath)
	p.container.ListPorts(absPath)
	return nil
}

var portsCmd = &cobra.Command{
	Use:   "ports [caminho]",
	Short: "Lista as portas mapeadas do container",
	Long:  "Inspeciona as interfaces de rede do Motor de containers e exibe o mapeamento ativo de portas e protocolos expostos/bindados entre a máquina host e a rede isolada do container do workspace atual.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()
		config := config.NewConfig()
		pather := pather.NewPather(
			pather.WithExecutor(executor),
		)

		container := container.NewContainerCLI(
			container.WithExecutor(executor),
			container.WithConfig(config),
			container.WithPather(pather),
		)

		return portsImpl(&portsImplParams{
			args:      args,
			pather:    pather,
			container: container,
		})
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
