package container

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func ExecDetached(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	}
	return cmd.Start()
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
	format := "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Label \"devcontainer.local_folder\"}}"
	cmd := exec.Command("docker", "ps", "--filter", "label=devcontainer.local_folder", "--format", format)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ShowLogs(path string, follow bool) error {
	filter := fmt.Sprintf("label=devcontainer.local_folder=%s", path)
	idCmd := exec.Command("docker", "ps", "-q", "--filter", filter)
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

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CleanResources() error {
	fmt.Println("Removendo containers parados...")
	exec.Command("docker", "container", "prune", "-f").Run()

	fmt.Println("Removendo redes não utilizadas...")
	exec.Command("docker", "network", "prune", "-f").Run()

	fmt.Println("Limpeza concluída.")
	return nil
}

func ListPorts(path string) error {
	filter := fmt.Sprintf("label=devcontainer.local_folder=%s", path)
	idCmd := exec.Command("docker", "ps", "-q", "--filter", filter)
	out, _ := idCmd.Output()
	id := strings.TrimSpace(string(out))

	if id == "" {
		return fmt.Errorf("nenhum container ativo encontrado para o caminho: %s", path)
	}

	fmt.Printf("Portas mapeadas para o container (%s):\n", id[:12])
	cmd := exec.Command("docker", "port", id)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func KillContainer(path string) error {
	filter := fmt.Sprintf("label=devcontainer.local_folder=%s", path)
	idCmd := exec.Command("docker", "ps", "-q", "--filter", filter)
	out, _ := idCmd.Output()
	id := strings.TrimSpace(string(out))

	if id == "" {
		return fmt.Errorf("nenhum container ativo encontrado para o caminho: %s", path)
	}

	fmt.Printf("Encerrando container: %s\n", id[:12])
	cmd := exec.Command("docker", "kill", id)
	return cmd.Run()
}
