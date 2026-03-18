package container

import (
	"fmt"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/constants"
	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
)

func (d *DockerClient) ListContainers() error {
	logs.Info("Listando containers de desenvolvimento")
	format := "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Label \"devcontainer.local_folder\"}}"
	logs.Verbose("executando: %s ps --filter label=%s --format ...", d.tool, constants.LabelDevContainerFolder)
	return d.executor.Run(d.tool, "ps", "--filter", "label="+constants.LabelDevContainerFolder, "--format", format)
}

func (d *DockerClient) ShowLogs(path string, follow bool) error {
	id, err := d.GetContainerID(path)
	if err != nil {
		return err
	}

	logs.Info("Exibindo logs do container %s", id[:12])

	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, id)

	logs.Verbose("executando: %s %s", d.tool, strings.Join(args, " "))
	return d.executor.Run(d.tool, args...)
}

func (d *DockerClient) ListPorts(path string) error {
	id, err := d.GetContainerID(path)
	if err != nil {
		return err
	}

	logs.Info("Portas mapeadas para o container %s", id[:12])
	logs.Verbose("executando: %s port %s", d.tool, id)
	return d.executor.Run(d.tool, "port", id)
}

func (d *DockerClient) CleanResources() error {
	logs.Info("Removendo containers parados")
	logs.Verbose("executando: %s container prune -f", d.tool)
	if err := d.executor.Run(d.tool, "container", "prune", "-f"); err != nil {
		return fmt.Errorf("falha ao remover containers parados: %w", err)
	}

	logs.Info("Removendo redes não utilizadas")
	logs.Verbose("executando: %s network prune -f", d.tool)
	if err := d.executor.Run(d.tool, "network", "prune", "-f"); err != nil {
		return fmt.Errorf("falha ao remover redes: %w", err)
	}

	logs.Success("Limpeza concluída")
	return nil
}

func (d *DockerClient) StopContainers(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	logs.Info("Parando %d container(s)", len(ids))
	for _, id := range ids {
		logs.Verbose("  - %s", id)
	}

	args := append([]string{"stop"}, ids...)
	logs.Verbose("executando: %s %s", d.tool, strings.Join(args, " "))
	return d.executor.Run(d.tool, args...)
}

func (d *DockerClient) RemoveContainers(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	logs.Info("Removendo %d container(s)", len(ids))
	for _, id := range ids {
		logs.Verbose("  - %s", id)
	}

	args := append([]string{"rm", "-f"}, ids...)
	logs.Verbose("executando: %s %s", d.tool, strings.Join(args, " "))
	return d.executor.Run(d.tool, args...)
}
