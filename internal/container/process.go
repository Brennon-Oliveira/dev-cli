package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Brennon-Oliveira/dev-cli/internal/config"
)

func ExecDetached(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	applyDetachedAttr(cmd)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Process.Release()
}

func RunUpSync(path string) error {
	cmd := exec.Command("devcontainer", "up", "--workspace-folder", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunInteractive(path string, command []string) error {
	args := append([]string{"exec", "--workspace-folder", path}, command...)
	cmd := exec.Command("devcontainer", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ListContainers() error {
	tool := config.Load().Core.Tool
	format := "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Label \"devcontainer.local_folder\"}}"
	cmd := exec.Command(tool, "ps", "--filter", "label=devcontainer.local_folder", "--format", format)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ShowLogs(path string, follow bool) error {
	tool := config.Load().Core.Tool
	filter := fmt.Sprintf("label=devcontainer.local_folder=%s", path)
	idCmd := exec.Command(tool, "ps", "-q", "--filter", filter)
	out, _ := idCmd.Output()
	id := strings.TrimSpace(string(out))

	if id == "" {
		return fmt.Errorf("nenhum container ativo encontrado para o caminho: %s", path)
	}

	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, id)

	cmd := exec.Command(tool, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CleanResources() error {
	tool := config.Load().Core.Tool
	fmt.Println("Removendo containers parados...")
	exec.Command(tool, "container", "prune", "-f").Run()

	fmt.Println("Removendo redes não utilizadas...")
	exec.Command(tool, "network", "prune", "-f").Run()

	fmt.Println("Limpeza concluída.")
	return nil
}

func ListPorts(path string) error {
	tool := config.Load().Core.Tool
	filter := fmt.Sprintf("label=devcontainer.local_folder=%s", path)
	idCmd := exec.Command(tool, "ps", "-q", "--filter", filter)
	out, _ := idCmd.Output()
	id := strings.TrimSpace(string(out))

	if id == "" {
		return fmt.Errorf("nenhum container ativo encontrado para o caminho: %s", path)
	}

	fmt.Printf("Portas mapeadas para o container (%s):\n", id[:12])
	cmd := exec.Command(tool, "port", id)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getContainerIDs(path string) (string, error) {
	tool := config.Load().Core.Tool
	pathsToTry := []string{path}

	// Tenta filtrar tanto pelo path nativo quanto pelo path convertido (WSL -> Windows)
	hostPath := GetHostPath(path)
	if hostPath != path {
		pathsToTry = append(pathsToTry, hostPath)
	}

	for _, p := range pathsToTry {
		filter := fmt.Sprintf("label=devcontainer.local_folder=%s", p)
		cmd := exec.Command(tool, "ps", "-q", "--filter", filter)
		out, _ := cmd.Output()
		id := strings.TrimSpace(string(out))

		if id != "" {
			id = strings.ReplaceAll(id, "\r\n", " ")
			id = strings.ReplaceAll(id, "\n", " ")
			return id, nil
		}
	}

	return "", fmt.Errorf("nenhum container ativo encontrado para o caminho: %s", path)
}

func getAllRelatedContainers(path string) ([]string, error) {
	tool := config.Load().Core.Tool
	pathsToTry := []string{path}

	hostPath := GetHostPath(path)
	if hostPath != path {
		pathsToTry = append(pathsToTry, hostPath)
	}

	var mainIDs []string
	for _, p := range pathsToTry {
		filter := fmt.Sprintf("label=devcontainer.local_folder=%s", p)
		// Usa -a para pegar containers que já estão parados, garantindo que o rm exclua tudo
		cmd := exec.Command(tool, "ps", "-a", "-q", "--filter", filter)
		out, _ := cmd.Output()
		idStr := strings.TrimSpace(string(out))
		if idStr != "" {
			idStr = strings.ReplaceAll(idStr, "\r\n", "\n")
			mainIDs = strings.Split(idStr, "\n")
			break
		}
	}

	if len(mainIDs) == 0 {
		return nil, fmt.Errorf("nenhum container encontrado para o caminho: %s", path)
	}

	allIDsMap := make(map[string]bool)
	for _, id := range mainIDs {
		allIDsMap[id] = true

		// Verifica se o container principal pertence a um compose
		cmd := exec.Command(tool, "inspect", "-f", `{{ if .Config.Labels }}{{ index .Config.Labels "com.docker.compose.project" }}{{ end }}`, id)
		out, _ := cmd.Output()
		project := strings.TrimSpace(string(out))

		if project != "" && project != "<no value>" {
			// Busca todos os containers amarrados a este projeto do compose
			filter := fmt.Sprintf("label=com.docker.compose.project=%s", project)
			cmd = exec.Command(tool, "ps", "-a", "-q", "--filter", filter)
			out2, _ := cmd.Output()
			compIDsStr := strings.TrimSpace(string(out2))
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

	return finalIDs, nil
}

func KillContainer(path string) error {
	tool := config.Load().Core.Tool
	ids, err := getAllRelatedContainers(path)
	if err != nil {
		return err
	}

	fmt.Printf("Forçando parada e excluindo (rm -f) o(s) container(s):\n%s\n", strings.Join(ids, "\n"))
	args := append([]string{"rm", "-f"}, ids...)
	cmd := exec.Command(tool, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DownContainer(path string) error {
	tool := config.Load().Core.Tool
	ids, err := getAllRelatedContainers(path)
	if err != nil {
		return err
	}

	fmt.Printf("Parando graciosamente (stop) o(s) container(s):\n%s\n", strings.Join(ids, "\n"))
	args := append([]string{"stop"}, ids...)
	cmd := exec.Command(tool, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
