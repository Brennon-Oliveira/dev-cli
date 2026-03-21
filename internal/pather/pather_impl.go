package pather

import (
	"bytes"
	"path/filepath"
	"strings"
)

func (p *realPather) GetAbsPath(target string) (string, error) {
	if target == "" {
		target = "."
	}
	return filepath.Abs(target)
}

func (p *realPather) GetPathFromArgs(args []string) string {
	path := ""
	if len(args) > 0 {
		path = args[0]
	}
	return path
}

func (p *realPather) GetRealPath(absPath string) (string, error) {
	if _, isWSL := p.lookupEnv("WSL_DISTRO_NAME"); !isWSL {
		return absPath, nil
	}
	var out bytes.Buffer

	if err := p.executor.RunWithOutput(&out, "wslpath", "-w", absPath); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
