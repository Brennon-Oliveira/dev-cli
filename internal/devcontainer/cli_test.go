package devcontainer

import (
	"errors"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
)

func TestDevContainerCLIImpl_Up(t *testing.T) {
	mock := exec.NewMockExecutor()
	cli := NewDevContainerCLI(mock)

	err := cli.Up("/workspace/project")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mock.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mock.Calls))
	}

	want := "devcontainer up --workspace-folder /workspace/project"
	if mock.Calls[0] != want {
		t.Errorf("expected %q, got %q", want, mock.Calls[0])
	}
}

func TestDevContainerCLIImpl_UpWithError(t *testing.T) {
	mock := exec.NewMockExecutor()
	mock.CombinedOutErr = errors.New("build failed")
	cli := NewDevContainerCLI(mock)

	err := cli.Up("/workspace/project")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDevContainerCLIImpl_Exec(t *testing.T) {
	mock := exec.NewMockExecutor()
	cli := NewDevContainerCLI(mock)

	err := cli.Exec("/workspace/project", []string{"npm", "test"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mock.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mock.Calls))
	}

	want := "devcontainer exec --workspace-folder /workspace/project npm test"
	if mock.Calls[0] != want {
		t.Errorf("expected %q, got %q", want, mock.Calls[0])
	}
}

func TestDevContainerCLIImpl_ReadConfiguration_Success(t *testing.T) {
	mock := exec.NewMockExecutor()
	mock.OutputResult = `{"workspace":{"workspaceFolder":"/workspace/custom"}}`
	cli := NewDevContainerCLI(mock)

	config, err := cli.ReadConfiguration("/workspace/project")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.WorkspaceFolder != "/workspace/custom" {
		t.Errorf("expected /workspace/custom, got %q", config.WorkspaceFolder)
	}
}

func TestDevContainerCLIImpl_ReadConfiguration_Fallback(t *testing.T) {
	mock := exec.NewMockExecutor()
	mock.OutputErr = errors.New("command not found")
	cli := NewDevContainerCLI(mock)

	config, err := cli.ReadConfiguration("/workspace/project")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.WorkspaceFolder != "/workspaces" {
		t.Errorf("expected /workspaces fallback, got %q", config.WorkspaceFolder)
	}
}

func TestDevContainerCLIImpl_ReadConfiguration_EmptyJSON(t *testing.T) {
	mock := exec.NewMockExecutor()
	mock.OutputResult = `{}`
	cli := NewDevContainerCLI(mock)

	config, err := cli.ReadConfiguration("/workspace/project")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.WorkspaceFolder != "/workspaces" {
		t.Errorf("expected /workspaces fallback, got %q", config.WorkspaceFolder)
	}
}

func TestMockDevContainerCLI(t *testing.T) {
	mock := NewMockDevContainerCLI()

	_ = mock.Up("/test")
	_ = mock.Exec("/test", []string{"ls"})
	_, _ = mock.ReadConfiguration("/test")

	if len(mock.Calls) != 3 {
		t.Errorf("expected 3 calls, got %d", len(mock.Calls))
	}
}

func TestArgsToString(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{[]string{"docker", "ps"}, "docker ps"},
		{[]string{"echo"}, "echo"},
		{[]string{"git", "commit", "-m", "message"}, "git commit -m message"},
	}

	for _, tt := range tests {
		got := argsToString(tt.args)
		if got != tt.expected {
			t.Errorf("argsToString(%v) = %q, want %q", tt.args, got, tt.expected)
		}
	}
}
