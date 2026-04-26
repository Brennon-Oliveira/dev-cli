package config

type realConfig struct {
}

type Option func(*realConfig)

func NewConfig(opts ...Option) *realConfig {
	c := &realConfig{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
