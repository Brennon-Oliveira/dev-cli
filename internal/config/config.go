package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Brennon-Oliveira/dev-cli/internal/constants"
)

type GlobalConfig struct {
	Core struct {
		Tool string `json:"tool"`
	} `json:"core"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, constants.ConfigDirName, constants.ConfigFileName), nil
}

func Load() GlobalConfig {
	cfg := GlobalConfig{}
	cfg.Core.Tool = constants.ToolDocker

	path, err := getConfigPath()
	if err != nil {
		return cfg
	}

	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &cfg)
	}

	if cfg.Core.Tool == "" {
		cfg.Core.Tool = constants.ToolDocker
	}

	return cfg
}

func Save(cfg GlobalConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
