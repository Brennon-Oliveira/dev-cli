package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "[EXPERIMENTAL] Baixa a última versão da CLI e prepara para instalação",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("Verificando atualizações...")

		resp, err := http.Get("https://api.github.com/repos/Brennon-Oliveira/dev-cli/releases/latest")
		if err != nil {
			logger.Error("Falha ao buscar última versão: %v", err)
			return fmt.Errorf("falha ao buscar última versão: %v", err)
		}
		defer resp.Body.Close()

		var release struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			logger.Error("Falha ao ler resposta da API: %v", err)
			return fmt.Errorf("falha ao ler resposta da API: %v", err)
		}

		if release.TagName == Version {
			logger.Info("A CLI já está na última versão (%s).\n", Version)
			return nil
		}

		logger.Info("Nova versão encontrada: %s (Atual: %s)\nBaixando...\n", release.TagName, Version)
		logger.Info("Atualize pela plataforma que instalou!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
