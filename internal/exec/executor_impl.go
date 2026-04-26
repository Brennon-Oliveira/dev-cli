package exec

import (
	"io"
	"os/exec"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
)

func (e *realExecutor) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = e.stdin
	cmd.Stdout = e.stdout
	logger.Verbose("Rodando: %s", strings.Join(append([]string{name}, args...), " "))

	err := cmd.Run()
	logger.Verbose("---")

	// fmt.Printf("%s %s %s %s", name, args[0], args[1], args[2])

	return err
}

func (e *realExecutor) RunWithOutput(output io.Writer, name string, args ...string) error {
	out := io.MultiWriter(e.stdout, output)
	cmd := exec.Command(name, args...)
	cmd.Stdin = e.stdin
	cmd.Stdout = out
	logger.Verbose("Rodando: %s", strings.Join(append([]string{name}, args...), " "))
	return cmd.Run()
}

// CombinedOutput implements [Executor].
func (e *realExecutor) CombinedOutput(name string, args ...string) (string, error) {
	panic("unimplemented")
}

func (e *realExecutor) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	return out, err
}

func (e *realExecutor) RunDetached(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	applyDetachedAttr(cmd)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Process.Release()
}

// RunInteractive implements [Executor].
func (e *realExecutor) RunInteractive(name string, args ...string) error {
	panic("unimplemented")
}
