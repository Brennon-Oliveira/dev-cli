package container

import (
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
