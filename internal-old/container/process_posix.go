//go:build !windows

package container

import (
	"os/exec"
	"syscall"
)

func applyDetachedAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
