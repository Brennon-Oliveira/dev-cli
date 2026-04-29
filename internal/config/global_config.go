package config

type GlobalConfig struct {
	Core struct {
		Tool string `json:"tool"`
	} `json:"core"`
}

type ConfigHandler struct {
	ValidValues []string
	Label       string
	Get         func(cfg *GlobalConfig) string
	Set         func(cfg *GlobalConfig, val string)
}
