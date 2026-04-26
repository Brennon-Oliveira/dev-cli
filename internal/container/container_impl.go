package container

import (
	"fmt"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
)

func (c *realContainerCLI) ListContainersOfActiveDevcontainers() error {
	tool := c.config.Load().Core.Tool
	format := "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Label \"devcontainer.local_folder\"}}"
	output, err := c.executor.Output(tool, "ps", "--filter", "name=devcontainer", "--format", format)

	if err != nil {
		logger.Error("Houve um erro ao ler os DevContainers ativos!")
		return err
	}

	containers := c.parseContainerOutput(string(output))
	formatted := c.formatGroupedContainers(containers)

	logger.Info(formatted)
	return nil
}

func (c *realContainerCLI) CleanResources() error {
	tool := c.config.Load().Core.Tool

	logger.Info("Removendo containers parados...")
	err := c.executor.Run(tool, "container", "prune", "-f")

	if err != nil {
		logger.Error("Houve um erro ao remover os containers parados.")
		logger.Verbose(err.Error())
		return nil
	}

	logger.Info("Removendo redes não utilizadas...")
	err = c.executor.Run(tool, "network", "prune", "-f")

	if err != nil {
		logger.Error("Houve um erro ao remover as redes não utilizadas.")
		logger.Verbose(err.Error())
		return err
	}

	logger.Success("Limpeza concluída.")
	return nil
}

func (c *realContainerCLI) GetAllRelatedContainers(path string) ([]string, error) {
	logger.Info("Procurando containers relacionados ao projeto")
	tool := c.config.Load().Core.Tool
	pathsToTry := []string{path}

	realPath, _ := c.pather.GetRealPath(path)

	if realPath != path {
		pathsToTry = append(pathsToTry, realPath)
	}

	var mainIDs []string
	for _, p := range pathsToTry {
		logger.Verbose("Verificando caminho: %s", p)
		filter := fmt.Sprintf("label=devcontainer.local_folder=%s", p)
		logger.Verbose("Filtro aplicado: %s", filter)

		out, err := c.executor.Output(tool, "ps", "-a", "-q", "--filter", filter)
		if err != nil {
			logger.Error("Houve um erro ao buscar os containers relacionados")
			return nil, err
		}

		idStr := strings.TrimSpace(string(out))
		if idStr != "" {
			idStr = strings.ReplaceAll(idStr, "\r\n", "\n")
			mainIDs = strings.Split(idStr, "\n")
			logger.Verbose("Containers encontrados:")
			logger.Verbose(idStr)
			break
		}
	}

	if len(mainIDs) == 0 {
		err := fmt.Errorf("Nenhum container contrado para o caminho: %s", path)
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("Buscando composes dos containers encontrados")
	allIDsMap := make(map[string]bool)
	for _, id := range mainIDs {
		logger.Verbose("Verificando container '%s'", id)
		allIDsMap[id] = true

		out, err := c.executor.Output(tool, "inspect", "-f", `{{ if .Config.Labels }}{{ index .Config.Labels "com.docker.compose.project" }}{{ end }}`, id)
		if err != nil {
			logger.Error("Houve um erro ao inspecionar o container %s", id)
			return nil, err
		}

		project := strings.TrimSpace(string(out))

		if project == "" || project == "<no value>" {
			logger.Verbose("Sem projeto associado")
			continue
		}

		logger.Verbose("Projeto encontrado: %s", project)
		filter := fmt.Sprintf("label=com.docker.compose.project=%s", project)
		out, err = c.executor.Output(tool, "ps", "-a", "-q", "--filter", filter)

		if err != nil {
			logger.Error("Houve um erro ao buscar os containers do projeto %s", project)
			return nil, err
		}

		compIDsStr := strings.TrimSpace(string(out))

		if compIDsStr == "" {
			continue
		}

		compIDsStr = strings.ReplaceAll(compIDsStr, "\r\n", "\n")

		for _, cid := range strings.Split(compIDsStr, "\n") {
			allIDsMap[cid] = true
		}
	}

	var finalIDs []string
	for id := range allIDsMap {
		if id == "" {
			continue
		}

		finalIDs = append(finalIDs, id)
	}

	return finalIDs, nil
}

func (c *realContainerCLI) DownContainer(path string) error {
	tool := c.config.Load().Core.Tool
	ids, err := c.GetAllRelatedContainers(path)
	if err != nil {
		return err
	}

	logger.Info("Parando graciosamente (stop) o(s) container(s):\n%s\n", strings.Join(ids, "\n"))

	args := append([]string{"stop"}, ids...)

	err = c.executor.Run(tool, args...)

	if err != nil {
		logger.Error("Não foi possível parar os containers.")
		return err
	}

	logger.Info("%d containers parados com sucesso.", len(ids))
	return nil

}
