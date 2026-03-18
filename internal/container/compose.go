package container

import (
	"fmt"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/constants"
	"github.com/Brennon-Oliveira/dev-cli/internal/logs"
	"github.com/Brennon-Oliveira/dev-cli/internal/paths"
)

func (d *DockerClient) GetAllRelatedContainers(path string) ([]string, error) {
	pathsToTry := []string{path}

	hostPath := paths.GetHostPath(path)
	if hostPath != path {
		pathsToTry = append(pathsToTry, hostPath)
	}

	var mainIDs []string
	for _, p := range pathsToTry {
		filter := fmt.Sprintf("label=%s=%s", constants.LabelDevContainerFolder, p)
		logs.Verbose("executando: %s ps -a -q --filter %s", d.tool, filter)

		out, err := d.executor.Output(d.tool, "ps", "-a", "-q", "--filter", filter)
		if err != nil {
			continue
		}

		idStr := strings.TrimSpace(out)
		if idStr != "" {
			idStr = strings.ReplaceAll(idStr, "\r\n", "\n")
			mainIDs = strings.Split(idStr, "\n")
			break
		}
	}

	if len(mainIDs) == 0 {
		return nil, fmt.Errorf("nenhum container encontrado para o caminho: %s", path)
	}

	logs.Debug("encontrados %d container(s) principal(is)", len(mainIDs))

	allIDsMap := make(map[string]bool)
	for _, id := range mainIDs {
		allIDsMap[id] = true

		project, err := d.InspectLabel(id, constants.LabelComposeProject)
		if err != nil {
			continue
		}

		if project != "" && project != "<no value>" {
			logs.Verbose("container %s pertence ao compose project: %s", id, project)
			filter := fmt.Sprintf("label=%s=%s", constants.LabelComposeProject, project)
			logs.Verbose("executando: %s ps -a -q --filter %s", d.tool, filter)

			out, err := d.executor.Output(d.tool, "ps", "-a", "-q", "--filter", filter)
			if err != nil {
				continue
			}

			compIDsStr := strings.TrimSpace(out)
			if compIDsStr != "" {
				compIDsStr = strings.ReplaceAll(compIDsStr, "\r\n", "\n")
				for _, cid := range strings.Split(compIDsStr, "\n") {
					allIDsMap[cid] = true
				}
			}
		}
	}

	var finalIDs []string
	for id := range allIDsMap {
		if id != "" {
			finalIDs = append(finalIDs, id)
		}
	}

	logs.Debug("total de containers relacionados: %d", len(finalIDs))
	return finalIDs, nil
}
