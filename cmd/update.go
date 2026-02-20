package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "[EXPERIMENTAL] Baixa a última versão da CLI e prepara para instalação",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Verificando atualizações...")

		resp, err := http.Get("https://api.github.com/repos/Brennon-Oliveira/dev-cli/releases/latest")
		if err != nil {
			return fmt.Errorf("falha ao buscar última versão: %v", err)
		}
		defer resp.Body.Close()

		var release struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return fmt.Errorf("falha ao ler resposta da API: %v", err)
		}

		if release.TagName == Version {
			fmt.Printf("A CLI já está na última versão (%s).\n", Version)
			return nil
		}

		fmt.Printf("Nova versão encontrada: %s (Atual: %s)\nBaixando...\n", release.TagName, Version)

		ext := ".tar.gz"
		if runtime.GOOS == "windows" {
			ext = ".zip"
		}

		fileName := fmt.Sprintf("dev-%s-%s%s", runtime.GOOS, runtime.GOARCH, ext)
		downloadURL := fmt.Sprintf("https://github.com/Brennon-Oliveira/dev-cli/releases/download/%s/%s", release.TagName, fileName)

		tmpArchive, err := os.CreateTemp("", "dev-cli-update-*"+ext)
		if err != nil {
			return err
		}
		defer os.Remove(tmpArchive.Name())
		defer tmpArchive.Close()

		fileResp, err := http.Get(downloadURL)
		if err != nil || fileResp.StatusCode != 200 {
			return fmt.Errorf("falha ao baixar o arquivo: %s", downloadURL)
		}
		defer fileResp.Body.Close()

		if _, err := io.Copy(tmpArchive, fileResp.Body); err != nil {
			return err
		}
		tmpArchive.Close()

		// Define o caminho temporário para o binário extraído
		tempDir := os.TempDir()
		extractedBin := filepath.Join(tempDir, "dev_new")
		if runtime.GOOS == "windows" {
			extractedBin += ".exe"
		}

		fmt.Println("Extraindo arquivos...")
		if ext == ".zip" {
			err = extractZip(tmpArchive.Name(), extractedBin)
		} else {
			err = extractTarGz(tmpArchive.Name(), extractedBin)
		}

		if err != nil {
			return fmt.Errorf("falha ao extrair nova versão: %v", err)
		}

		if runtime.GOOS != "windows" {
			os.Chmod(extractedBin, 0755)
		}

		execPath, err := os.Executable()
		if err != nil {
			return err
		}

		fmt.Println("\n✅ Download e extração concluídos!")
		fmt.Println("Para finalizar a atualização, copie e execute o comando abaixo no seu terminal:")
		fmt.Println(strings.Repeat("-", 60))

		if runtime.GOOS == "windows" {
			fmt.Printf("Move-Item -Force '%s' '%s'\n", extractedBin, execPath)
		} else {
			// No Linux/macOS, pode ser necessário sudo dependendo de onde o binário atual está
			fmt.Printf("sudo mv -f %s %s\n", extractedBin, execPath)
		}

		fmt.Println(strings.Repeat("-", 60))
		return nil
	},
}

func extractZip(zipPath, targetPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "dev.exe") || f.Name == "dev" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			return err
		}
	}
	return fmt.Errorf("executável não encontrado dentro do zip")
}

func extractTarGz(tarGzPath, targetPath string) error {
	file, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg && (strings.HasSuffix(header.Name, "dev") || strings.HasSuffix(header.Name, "dev.exe")) {
			outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tr)
			return err
		}
	}
	return fmt.Errorf("executável não encontrado dentro do tar.gz")
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
