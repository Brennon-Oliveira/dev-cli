package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func writeExecutable(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		t.Fatalf("failed to write executable %s: %v", name, err)
	}
}

func setupCommandEnv(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	bin := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(bin, 0755); err != nil {
		t.Fatalf("failed to create bin dir: %v", err)
	}
	callLog := filepath.Join(t.TempDir(), "calls.log")

	t.Setenv("HOME", home)
	t.Setenv("DEVCLI_CALL_LOG", callLog)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	writeExecutable(t, bin, "docker", `#!/bin/sh
echo "$0 $*" >> "$DEVCLI_CALL_LOG"
case "$1" in
  ps)
    case "$*" in
      *'label=com.docker.compose.project='*)
        printf 'main1234567890\nside0987654321\n'
        ;;
      *'-a -q --filter label=devcontainer.local_folder='*)
        printf 'main1234567890\n'
        ;;
      *'-q --filter label=devcontainer.local_folder='*)
        printf 'main1234567890\n'
        ;;
    esac
    ;;
  inspect)
    printf 'projx\n'
    ;;
  port)
    printf '8080/tcp -> 0.0.0.0:8080\n'
    ;;
esac
exit 0
`)

	writeExecutable(t, bin, "devcontainer", `#!/bin/sh
echo "$0 $*" >> "$DEVCLI_CALL_LOG"
if [ "$1" = "read-configuration" ]; then
  printf '{"workspace":{"workspaceFolder":"/workspaces"}}'
fi
exit 0
`)

	writeExecutable(t, bin, "code", `#!/bin/sh
echo "$0 $*" >> "$DEVCLI_CALL_LOG"
exit 0
`)

	writeExecutable(t, bin, "powershell", `#!/bin/sh
echo "$0 $*" >> "$DEVCLI_CALL_LOG"
printf '/tmp/devcli-profile.ps1\n'
exit 0
`)

	return callLog
}

func executeRoot(args ...string) error {
	globalFlag = false
	interactiveFlag = false
	follow = false
	execPath = ""
	resetCommandBoolFlags(rootCmd, "help")
	resetCommandBoolFlags(rootCmd, "version")
	rootCmd.SetArgs(args)
	_, err := rootCmd.ExecuteC()
	return err
}

func resetCommandBoolFlags(c *cobra.Command, name string) {
	if f := c.Flags().Lookup(name); f != nil {
		_ = f.Value.Set("false")
	}
	if f := c.PersistentFlags().Lookup(name); f != nil {
		_ = f.Value.Set("false")
	}
	for _, sub := range c.Commands() {
		resetCommandBoolFlags(sub, name)
	}
}

func TestCommands_AllMainCommandsExecuteAtLeastOnce(t *testing.T) {
	callLog := setupCommandEnv(t)

	commands := [][]string{
		{"list"},
		{"clean"},
		{"up", "."},
		{"run", "."},
		{"open", "."},
		{"shell", "."},
		{"exec", "--path", ".", "echo", "ok"},
		{"ports", "."},
		{"logs", "-f", "."},
		{"down", "."},
		{"kill", "."},
		{"config", "core.tool", "docker", "--global"},
		{"config", "core.tool", "--global"},
		{"add-completion", "bash"},
		{"add-completion", "zsh"},
		{"add-completion", "powershell"},
	}

	for _, cmdArgs := range commands {
		if err := executeRoot(cmdArgs...); err != nil {
			t.Fatalf("command %q failed: %v", strings.Join(cmdArgs, " "), err)
		}
	}

	data, err := os.ReadFile(callLog)
	if err != nil {
		t.Fatalf("failed to read call log: %v", err)
	}
	logText := string(data)

	checks := []string{
		"docker ps --filter label=devcontainer.local_folder",
		"devcontainer up --workspace-folder",
		"devcontainer exec --workspace-folder",
		"docker logs -f",
		"docker stop",
		"docker rm -f",
		"code --folder-uri",
		"powershell -NoProfile -Command Write-Host $PROFILE",
	}

	for _, c := range checks {
		if !strings.Contains(logText, c) {
			t.Fatalf("expected call log to contain %q\nlog:\n%s", c, logText)
		}
	}
}

func TestCommands_AllCmdSubcommandsAreRegistered(t *testing.T) {
	want := []string{
		"add-completion",
		"clean",
		"config",
		"down",
		"exec",
		"kill",
		"list",
		"logs",
		"open",
		"ports",
		"run",
		"shell",
		"up",
		"update",
	}

	haveMap := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		haveMap[c.Name()] = true
	}

	var missing []string
	for _, name := range want {
		if !haveMap[name] {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf("missing registered subcommands: %v", missing)
	}
}

