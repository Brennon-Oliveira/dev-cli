package container

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func writeExecutable(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
	return path
}

func setupPath(t *testing.T) string {
	t.Helper()
	bin := t.TempDir()
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return bin
}

func TestGetAbsPath_DefaultAndProvided(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	gotDefault, err := GetAbsPath("")
	if err != nil {
		t.Fatalf("GetAbsPath default returned error: %v", err)
	}
	if gotDefault != cwd {
		t.Fatalf("expected %q, got %q", cwd, gotDefault)
	}

	gotProvided, err := GetAbsPath(".")
	if err != nil {
		t.Fatalf("GetAbsPath provided returned error: %v", err)
	}
	if gotProvided != cwd {
		t.Fatalf("expected %q, got %q", cwd, gotProvided)
	}
}

func TestGetHostPath_NoWSLReturnsOriginalPath(t *testing.T) {
	t.Setenv("WSL_DISTRO_NAME", "")
	path := "/tmp/my-project"
	if got := GetHostPath(path); got != path {
		t.Fatalf("expected same path %q, got %q", path, got)
	}
}

func TestGetHostPath_WithWSLUsesWslpath(t *testing.T) {
	bin := setupPath(t)
	t.Setenv("WSL_DISTRO_NAME", "Ubuntu")

	writeExecutable(t, bin, "wslpath", "#!/bin/sh\nprintf 'C:\\\\Repo\\\\project\\n'\n")

	got := GetHostPath("/home/user/project")
	want := "C:\\Repo\\project"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestGetContainerURI_FallbackWorkspace(t *testing.T) {
	absPath := "/tmp/project"
	uri := GetContainerURI(absPath)

	hexPath := hex.EncodeToString([]byte(absPath))
	want := "vscode-remote://dev-container+" + hexPath + "/workspaces"
	if uri != want {
		t.Fatalf("expected %q, got %q", want, uri)
	}
}

func TestGetContainerURI_AppendsDoubleSlashOutsideWorkspaces(t *testing.T) {
	bin := setupPath(t)

	writeExecutable(t, bin, "devcontainer", "#!/bin/sh\nif [ \"$1\" = \"read-configuration\" ]; then\n  printf '{\"workspace\":{\"workspaceFolder\":\"/workspace/custom\"}}'\n  exit 0\nfi\nexit 1\n")

	absPath := "/tmp/project"
	uri := GetContainerURI(absPath)

	hexPath := hex.EncodeToString([]byte(absPath))
	want := "vscode-remote://dev-container+" + hexPath + "/workspace/custom//"
	if uri != want {
		t.Fatalf("expected %q, got %q", want, uri)
	}
}

func TestShowLogs_ReturnsErrorWhenNoContainerFound(t *testing.T) {
	bin := setupPath(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeExecutable(t, bin, "docker", "#!/bin/sh\nif [ \"$1\" = \"ps\" ]; then\n  exit 0\nfi\nexit 0\n")

	err := ShowLogs("/tmp/project", false)
	if err == nil {
		t.Fatal("expected error when no container ID is found")
	}
	if !strings.Contains(err.Error(), "nenhum container ativo encontrado") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetContainerIDs_NormalizesCRLFAndUsesHostPathFallback(t *testing.T) {
	bin := setupPath(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("WSL_DISTRO_NAME", "Ubuntu")

	writeExecutable(t, bin, "wslpath", "#!/bin/sh\nprintf 'C:\\\\work\\\\project\\n'\n")
	writeExecutable(t, bin, "docker", `#!/bin/sh
if [ "$1" = "ps" ] && [ "$2" = "-q" ]; then
  case "$*" in
    *'/home/user/project'*)
      ;;
    *)
      printf 'id1\r\nid2\r\n'
      ;;
  esac
fi
exit 0
`)

	ids, err := getContainerIDs("/home/user/project")
	if err != nil {
		t.Fatalf("getContainerIDs returned error: %v", err)
	}
	if ids != "id1 id2" {
		t.Fatalf("expected normalized ids 'id1 id2', got %q", ids)
	}
}

func TestGetAllRelatedContainers_IncludesComposeContainersWithoutDuplicates(t *testing.T) {
	bin := setupPath(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeExecutable(t, bin, "docker", `#!/bin/sh
case "$*" in
  'ps -a -q --filter label=devcontainer.local_folder=/tmp/project')
    printf 'main1\nmain2\n'
    ;;
  'inspect -f {{ if .Config.Labels }}{{ index .Config.Labels "com.docker.compose.project" }}{{ end }} main1')
    printf 'projA\n'
    ;;
  'inspect -f {{ if .Config.Labels }}{{ index .Config.Labels "com.docker.compose.project" }}{{ end }} main2')
    printf '<no value>\n'
    ;;
  'ps -a -q --filter label=com.docker.compose.project=projA')
    printf 'main1\nside1\n'
    ;;
esac
exit 0
`)

	ids, err := getAllRelatedContainers("/tmp/project")
	if err != nil {
		t.Fatalf("getAllRelatedContainers returned error: %v", err)
	}

	sort.Strings(ids)
	want := []string{"main1", "main2", "side1"}
	if strings.Join(ids, ",") != strings.Join(want, ",") {
		t.Fatalf("expected ids %v, got %v", want, ids)
	}
}
