package cmd

import (
	"github.com/Brennon-Oliveira/dev-cli/internal-old/container"
	"github.com/Brennon-Oliveira/dev-cli/internal/path"
	"github.com/spf13/cobra"
)

type runImplParams struct {
	args   []string
	pather path.Pather
}

func runImpl(params *runImplParams) error {
	args := params.args
	pather := params.pather
	path := pather.GetPathFromArgs(args)
	absPath, _ := pather.GetAbsPath(path)

	if err := container.RunUpSync(absPath); err != nil {
		return err
	}

	uri := container.GetContainerURI(absPath)
	return container.ExecDetached("code", "--folder-uri", uri)
}

var runCmd = &cobra.Command{
	Use:   "run [caminho]",
	Short: "Sobe o container e abre o VS Code",
	Long:  "Executa a rotina completa de inicialização: provisiona o container (equivalente ao 'up') e imediatamente anexa o VS Code ao ambiente remoto. Resolve dinamicamente o workspaceFolder e contorna falhas de URI em ambientes como WSL.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		pather := path.NewPather()

		return runImpl(&runImplParams{
			args:   args,
			pather: pather,
		})
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
