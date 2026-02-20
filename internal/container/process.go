package container

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func ExecDetached(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	}
	return cmd.Start()
}

func RunUpSync(path string) error {
	cmd := exec.Command("devcontainer", "up", "--workspace-folder", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
