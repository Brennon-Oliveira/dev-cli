package completer

type Shell string

const (
	Zsh        Shell = "zsh"
	Bash       Shell = "bash"
	PowerShell Shell = "powershell"
)

type Completer interface {
	DetectShell() Shell
	GetHomeDir() (string, error)
	GetDevDir(homeDir string) (string, error)
	InstallInShell(shell Shell, devDir string, homeDir string) error
	AppendToFileIfMissing(filePath, line string) error
}
