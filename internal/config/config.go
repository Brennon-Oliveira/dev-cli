package config

type Config interface {
	Load() GlobalConfig
}
