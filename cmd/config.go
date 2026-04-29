package cmd

import (
	"fmt"
	"strings"

	conf "github.com/Brennon-Oliveira/dev-cli/internal-old/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/spf13/cobra"
)

var globalFlag bool
var interactiveFlag bool

type configImplParams struct {
	args   []string
	config config.Config
}

func configImpl(p *configImplParams) error {
	if !globalFlag {
		logger.Error("Atualmente apenas a flag --global é suportada")
		return nil
	}

	key := p.args[0]
	if !p.config.ValidateKey(key) {
		logger.Error("Chave `%s` de configuração desconhecida", key)
		return nil
	}

	if len(p.args) == 2 || interactiveFlag {
		value, err := p.config.TrySave(key, p.args[1])
		if err != nil {
			logger.Error("Erro ao salvar a configuração `%s` com o valor `%s`", key, value)
			return nil
		}

		logger.Info("Configuração '%s' atualizada para: %s", key, value)
		return nil
	}

	value := p.config.LoadByKey(key)

	logger.Info(value)

	return nil
}

type configHandler struct {
	ValidValues []string
	Label       string
	Get         func(cfg conf.GlobalConfig) string
	Set         func(cfg *conf.GlobalConfig, val string)
}

var handlers = map[string]configHandler{
	"core.tool": {
		ValidValues: []string{"docker", "podman"},
		Label:       "Selecione o motor de containers padrão",
		Get: func(cfg conf.GlobalConfig) string {
			return cfg.Core.Tool
		},
		Set: func(cfg *conf.GlobalConfig, val string) {
			cfg.Core.Tool = val
		},
	},
}

var configCmd = &cobra.Command{
	Use:          "config [chave] [valor]",
	Short:        "Gerencia as configurações da CLI",
	SilenceUsage: true,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			var keys []string
			for k := range handlers {
				if strings.HasPrefix(k, toComplete) {
					keys = append(keys, k)
				}
			}
			return keys, cobra.ShellCompDirectiveNoFileComp
		}

		if len(args) == 1 {
			key := args[0]
			if handler, exists := handlers[key]; exists {
				var values []string
				for _, v := range handler.ValidValues {
					if strings.HasPrefix(v, toComplete) {
						values = append(values, v)
					}
				}
				return values, cobra.ShellCompDirectiveNoFileComp
			}
		}

		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return fmt.Errorf("argumentos ausentes ou em excesso.\n\nUso correto:\n  dev config <chave> <valor> [flags]")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.NewConfig(
			config.WithConfigFlags(
				&config.ConfigFlags{
					Global:     globalFlag,
					Interative: interactiveFlag,
				},
			),
		)

		return configImpl(&configImplParams{
			args:   args,
			config: config,
		})
	},
}

func init() {
	configCmd.Flags().BoolVar(&globalFlag, "global", false, "Aplica a configuração no escopo global")
	configCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "i", false, "Abre um menu interativo para seleção de opções válidas")
	rootCmd.AddCommand(configCmd)
}
