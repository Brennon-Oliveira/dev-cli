package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
)

func (u *GitHubUpdater) Extract(archivePath string) (string, error) {
	logs.Info("Extraindo arquivos...")

	extractedBin := filepath.Join(os.TempDir(), "dev_new")
	if runtime.GOOS == "windows" {
		extractedBin += ".exe"
	}

	var err error
	if isWindowsArchive(archivePath) {
		err = extractZip(archivePath, extractedBin)
	} else {
		err = extractTarGz(archivePath, extractedBin)
	}

	if err != nil {
		return "", fmt.Errorf("falha ao extrair nova versão: %w", err)
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(extractedBin, 0755); err != nil {
			return "", err
		}
	}

	logs.Verbose("binário extraído: %s", extractedBin)
	return extractedBin, nil
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
