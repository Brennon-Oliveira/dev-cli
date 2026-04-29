package config

import (
	"os"
)

type realConfig struct {
	userHomeDir   func() (string, error)
	readFile      func(name string) ([]byte, error)
	mkdirAll      func(path string, perm os.FileMode) error
	writeFile     func(name string, data []byte, perm os.FileMode) error
	stat          func(name string) (os.FileInfo, error)
	getDefault    func() *GlobalConfig
	getHandlers   func() *map[string]ConfigHandler
	isAValidValue func(handler *ConfigHandler, value string) bool
	flags         *ConfigFlags
}

type Option func(*realConfig)

func NewConfig(opts ...Option) *realConfig {
	c := &realConfig{
		userHomeDir:   os.UserHomeDir,
		readFile:      os.ReadFile,
		mkdirAll:      os.MkdirAll,
		writeFile:     os.WriteFile,
		stat:          os.Stat,
		getDefault:    getDefaultConfig,
		getHandlers:   GetHandlers,
		isAValidValue: IsAValidValue,
		flags: &ConfigFlags{
			Global:     false,
			Interative: false,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithUserHomeDir(f func() (string, error)) Option {
	return func(c *realConfig) {
		c.userHomeDir = f
	}
}

func WithReadFile(f func(name string) ([]byte, error)) Option {
	return func(c *realConfig) {
		c.readFile = f
	}
}

func WithMkdirAll(f func(path string, perm os.FileMode) error) Option {
	return func(c *realConfig) {
		c.mkdirAll = f
	}
}

func WithWriteFile(f func(name string, data []byte, perm os.FileMode) error) Option {
	return func(c *realConfig) {
		c.writeFile = f
	}
}

func WithStat(f func(name string) (os.FileInfo, error)) Option {
	return func(c *realConfig) {
		c.stat = f
	}
}

func WithGetDefault(f func() *GlobalConfig) Option {
	return func(c *realConfig) {
		c.getDefault = f
	}
}

func WithGetHandlers(f func() *map[string]ConfigHandler) Option {
	return func(c *realConfig) {
		c.getHandlers = f
	}
}

func WithConfigFlags(flags *ConfigFlags) Option {
	return func(c *realConfig) {
		c.flags = flags
	}
}
