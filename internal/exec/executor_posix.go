//go:build !windows

package exec

import (
	"os/exec"
	"syscall"
)

func applyDetachedAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
