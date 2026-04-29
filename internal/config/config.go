package config

type ConfigFlags struct {
	Global     bool
	Interative bool
}

type Config interface {
	GetConfigPath() (string, error)
	HasConfigFile() bool
	Load() GlobalConfig
	LoadByKey(key string) string
	TrySave(key string, value string) (string, error)
	Save(key string, value string) error
	InterativeSelect(key string) (string, error)
	ValidateKey(key string) bool
}
