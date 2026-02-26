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

	return "/workspaces"
}

// GetHostPath resolve o caminho real considerando o ambiente WSL
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

func GetContainerURI(absPath string) string {
	hostPath := GetHostPath(absPath)
	hexPath := hex.EncodeToString([]byte(hostPath))
	containerPath := getWorkspaceFolder(absPath)

	if strings.HasSuffix(containerPath, "workspaces/") {
		containerPath += "/"
	} else if !strings.HasSuffix(containerPath, "workspaces") {
		containerPath += "//"
	}

	final := fmt.Sprintf("vscode-remote://dev-container+%s%s", hexPath, containerPath)

	fmt.Println(final)

	return final
}
