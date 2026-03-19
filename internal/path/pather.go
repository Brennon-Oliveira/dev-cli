package path

type Pather interface {
	GetAbsPath(target string) (string, error)
	GetHostPath(absPath string) (string, error)
	GetPathFromArgs(args []string) string
}
