package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/completer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/spf13/cobra"
)

type addCompletionImplParams struct {
	args      []string
	completer completer.Completer
}

func addCompletionImpl(p *addCompletionImplParams) error {
	var shell completer.Shell
	if len(p.args) > 0 {
		shell = completer.Shell(p.args[0])
	} else {
		shell = p.completer.DetectShell()
	}

	homeDir, err := p.completer.GetHomeDir()
	if err != nil {
		return err
	}

	devDir, err := p.completer.GetDevDir(homeDir)
	if err != nil {
		return err
	}

	return p.completer.InstallInShell(shell, devDir, homeDir)
}

var addCompletionCmd = &cobra.Command{
	Use:       "add-completion [bash|zsh|powershell]",
	Short:     "Configura o autocompletar da CLI automaticamente no seu shell",
	ValidArgs: []string{"bash", "zsh", "powershell"},
	Args:      cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := exec.NewExecutor()

		completer := completer.NewCompleter(
			completer.WithExecutor(executor),
			completer.WithGenBashCompletionFile(rootCmd.GenBashCompletionFile),
			completer.WithGenZshCompletionFile(rootCmd.GenZshCompletionFile),
			completer.WithGenPowerShellCompletionFile(rootCmd.GenPowerShellCompletionFile),
		)
		return addCompletionImpl(&addCompletionImplParams{
			args:      args,
			completer: completer,
		})
	},
}

func init() {
	rootCmd.AddCommand(addCompletionCmd)
}
