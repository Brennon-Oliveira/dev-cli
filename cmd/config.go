package cmd

import (
	"fmt"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var globalFlag bool
var interactiveFlag bool

type configHandler struct {
	ValidValues []string
	Label       string
	Get         func(cfg config.GlobalConfig) string
	Set         func(cfg *config.GlobalConfig, val string)
}

var handlers = map[string]configHandler{
	"core.tool": {
		ValidValues: []string{"docker", "podman"},
		Label:       "Selecione o motor de containers padrão",
		Get: func(cfg config.GlobalConfig) string {
			return cfg.Core.Tool
		},
		Set: func(cfg *config.GlobalConfig, val string) {
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
			return fmt.Errorf("argumentos ausentes ou em excesso.\n\nUso correto:\n  dev config [chave] [valor] [flags]")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !globalFlag {
			return fmt.Errorf("atualmente apenas a flag --global é suportada")
		}

		key := args[0]
		handler, exists := handlers[key]
		if !exists {
			var keys []string
			for k := range handlers {
				keys = append(keys, fmt.Sprintf("* %s", k))
			}
			return fmt.Errorf("chave desconhecida: %s.\n\nChaves suportadas:\n%s", key, strings.Join(keys, "\n"))
		}

		cfg := config.Load()

		if interactiveFlag {
			prompt := promptui.Select{
				Label: handler.Label,
				Items: handler.ValidValues,
			}

			_, result, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("seleção cancelada: %v", err)
			}

			handler.Set(&cfg, result)
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("erro ao salvar configuração: %v", err)
			}

			fmt.Printf("Configuração '%s' atualizada para: %s\n", key, result)
			return nil
		}

		if len(args) == 2 {
			val := args[1]

			isValid := false
			var optionsList []string
			for _, validVal := range handler.ValidValues {
				optionsList = append(optionsList, fmt.Sprintf("* %s", validVal))
				if val == validVal {
					isValid = true
				}
			}

			if !isValid {
				return fmt.Errorf("valor inválido para '%s'.\n\nOpções permitidas:\n%s", key, strings.Join(optionsList, "\n"))
			}

			handler.Set(&cfg, val)
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("erro ao salvar configuração: %v", err)
			}

			fmt.Printf("Configuração '%s' atualizada para: %s\n", key, val)
			return nil
		}

		fmt.Println(handler.Get(cfg))
		return nil
	},
}

func init() {
	configCmd.Flags().BoolVar(&globalFlag, "global", false, "Aplica a configuração no escopo global")
	configCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "i", false, "Abre um menu interativo para seleção de opções válidas")
	rootCmd.AddCommand(configCmd)
}
