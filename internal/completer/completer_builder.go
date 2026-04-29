package completer

import (
	"os"
	"runtime"

	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
)

type CompletionInstall func(devDir string, homeDir string) error
type Completions map[Shell]CompletionInstall

type realCompleter struct {
	executor                    exec.Executor
	GOOS                        string
	getenv                      func(key string) string
	userHomeDir                 func() (string, error)
	mkdirAll                    func(path string, perm os.FileMode) error
	genBashCompletionFile       func(filename string) error
	genZshCompletionFile        func(filename string) error
	genPowerShellCompletionFile func(filename string) error
	completions                 *Completions
}

type Option func(*realCompleter)

func NewCompleter(opts ...Option) *realCompleter {
	c := &realCompleter{
		GOOS:        runtime.GOOS,
		getenv:      os.Getenv,
		userHomeDir: os.UserHomeDir,
		mkdirAll:    os.MkdirAll,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.completions = &Completions{
		Bash:       c.installBash,
		Zsh:        c.installZsh,
		PowerShell: c.installPowerShell,
	}

	return c
}

func WithExecutor(e exec.Executor) Option {
	return func(c *realCompleter) {
		c.executor = e
	}
}

func WithGOOS(goos string) Option {
	return func(c *realCompleter) {
		c.GOOS = goos
	}
}

func WithGetenv(f func(key string) string) Option {
	return func(c *realCompleter) {
		c.getenv = f
	}
}

func WithUserHomeDir(f func() (string, error)) Option {
	return func(c *realCompleter) {
		c.userHomeDir = f
	}
}

func WithMkdirAll(f func(path string, perm os.FileMode) error) Option {
	return func(c *realCompleter) {
		c.mkdirAll = f
	}
}

func WithGenBashCompletionFile(f func(filename string) error) Option {
	return func(c *realCompleter) {
		c.genBashCompletionFile = f
	}
}

func WithGenZshCompletionFile(f func(filename string) error) Option {
	return func(c *realCompleter) {
		c.genZshCompletionFile = f
	}
}

func WithGenPowerShellCompletionFile(f func(filename string) error) Option {
	return func(c *realCompleter) {
		c.genPowerShellCompletionFile = f
	}
}
