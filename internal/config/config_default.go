package config

func getDefaultConfig() *GlobalConfig {
	cfg := &GlobalConfig{}

	cfg.Core.Tool = "docker"

	return cfg
}
