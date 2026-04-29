package main

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/cmd"
)

func main() {
	if runtime.GOOS != "windows" {
		out, _ := exec.Command("/bin/sh", "-i", "-c", "echo $PATH").Output()
		if realPath := strings.TrimSpace(string(out)); realPath != "" {
			os.Setenv("PATH", realPath)
		}
	}

	cmd.Execute()
}
