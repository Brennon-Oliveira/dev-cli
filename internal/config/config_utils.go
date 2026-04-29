package config

import (
	"fmt"
)

var handlers = map[string]ConfigHandler{
	"core.tool": {
		ValidValues: []string{"docker", "podman"},
		Label:       "Selecione o motor de containers padrão",
		Get: func(cfg *GlobalConfig) string {
			return cfg.Core.Tool
		},
		Set: func(cfg *GlobalConfig, val string) {
			cfg.Core.Tool = val
		},
	},
}

func GetHandlers() *map[string]ConfigHandler {
	return &handlers
}

func IsAValidValue(handler *ConfigHandler, value string) bool {
	isValid := false
	var optionsList []string
	for _, validVal := range handler.ValidValues {
		optionsList = append(optionsList, fmt.Sprintf("* %s", validVal))
		if value == validVal {
			isValid = true
		}
	}

	if !isValid {
		return false
	}

	return true
}
