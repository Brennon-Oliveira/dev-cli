package pather

type Pather interface {
	GetAbsPath(target string) (string, error)
	GetRealPath(absPath string) (string, error)
	GetPathFromArgs(args []string) string
}
