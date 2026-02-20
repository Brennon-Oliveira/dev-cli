package container

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

func getWorkspaceFolder(absPath string) string {
	cmd := exec.Command("devcontainer", "read-configuration", "--workspace-folder", absPath)
	out, err := cmd.Output()
	if err == nil {
		var config struct {
			Workspace struct {
				WorkspaceFolder string `json:"workspaceFolder"`
			} `json:"workspace"`
		}

		if err := json.Unmarshal(out, &config); err == nil && config.Workspace.WorkspaceFolder != "" {
			return config.Workspace.WorkspaceFolder
		}
	}

	return fmt.Sprintf("/workspaces/%s", filepath.Base(absPath))
}

func GetContainerURI(absPath string) string {
	hostPath := absPath

	if _, isWSL := os.LookupEnv("WSL_DISTRO_NAME"); isWSL {
		cmd := exec.Command("wslpath", "-w", absPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			hostPath = strings.TrimSpace(out.String())
		}
	}

	hexPath := hex.EncodeToString([]byte(hostPath))
	containerPath := getWorkspaceFolder(absPath)

	// Validação para garantir que a URI termine sempre com "//"
	if strings.HasSuffix(containerPath, "/") && !strings.HasSuffix(containerPath, "//") {
		containerPath += "/"
	} else if !strings.HasSuffix(containerPath, "/") {
		containerPath += "//"
	}

	return fmt.Sprintf("vscode-remote://dev-container+%s%s", hexPath, containerPath)
}