func TestCommands_AllCommandHelpEntriesWork(t *testing.T) {
	setupCommandEnv(t)

	commands := [][]string{
		{"--help"},
		{"--version"},
		{"add-completion", "--help"},
		{"clean", "--help"},
		{"config", "--help"},
		{"down", "--help"},
		{"exec", "--help"},
		{"kill", "--help"},
		{"list", "--help"},
		{"logs", "--help"},
		{"open", "--help"},
		{"ports", "--help"},
		{"run", "--help"},
		{"shell", "--help"},
		{"up", "--help"},
		{"update", "--help"},
	}

	for _, args := range commands {
		if err := executeRoot(args...); err != nil {
			t.Fatalf("help/version failed for %q: %v", strings.Join(args, " "), err)
		}
	}
}

func TestCommandEdgeCases_ValidationAndErrors(t *testing.T) {
	setupCommandEnv(t)

	tests := []struct {
		name string
		args []string
	}{
		{name: "config requires global flag", args: []string{"config", "core.tool", "docker"}},
		{name: "config invalid key", args: []string{"config", "bad.key", "docker", "--global"}},
		{name: "config invalid value", args: []string{"config", "core.tool", "bad", "--global"}},
		{name: "config too many args", args: []string{"config", "a", "b", "c"}},
		{name: "down too many args", args: []string{"down", "a", "b"}},
		{name: "exec requires command", args: []string{"exec"}},
		{name: "unsupported completion shell", args: []string{"add-completion", "fish"}},
	}

	for _, tt := range tests {
		err := executeRoot(tt.args...)
		if err == nil {
			t.Fatalf("expected error for %s", tt.name)
		}
	}
}

func TestAddCompletion_AutoDetectsShellWhenNoArgs(t *testing.T) {
	setupCommandEnv(t)
	t.Setenv("SHELL", "/bin/zsh")

	if err := executeRoot("add-completion"); err != nil {
		t.Fatalf("add-completion with detected shell failed: %v", err)
	}
}

func TestConfigCommand_ValidArgsFunction(t *testing.T) {
	keys, directive := configCmd.ValidArgsFunction(configCmd, []string{}, "core")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("unexpected shell directive for keys: %v", directive)
	}
	if len(keys) != 1 || keys[0] != "core.tool" {
		t.Fatalf("expected key suggestion [core.tool], got %v", keys)
	}

	values, directive := configCmd.ValidArgsFunction(configCmd, []string{"core.tool"}, "do")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("unexpected shell directive for values: %v", directive)
	}
	if len(values) != 1 || values[0] != "docker" {
		t.Fatalf("expected value suggestion [docker], got %v", values)
	}

	none, directive := configCmd.ValidArgsFunction(configCmd, []string{"unknown"}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("unexpected shell directive for unknown key: %v", directive)
	}
	if none != nil {
		t.Fatalf("expected nil suggestions for unknown key, got %v", none)
	}
}

func TestUpdateCommand_DownloadAndExtractFlow(t *testing.T) {
	setupCommandEnv(t)
	oldTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = oldTransport
		_ = os.Remove(filepath.Join(os.TempDir(), "dev_new"))
		_ = os.Remove(filepath.Join(os.TempDir(), "dev_new.exe"))
	})

	tarData := buildTarGzWithBinary(t, "dev", []byte("binary-content"))

	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case "https://api.github.com/repos/Brennon-Oliveira/dev-cli/releases/latest":
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"tag_name":"v9.9.9"}`)),
				Header:     make(http.Header),
			}, nil
		case "https://github.com/Brennon-Oliveira/dev-cli/releases/download/v9.9.9/dev-linux-amd64.tar.gz":
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(tarData)),
				Header:     make(http.Header),
			}, nil
		default:
			t.Fatalf("unexpected URL requested: %s", req.URL.String())
			return nil, nil
		}
	})

	if err := executeRoot("update"); err != nil {
		t.Fatalf("update command failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(os.TempDir(), "dev_new")); err != nil {
		t.Fatalf("expected extracted binary at %s: %v", filepath.Join(os.TempDir(), "dev_new"), err)
	}
}

func TestUpdateCommand_NoUpdateWhenAlreadyLatest(t *testing.T) {
	setupCommandEnv(t)
	oldTransport := http.DefaultTransport
	oldVersion := Version
	Version = "v1.2.3"
	t.Cleanup(func() {
		http.DefaultTransport = oldTransport
		Version = oldVersion
	})

	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://api.github.com/repos/Brennon-Oliveira/dev-cli/releases/latest" {
			t.Fatalf("unexpected URL requested: %s", req.URL.String())
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"tag_name":"v1.2.3"}`)),
			Header:     make(http.Header),
		}, nil
	})

	if err := executeRoot("update"); err != nil {
		t.Fatalf("expected no error when already latest, got: %v", err)
	}
}

func buildTarGzWithBinary(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	buf := &bytes.Buffer{}
	gzw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gzw)

	hdr := &tar.Header{Name: name, Mode: 0755, Size: int64(len(content))}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("failed to write tar content: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("failed to close tar writer: %v", err)
	}
	if err := gzw.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}

	return buf.Bytes()
}
