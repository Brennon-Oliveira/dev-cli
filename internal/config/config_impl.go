package config

func (c *realConfig) Load() GlobalConfig {
	config := GlobalConfig{}

	config.Core.Tool = "docker"

	return config
}
