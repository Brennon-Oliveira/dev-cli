package cmd

import (
	"os"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "dev",
	Short:   "CLI para gerenciar Dev Containers",
	Version: Version,
	Long: `O Dev CLI é uma interface de linha de comando para orquestrar o ciclo de vida de Dev Containers e a integração nativa com o VS Code.

Ele permite provisionar, acessar e destruir ambientes de desenvolvimento isolados diretamente pelo terminal, sem depender da interface gráfica do editor para a gestão do Motor de containers. Possui suporte avançado para roteamento dinâmico de caminhos no WSL, leitura direta do devcontainer.json para montagem de workspaces e controle de ecossistemas acoplados (composer do Motor de containers).`,
	Example: `  # Provisiona o container no diretório atual e abre o VS Code
  dev run .

  # Inicia a infraestrutura em background sem abrir o editor
  dev up /caminho/do/projeto

  # Abre o VS Code conectado a um container que já está rodando
  dev open .

  # Injeta um terminal interativo no container atual
  dev shell

  # Derruba e exclui o container e todos os serviços acoplados
  dev kill .`,
	SilenceUsage: true,
}

func initLogger() {
	logger.InitLogger(nil)
	var verboseFlag bool
	rootCmd.Flags().BoolVar(&verboseFlag, "verbose", false, "Roda o comando em modo verbose (mostra todos os logs executados por baixo dos panos)")
	if verboseFlag {
		logger.SetVerbose(verboseFlag)
	}
	logger.SetOutput(os.Stdout)
}

func Execute() {
	initLogger()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
