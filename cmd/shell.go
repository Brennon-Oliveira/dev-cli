package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal/container"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell [caminho]",
	Short: "Abre um shell interativo dentro do container",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		absPath, _ := container.GetAbsPath(path)

		// Tenta abrir shells em ordem de preferência
		// Se o zsh falhar (container não tem), o erro do RunE será propagado.
		// Para maior robustez, poderíamos encadear tentativas, mas aqui segue o padrão direto.
		return container.RunInteractive(absPath, []string{"/bin/sh", "-c", "if command -v zsh >/dev/null 2>&1; then zsh; elif command -v bash >/dev/null 2>&1; then bash; else sh; fi"})
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
