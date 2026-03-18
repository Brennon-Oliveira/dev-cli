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

	args := d.buildArgs("ps", "--filter", "label="+constants.LabelDevContainerFolder, "--format", format)
	return wrapPermissionError(d.executor.Run(args[0], args[1:]...))
}

func (d *DockerClient) ShowLogs(path string, follow bool) error {
	id, err := d.GetContainerID(path)
	if err != nil {
		return err
	}

	logs.Info("Exibindo logs do container %s", id[:12])

	dockerArgs := []string{"logs"}
	if follow {
		dockerArgs = append(dockerArgs, "-f")
	}
	dockerArgs = append(dockerArgs, id)

	logs.Verbose("executando: %s %s", d.tool, strings.Join(dockerArgs, " "))
	args := d.buildArgs(dockerArgs...)
	return wrapPermissionError(d.executor.Run(args[0], args[1:]...))
}

func (d *DockerClient) ListPorts(path string) error {
	id, err := d.GetContainerID(path)
	if err != nil {
		return err
	}

	logs.Info("Portas mapeadas para o container %s", id[:12])
	logs.Verbose("executando: %s port %s", d.tool, id)

	args := d.buildArgs("port", id)
	return wrapPermissionError(d.executor.Run(args[0], args[1:]...))
}

func (d *DockerClient) CleanResources() error {
	logs.Info("Removendo containers parados")
	logs.Verbose("executando: %s container prune -f", d.tool)

	args := d.buildArgs("container", "prune", "-f")
	if err := wrapPermissionError(d.executor.Run(args[0], args[1:]...)); err != nil {
		return fmt.Errorf("falha ao remover containers parados: %w", err)
	}

	logs.Info("Removendo redes não utilizadas")
	logs.Verbose("executando: %s network prune -f", d.tool)

	args = d.buildArgs("network", "prune", "-f")
	if err := wrapPermissionError(d.executor.Run(args[0], args[1:]...)); err != nil {
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

	dockerArgs := append([]string{"stop"}, ids...)
	logs.Verbose("executando: %s %s", d.tool, strings.Join(dockerArgs, " "))
	args := d.buildArgs(dockerArgs...)
	return wrapPermissionError(d.executor.Run(args[0], args[1:]...))
}

func (d *DockerClient) RemoveContainers(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	logs.Info("Removendo %d container(s)", len(ids))
	for _, id := range ids {
		logs.Verbose("  - %s", id)
	}

	dockerArgs := append([]string{"rm", "-f"}, ids...)
	logs.Verbose("executando: %s %s", d.tool, strings.Join(dockerArgs, " "))
	args := d.buildArgs(dockerArgs...)
	return wrapPermissionError(d.executor.Run(args[0], args[1:]...))
}
