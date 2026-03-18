package vscode

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/constants"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
)

func GetContainerURI(absPath string, workspaceFolder string) string {
	hostPath := paths.GetHostPath(absPath)
	hexPath := hex.EncodeToString([]byte(hostPath))

	containerPath := workspaceFolder
	if containerPath == "" {
		containerPath = constants.DefaultWorkspaceFolder
	}

	if strings.HasSuffix(containerPath, "/") {
	} else if !strings.HasSuffix(containerPath, "workspaces") {
		containerPath += "//"
	}

	return fmt.Sprintf("vscode-remote://dev-container+%s%s", hexPath, containerPath)
}
