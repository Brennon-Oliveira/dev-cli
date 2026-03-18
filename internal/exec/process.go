package exec

import (
	"os"
	"os/exec"
)

func (e *RealExecutor) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (e *RealExecutor) RunInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (e *RealExecutor) RunDetached(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	applyDetachedAttr(cmd)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Process.Release()
}

func (e *RealExecutor) Output(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	return string(out), err
}

func (e *RealExecutor) CombinedOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
