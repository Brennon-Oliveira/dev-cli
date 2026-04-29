package container_utils

import (
	"fmt"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
)

type ParseContainerOutputFunc func(output string) map[string][]*Container
type FormatGroupedContainersFunc func(grouped map[string][]*Container) string

func ParseContainerOutput(output string) map[string][]*Container {
	logger.Verbose("═══════════════════════════════════════════════════════════════")
	logger.Verbose("PARSING DE OUTPUT DE CONTAINERS")
	logger.Verbose("═══════════════════════════════════════════════════════════════")

	grouped := make(map[string][]*Container)
	lines := strings.Split(output, "\n")
	var containers []*Container

	logger.Verbose(fmt.Sprintf("Total de linhas no output: %d", len(lines)))
	logger.Verbose("")

	logger.Verbose("▶ Parseando linhas...")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		container := parseContainerLine(line)
		if container == nil {
			logger.Verbose(fmt.Sprintf("  ⚠ Linha %d não pôde ser parseada", i))
			continue
		}

		logger.Verbose(fmt.Sprintf("  ✓ Container parseado: %s", container.Names))

		containers = append(containers, container)
	}

	logger.Verbose(fmt.Sprintf("✓ Total de containers parseados: %d", len(containers)))
	logger.Verbose("")

	logger.Verbose("▶ Agrupando containers com local_folder (principais)...")
	for _, container := range containers {
		if container.LocalFolder != "" {
			logger.Verbose(fmt.Sprintf("  ✓ %s -> %s", container.Names, container.LocalFolder))
			grouped[container.LocalFolder] = append(grouped[container.LocalFolder], container)
		}
	}

	logger.Verbose("▶ Agrupando containers sem local_folder (auxiliares)...")
	for _, container := range containers {
		if container.LocalFolder == "" {
			prefix := extractDevcontainerPrefix(container.Names)
			logger.Verbose(fmt.Sprintf("  ℹ %s (prefixo: %s)", container.Names, prefix))

			mainFolder := findMainContainerFolder(prefix, containers)
			if mainFolder != "" {
				logger.Verbose(fmt.Sprintf("    ✓ Agrupado em: %s", mainFolder))
				grouped[mainFolder] = append(grouped[mainFolder], container)
			} else {
				logger.Verbose("    ⚠ Nenhum container principal encontrado")
				grouped["[sem pasta local mapeada]"] = append(grouped["[sem pasta local mapeada]"], container)
			}
		}
	}

	logger.Verbose("")
	logger.Verbose(fmt.Sprintf("✓ Total de grupos criados: %d", len(grouped)))
	logger.Verbose("")
	logger.Verbose("═══════════════════════════════════════════════════════════════")
	logger.Verbose("")
	return grouped
}

func FormatGroupedContainers(grouped map[string][]*Container) string {
	logger.Verbose("═══════════════════════════════════════════════════════════════")
	logger.Verbose("FORMATAÇÃO DOS CONTAINERS")
	logger.Verbose("═══════════════════════════════════════════════════════════════")

	if len(grouped) == 0 {
		logger.Verbose("⚠ Nenhum container encontrado")
		logger.Verbose("")
		return "Nenhum DevContainer ativo encontrado."
	}

	logger.Verbose(fmt.Sprintf("▶ Formatando %d grupo(s) de containers", len(grouped)))
	logger.Verbose("")

	var output strings.Builder
	output.WriteString("Os containers atuais são:\n\n")

	for folder, containers := range grouped {
		logger.Verbose(fmt.Sprintf("▶ Pasta: %s", folder))
		logger.Verbose(fmt.Sprintf("  └─ %d container(s)", len(containers)))

		output.WriteString(folder + "\n")
		output.WriteString("---\n")

		output.WriteString(fmt.Sprintf("%-12s %-45s %-15s %s\n",
			"CONTAINER ID", "NAMES", "STATUS", "FOLDER"))

		for i, container := range containers {
			logger.Verbose(fmt.Sprintf("    %d. %s", i+1, container.Names))
			output.WriteString(fmt.Sprintf("%-12s %-45s %-15s %s\n",
				container.ID, container.Names, container.Status, container.LocalFolder))
		}

		output.WriteString("---\n\n")
	}

	logger.Verbose("")
	logger.Verbose("✓ Formatação concluída com sucesso")
	logger.Verbose("═══════════════════════════════════════════════════════════════\n")
	return output.String()
}

func parseContainerLine(line string) *Container {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil
	}

	id := parts[0]
	names := parts[1]

	var status string
	var localFolder string

	for j := 2; j < len(parts); j++ {
		if strings.HasPrefix(parts[j], "/") {
			localFolder = strings.Join(parts[j:], " ")
			break
		}
		if status != "" {
			status += " "
		}
		status += parts[j]
	}

	return &Container{
		ID:          id,
		Names:       names,
		Status:      status,
		LocalFolder: localFolder,
	}
}

func extractDevcontainerPrefix(names string) string {
	parts := strings.Split(names, "_devcontainer")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func findMainContainerFolder(prefix string, containers []*Container) string {
	for _, container := range containers {
		if container.LocalFolder != "" && strings.HasPrefix(container.Names, prefix) {
			return container.LocalFolder
		}
	}
	return ""
}
