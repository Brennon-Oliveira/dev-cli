package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/manifoldco/promptui"
)

func (c *realConfig) GetConfigPath() (string, error) {
	home, err := c.userHomeDir()
	if err != nil {
		logger.Error("Não foi possível determinar o diretório home do usuário")
		return "", err
	}

	return filepath.Join(home, ".dev-cli", "config.json"), nil
}

func (c *realConfig) HasConfigFile() bool {
	configPath, err := c.GetConfigPath()

	if err != nil || configPath == "" {
		return false
	}

	if _, err := c.stat(configPath); err != nil || errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (c *realConfig) Load() GlobalConfig {
	cfg := *c.getDefault()

	path, err := c.GetConfigPath()

	if err != nil || path == "" {
		return cfg
	}

	data, err := c.readFile(path)

	if err == nil {
		json.Unmarshal(data, &cfg)
	}

	return cfg
}

func (c *realConfig) LoadByKey(key string) string {
	cfg := c.Load()

	handlers := *c.getHandlers()

	return handlers[key].Get(&cfg)
}

func (c *realConfig) TrySave(key string, value string) (string, error) {
	if !c.flags.Global {
		logger.Error("atualmente apenas a flag --global é suportada")
		return "", fmt.Errorf("atualmente apenas a flag --global é suportada")
	}

	if c.flags.Interative {
		selectedValue, err := c.InterativeSelect(key)

		if err != nil {
			return "", err
		}

		value = selectedValue
	}

	err := c.Save(key, value)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *realConfig) Save(key string, value string) error {
	if !c.flags.Global {
		logger.Error("atualmente apenas a flag --global é suportada")
		return fmt.Errorf("atualmente apenas a flag --global é suportada")
	}

	path, err := c.GetConfigPath()

	if err != nil {
		return err
	}

	if err := c.mkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	cfg := c.Load()

	handlers := *c.getHandlers()

	handler := handlers[key]

	if !c.isAValidValue(&handler, value) {

		var optionsList []string
		for _, validVal := range handler.ValidValues {
			optionsList = append(optionsList, fmt.Sprintf("* %s", validVal))
		}

		logger.Error("Valor inválido para '%s'. Opções permitidas:\n%s", key, strings.Join(optionsList, "\n"))
		return fmt.Errorf("valor inválido para '%s'.\n\nOpções permitidas:\n%s", key, strings.Join(optionsList, "\n"))
	}

	handler.Set(&cfg, value)

	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	return c.writeFile(path, data, 0644)
}

func (c *realConfig) InterativeSelect(key string) (string, error) {
	handlers := *c.getHandlers()
	handler := handlers[key]

	prompt := promptui.Select{
		Label: handler.Label,
		Items: handler.ValidValues,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("seleção cancelada: %v", err)
	}

	return result, nil
}

func (c *realConfig) ValidateKey(key string) bool {
	handlers := *c.getHandlers()

	_, exists := handlers[key]

	if !exists {
		return false
	}

	return true
}
