package paths

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetAbsPath(target string) (string, error) {
	if target == "" {
		target = "."
	}
	return filepath.Abs(target)
}

func GetHostPath(absPath string) string {
	hostPath := absPath
	if _, isWSL := os.LookupEnv("WSL_DISTRO_NAME"); isWSL {
		cmd := exec.Command("wslpath", "-w", absPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			hostPath = strings.TrimSpace(out.String())
		}
	}
	return hostPath
}
