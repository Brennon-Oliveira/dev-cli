package path

import "path/filepath"

type realPather struct{}

func NewPather() *realPather {
	return &realPather{}
}

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

func (p *realPather) GetHostPath(string) (string, error) {
	return "", nil
}
