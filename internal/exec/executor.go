package exec

import "io"

type Executor interface {
	Run(name string, args ...string) error
	RunWithOutput(output io.Writer, name string, args ...string) error
	RunInteractive(name string, args ...string) error
	RunDetached(name string, args ...string) error
	Output(name string, args ...string) ([]byte, error)
	CombinedOutput(name string, args ...string) (string, error)
}
