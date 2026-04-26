package config

type GlobalConfig struct {
	Core struct {
		Tool string `json:"tool"`
	} `json:"core"`
}
