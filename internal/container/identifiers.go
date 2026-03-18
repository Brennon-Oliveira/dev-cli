package container

import (
	"fmt"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/constants"
	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
)

func (d *DockerClient) GetContainerID(path string) (string, error) {
	pathsToTry := []string{path}

	hostPath := paths.GetHostPath(path)
	if hostPath != path {
		pathsToTry = append(pathsToTry, hostPath)
	}

	for _, p := range pathsToTry {
		filter := fmt.Sprintf("label=%s=%s", constants.LabelDevContainerFolder, p)
		logs.Verbose("executando: %s ps -q --filter %s", d.tool, filter)

		out, err := d.executor.Output(d.tool, "ps", "-q", "--filter", filter)
		if err != nil {
			continue
		}

		id := strings.TrimSpace(out)
		if id != "" {
			id = strings.ReplaceAll(id, "\r\n", " ")
			id = strings.ReplaceAll(id, "\n", " ")
			id = strings.TrimSpace(id)
			if id != "" {
				logs.Debug("container ID encontrado: %s", id)
				return id, nil
			}
		}
	}

	return "", fmt.Errorf("nenhum container ativo encontrado para o caminho: %s", path)
}

func (d *DockerClient) InspectLabel(id, label string) (string, error) {
	format := fmt.Sprintf("{{ if .Config.Labels }}{{ index .Config.Labels %q }}{{ end }}", label)
	logs.Verbose("executando: %s inspect -f '%s' %s", d.tool, format, id)

	out, err := d.executor.Output(d.tool, "inspect", "-f", format, id)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
