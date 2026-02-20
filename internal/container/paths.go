package container

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
)

func GetAbsPath(target string) (string, error) {
	if target == "" {
		target = "."
	}
	return filepath.Abs(target)
}

func GetContainerURI(absPath string) string {
	hexPath := hex.EncodeToString([]byte(absPath))
	dirName := filepath.Base(absPath)
	return fmt.Sprintf("vscode-remote://dev-container+%s/workspaces/%s", hexPath, dirName)
}
