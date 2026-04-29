package container_utils

import (
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
)

type TryPathsFunc func(path string, pth pather.Pather) []string
type FindMainContainersForPathFunc func(tool string, p string, executor exec.Executor) ([]string, error)
type ExtractProjectFromContainerFunc func(tool string, id string, executor exec.Executor) (string, error)
type FindComposeContainersForProjectFunc func(tool string, project string, executor exec.Executor) ([]string, error)
type DeduplicateAndFilterContainerIDsFunc func(idMap map[string]bool) []string

func TryPaths(path string, pth pather.Pather) []string {
	pathsToTry := []string{path}

	realPath, _ := pth.GetRealPath(path)

	if realPath != path {
		pathsToTry = append(pathsToTry, realPath)
	}

	return pathsToTry
}

func FindMainContainersForPath(tool string, p string, executor exec.Executor) ([]string, error) {
	logger.Verbose("Verificando caminho: %s", p)
	filter := "label=devcontainer.local_folder=" + p
	logger.Verbose("Filtro aplicado: %s", filter)

	out, err := executor.Output(tool, "ps", "-a", "-q", "--filter", filter)
	if err != nil {
		logger.Error("Houve um erro ao buscar os containers relacionados")
		return nil, err
	}

	idStr := strings.TrimSpace(string(out))
	if idStr == "" {
		return nil, nil
	}

	idStr = strings.ReplaceAll(idStr, "\r\n", "\n")
	mainIDs := strings.Split(idStr, "\n")
	logger.Verbose("Containers encontrados:")
	logger.Verbose(idStr)
	return mainIDs, nil
}

func ExtractProjectFromContainer(tool string, id string, executor exec.Executor) (string, error) {
	logger.Verbose("Verificando container '%s'", id)

	out, err := executor.Output(tool, "inspect", "-f", `{{ if .Config.Labels }}{{ index .Config.Labels "com.docker.compose.project" }}{{ end }}`, id)
	if err != nil {
		logger.Error("Houve um erro ao inspecionar o container %s", id)
		return "", err
	}

	project := strings.TrimSpace(string(out))

	if project == "" || project == "<no value>" {
		logger.Verbose("Sem projeto associado")
		return "", nil
	}

	logger.Verbose("Projeto encontrado: %s", project)
	return project, nil
}

func FindComposeContainersForProject(tool string, project string, executor exec.Executor) ([]string, error) {
	filter := "label=com.docker.compose.project=" + project
	out, err := executor.Output(tool, "ps", "-a", "-q", "--filter", filter)

	if err != nil {
		logger.Error("Houve um erro ao buscar os containers do projeto %s", project)
		return nil, err
	}

	compIDsStr := strings.TrimSpace(string(out))

	if compIDsStr == "" {
		return nil, nil
	}

	compIDsStr = strings.ReplaceAll(compIDsStr, "\r\n", "\n")
	compIDs := strings.Split(compIDsStr, "\n")

	return compIDs, nil
}

func DeduplicateAndFilterContainerIDs(idMap map[string]bool) []string {
	var finalIDs []string
	for id := range idMap {
		if id == "" {
			continue
		}

		finalIDs = append(finalIDs, id)
	}

	return finalIDs
}
