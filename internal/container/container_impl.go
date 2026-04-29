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

	pathsToTry := c.tryPaths(path, c.pather)

	var mainIDs []string
	for _, p := range pathsToTry {
		ids, err := c.findMainContainersForPath(tool, p, c.executor)
		if err != nil {
			return nil, err
		}

		if len(ids) > 0 {
			mainIDs = ids
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
		allIDsMap[id] = true

		project, err := c.extractProjectFromContainer(tool, id, c.executor)
		if err != nil {
			return nil, err
		}

		if project == "" {
			continue
		}

		compIDs, err := c.findComposeContainersForProject(tool, project, c.executor)
		if err != nil {
			return nil, err
		}

		for _, cid := range compIDs {
			allIDsMap[cid] = true
		}
	}

	finalIDs := c.deduplicateAndFilterContainerIDs(allIDsMap)
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

func (c *realContainerCLI) KillContainer(path string) error {

	tool := c.config.Load().Core.Tool
	ids, err := c.GetAllRelatedContainers(path)
	if err != nil {
		return err
	}

	logger.Info("Forçando a parada e excluindo (rm -f) o(s) container(s):\n%s", strings.Join(ids, "\n"))
	args := append([]string{"rm", "-f"}, ids...)

	err = c.executor.Run(tool, args...)

	if err != nil {
		logger.Error("Não foi possível forçar a parada e remoção dos containers.")
		return err
	}

	logger.Info("%d containers removidos com sucesso.", len(ids))

	return nil
}

func (c *realContainerCLI) ShowLogs(path string, follow bool) error {
	tool := c.config.Load().Core.Tool
	filter := fmt.Sprintf("label=devcontainer.local_folder=%s", path)

	getIdArgs := []string{"ps", "-q", "--filter", filter}
	out, err := c.executor.Output(tool, getIdArgs...)
	if err != nil {
		logger.Error("Não foi possível obter os containers para mostrar os logs.")
		return err
	}

	id := strings.TrimSpace(string(out))

	if id == "" {
		logger.Error("Nenhum container encontrado para o caminho especificado.")
		return fmt.Errorf("nenhum container encontrado para o caminho: %s", path)
	}

	logger.Info("Logs do container %s:", id)

	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, id)

	err = c.executor.Run(tool, args...)

	if err != nil {
		logger.Error("Não foi possível mostrar os logs do container.")
		return err
	}

	return nil
}
