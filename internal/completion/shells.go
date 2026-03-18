package completion

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/constants"
	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
)

func (c *CompletionInstaller) Install(shell string) error {
	if shell == "" {
		shell = c.DetectShell()
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("não foi possível determinar o diretório home: %w", err)
	}

	devDir := filepath.Join(homeDir, constants.ConfigDirName)
	if err := os.MkdirAll(devDir, 0755); err != nil {
		return err
	}

	switch shell {
	case "zsh":
		return c.installZsh(devDir, homeDir)
	case "bash":
		return c.installBash(devDir, homeDir)
	case "powershell":
		return c.installPowerShell(devDir)
	default:
		return fmt.Errorf("shell não suportado ou não detectado: %s. Use: dev add-completion [bash|zsh|powershell]", shell)
	}
}

func (c *CompletionInstaller) DetectShell() string {
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	shellEnv := os.Getenv("SHELL")
	if strings.Contains(shellEnv, "zsh") {
		return "zsh"
	}
	if strings.Contains(shellEnv, "bash") {
		return "bash"
	}
	return ""
}

func (c *CompletionInstaller) installZsh(devDir, homeDir string) error {
	compFile := filepath.Join(devDir, "completion.zsh")
	logs.Info("Gerando completion para zsh")
	logs.Verbose("arquivo: %s", compFile)

	if err := c.rootCmd.GenZshCompletionFile(compFile); err != nil {
		return err
	}

	rcFile := filepath.Join(homeDir, ".zshrc")
	line := fmt.Sprintf("source %s", compFile)

	if err := appendToFileIfMissing(rcFile, line); err != nil {
		return err
	}

	logs.Success("Autocompletar configurado no %s", rcFile)
	logs.Info("Rode 'source %s' ou reinicie o terminal", rcFile)
	return nil
}

func (c *CompletionInstaller) installBash(devDir, homeDir string) error {
	compFile := filepath.Join(devDir, "completion.bash")
	logs.Info("Gerando completion para bash")
	logs.Verbose("arquivo: %s", compFile)

	if err := c.rootCmd.GenBashCompletionFile(compFile); err != nil {
		return err
	}

	rcFile := filepath.Join(homeDir, ".bashrc")
	line := fmt.Sprintf("source %s", compFile)

	if err := appendToFileIfMissing(rcFile, line); err != nil {
		return err
	}

	logs.Success("Autocompletar configurado no %s", rcFile)
	logs.Info("Rode 'source %s' ou reinicie o terminal", rcFile)
	return nil
}

func (c *CompletionInstaller) installPowerShell(devDir string) error {
	compFile := filepath.Join(devDir, "completion.ps1")
	logs.Info("Gerando completion para PowerShell")
	logs.Verbose("arquivo: %s", compFile)

	if err := c.rootCmd.GenPowerShellCompletionFile(compFile); err != nil {
		return err
	}

	out, err := exec.Command("powershell", "-NoProfile", "-Command", "Write-Host $PROFILE").Output()
	if err != nil {
		return fmt.Errorf("falha ao localizar o $PROFILE do PowerShell: %w", err)
	}

	profilePath := strings.TrimSpace(string(out))
	if profilePath == "" {
		return fmt.Errorf("caminho do $PROFILE retornou vazio")
	}

	if err := os.MkdirAll(filepath.Dir(profilePath), 0755); err != nil {
		return err
	}

	line := fmt.Sprintf(". %s", compFile)
	if err := appendToFileIfMissing(profilePath, line); err != nil {
		return err
	}

	logs.Success("Autocompletar configurado no %s", profilePath)
	logs.Info("Reinicie o PowerShell para aplicar")
	return nil
}

func appendToFileIfMissing(filePath, line string) error {
	content, err := os.ReadFile(filePath)
	if err == nil && strings.Contains(string(content), line) {
		return nil
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("\n# Dev CLI Autocompletion\n%s\n", line))
	return err
}
