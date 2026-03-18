package config

import (
	"fmt"

	"github.com/Brennon-Oliveira/dev-cli/internal/constants"
)

type ConfigHandler struct {
	ValidValues []string
	Label       string
	Get         func(cfg GlobalConfig) string
	Set         func(cfg *GlobalConfig, val string)
}

var Handlers = map[string]ConfigHandler{
	"core.tool": {
		ValidValues: constants.ValidTools,
		Label:       "Selecione o motor de containers padrão",
		Get: func(cfg GlobalConfig) string {
			return cfg.Core.Tool
		},
		Set: func(cfg *GlobalConfig, val string) {
			cfg.Core.Tool = val
		},
	},
	"core.use-sudo": {
		ValidValues: constants.ValidBoolValues,
		Label:       "Usar sudo para comandos Docker/Podman?",
		Get: func(cfg GlobalConfig) string {
			return fmt.Sprintf("%v", cfg.Core.UseSudo)
		},
		Set: func(cfg *GlobalConfig, val string) {
			cfg.Core.UseSudo = val == "true"
		},
	},
}

func GetHandlerKeys() []string {
	keys := make([]string, 0, len(Handlers))
	for k := range Handlers {
		keys = append(keys, k)
	}
	return keys
}
