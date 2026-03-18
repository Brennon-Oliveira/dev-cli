package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
)

const (
	githubAPIURL   = "https://api.github.com/repos/Brennon-Oliveira/dev-cli/releases/latest"
	githubReleases = "https://github.com/Brennon-Oliveira/dev-cli/releases/download"
)

func (u *GitHubUpdater) CheckForUpdate(currentVersion string) (*ReleaseInfo, error) {
	logs.Info("Verificando atualizações")
	logs.Verbose("consultando: %s", githubAPIURL)

	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar última versão: %w", err)
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("falha ao ler resposta da API: %w", err)
	}

	if release.TagName == currentVersion {
		logs.Info("A CLI já está na última versão (%s)", currentVersion)
		return nil, nil
	}

	logs.Info("Nova versão encontrada: %s (Atual: %s)", release.TagName, currentVersion)

	info := &ReleaseInfo{
		TagName: release.TagName,
		Assets:  make([]Asset, len(release.Assets)),
	}

	for i, a := range release.Assets {
		info.Assets[i] = Asset{Name: a.Name, URL: a.URL}
	}

	return info, nil
}

func (u *GitHubUpdater) Download(assetURL string) (string, error) {
	logs.Info("Baixando atualização...")
	logs.Verbose("URL: %s", assetURL)

	resp, err := http.Get(assetURL)
	if err != nil || resp.StatusCode != 200 {
		return "", fmt.Errorf("falha ao baixar o arquivo: %s", assetURL)
	}
	defer resp.Body.Close()

	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}

	tmpArchive, err := os.CreateTemp("", "dev-cli-update-*"+ext)
	if err != nil {
		return "", err
	}
	defer tmpArchive.Close()

	if _, err := io.Copy(tmpArchive, resp.Body); err != nil {
		return "", err
	}

	logs.Verbose("arquivo baixado: %s", tmpArchive.Name())
	return tmpArchive.Name(), nil
}

func getAssetName(version string) string {
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	return fmt.Sprintf("dev-%s-%s%s", runtime.GOOS, runtime.GOARCH, ext)
}

func findAsset(info *ReleaseInfo) string {
	targetName := getAssetName(info.TagName)
	for _, asset := range info.Assets {
		if asset.Name == targetName {
			return asset.URL
		}
	}
	return ""
}

func getDownloadURL(version string) string {
	return fmt.Sprintf("%s/%s/%s", githubReleases, version, getAssetName(version))
}

func isWindowsArchive(filename string) bool {
	return strings.HasSuffix(filename, ".zip")
}
