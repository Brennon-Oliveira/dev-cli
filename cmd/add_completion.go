package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var addCompletionCmd = &cobra.Command{
	Use:       "add-completion [bash|zsh|powershell]",
	Short:     "Configura o autocompletar da CLI automaticamente no seu shell",
	ValidArgs: []string{"bash", "zsh", "powershell"},
	Args:      cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := ""
		if len(args) > 0 {
			shell = args[0]
		} else {
			shell = detectShell()
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("não foi possível determinar o diretório home: %v", err)
		}

		devDir := filepath.Join(homeDir, ".dev-cli")
		if err := os.MkdirAll(devDir, 0755); err != nil {
			return err
		}

		switch shell {
		case "zsh":
			return installZsh(devDir, homeDir)
		case "bash":
			return installBash(devDir, homeDir)
		case "powershell":
			return installPowerShell(devDir)
		default:
			return fmt.Errorf("shell não suportado ou não detectado: %s. Use explícito: dev add-completion [bash|zsh|powershell]", shell)
		}
	},
}

func detectShell() string {
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

func installZsh(devDir, homeDir string) error {
	compFile := filepath.Join(devDir, "completion.zsh")
	if err := rootCmd.GenZshCompletionFile(compFile); err != nil {
		return err
	}

	rcFile := filepath.Join(homeDir, ".zshrc")
	line := fmt.Sprintf("source %s", compFile)

	if err := appendToFileIfMissing(rcFile, line); err != nil {
		return err
	}
	fmt.Printf("Autocompletar configurado no %s!\nRode 'source %s' ou reinicie o terminal.\n", rcFile, rcFile)
	return nil
}

func installBash(devDir, homeDir string) error {
	compFile := filepath.Join(devDir, "completion.bash")
	if err := rootCmd.GenBashCompletionFile(compFile); err != nil {
		return err
	}

	rcFile := filepath.Join(homeDir, ".bashrc")
	line := fmt.Sprintf("source %s", compFile)

	if err := appendToFileIfMissing(rcFile, line); err != nil {
		return err
	}
	fmt.Printf("Autocompletar configurado no %s!\nRode 'source %s' ou reinicie o terminal.\n", rcFile, rcFile)
	return nil
}

func installPowerShell(devDir string) error {
	compFile := filepath.Join(devDir, "completion.ps1")
	if err := rootCmd.GenPowerShellCompletionFile(compFile); err != nil {
		return err
	}

	out, err := exec.Command("powershell", "-NoProfile", "-Command", "Write-Host $PROFILE").Output()
	if err != nil {
		return fmt.Errorf("falha ao localizar o $PROFILE do PowerShell: %v", err)
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

	fmt.Printf("Autocompletar configurado no %s!\nReinicie o PowerShell para aplicar.\n", profilePath)
	return nil
}

func init() {
	rootCmd.AddCommand(addCompletionCmd)
}
