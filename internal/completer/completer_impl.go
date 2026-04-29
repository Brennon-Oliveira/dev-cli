package completer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
)

func (c *realCompleter) DetectShell() Shell {
	if c.GOOS == "windows" {
		return "powershell"
	}
	shellEnv := c.getenv("SHELL")
	if strings.Contains(shellEnv, "zsh") {
		return "zsh"
	}
	if strings.Contains(shellEnv, "bash") {
		return "bash"
	}
	return ""
}

func (c *realCompleter) GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Não foi possível determinar o diretório home")
		return "", fmt.Errorf("não foi possível determinar o diretório home: %v", err)
	}

	return homeDir, nil
}

func (c *realCompleter) GetDevDir(homeDir string) (string, error) {
	devDir := filepath.Join(homeDir, ".dev-cli")
	if err := c.mkdirAll(devDir, 0755); err != nil {
		logger.Error("Erro ao criar diretório de configuração da CLI")
		return "", err
	}

	return devDir, nil
}

func (c *realCompleter) AppendToFileIfMissing(filePath, line string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), line) {
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

func (c *realCompleter) InstallInShell(shell Shell, devDir string, homeDir string) error {
	install, exists := (*c.completions)[shell]
	if !exists {
		message := fmt.Sprintf("Shell não suportado ou não detectado: %s. Use explícito: dev add-completion [%s]", shell, strings.Join([]string{
			string(Bash),
			string(Zsh),
			string(PowerShell),
		}, "|"))
		logger.Error(message)
		return fmt.Errorf("%s", message)
	}

	return install(devDir, homeDir)
}

func (c *realCompleter) installBash(devDir string, homeDir string) error {
	compFile := filepath.Join(devDir, "completion.bash")

	if err := c.genBashCompletionFile(compFile); err != nil {
		logger.Error("Erro ao gerar arquivo de autocompletar para bash")
		return err
	}

	rcFile := filepath.Join(homeDir, ".bashrc")
	line := fmt.Sprintf("source %s", compFile)

	if err := c.AppendToFileIfMissing(rcFile, line); err != nil {
		logger.Error("Erro ao adicionar linha ao arquivo de configuração")
		return err
	}

	logger.Info("Autocompletar configurado no %s!\nRode 'source %s' ou reinicie o terminal.\n", rcFile, rcFile)
	return nil
}

func (c *realCompleter) installZsh(devDir string, homeDir string) error {
	compFile := filepath.Join(devDir, "completion.zsh")

	if err := c.genZshCompletionFile(compFile); err != nil {
		logger.Error("Erro ao gerar arquivo de autocompletar para zsh")
		return err
	}

	rcFile := filepath.Join(homeDir, ".zshrc")
	line := fmt.Sprintf("source %s", compFile)

	if err := c.AppendToFileIfMissing(rcFile, line); err != nil {
		logger.Error("Erro ao adicionar linha ao arquivo de configuração")
		return err
	}

	logger.Info("Autocompletar configurado no %s!\nRode 'source %s' ou reinicie o terminal.\n", rcFile, rcFile)
	return nil
}

func (c *realCompleter) installPowerShell(devDir string, homeDir string) error {
	compFile := filepath.Join(devDir, "completion.ps1")

	if err := c.genPowerShellCompletionFile(compFile); err != nil {
		logger.Error("Erro ao gerar arquivo de autocompletar para zsh")
		return err
	}

	out, err := c.executor.Output("powershell", "-NoProfile", "-Command", "Write-Host $PROFILE")
	if err != nil {
		logger.Error("Erro ao localizar o $PROFILE do PowerShell")
		return fmt.Errorf("falha ao localizar o $PROFILE do PowerShell: %v", err)
	}

	profilePath := strings.TrimSpace(string(out))
	if profilePath == "" {
		logger.Error("Caminho do $PROFILE retornou vazio")
		return fmt.Errorf("caminho do $PROFILE retornou vazio")
	}

	if err := c.mkdirAll(filepath.Dir(profilePath), 0755); err != nil {
		logger.Error("Erro ao criar diretório para o $PROFILE do PowerShell")
		return err
	}

	line := fmt.Sprintf(". %s", compFile)
	if err := c.AppendToFileIfMissing(profilePath, line); err != nil {
		logger.Error("Erro ao adicionar linha ao $PROFILE do PowerShell")
		return err
	}

	logger.Info("Autocompletar configurado no $PROFILE do PowerShell!\nRode 'powershell -NoProfile' para testar ou reinicie o terminal.\n")

	return nil
}
