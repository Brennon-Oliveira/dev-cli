package exec

import (
	"errors"
	"testing"
)

func TestMockExecutor_Run(t *testing.T) {
	m := NewMockExecutor()

	err := m.Run("docker", "ps")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(m.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(m.Calls))
	}
	if m.Calls[0] != "docker ps" {
		t.Errorf("expected 'docker ps', got %q", m.Calls[0])
	}
}

func TestMockExecutor_RunWithError(t *testing.T) {
	m := NewMockExecutor()
	m.RunErr = errors.New("command failed")

	err := m.Run("docker", "ps")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestMockExecutor_RunInteractive(t *testing.T) {
	m := NewMockExecutor()

	err := m.RunInteractive("docker", "exec", "-it", "container", "bash")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(m.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(m.Calls))
	}
	if m.Calls[0] != "docker exec -it container bash" {
		t.Errorf("unexpected call: %q", m.Calls[0])
	}
}

func TestMockExecutor_RunDetached(t *testing.T) {
	m := NewMockExecutor()

	err := m.RunDetached("code", "--folder-uri", "vscode-remote://...")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(m.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(m.Calls))
	}
}

func TestMockExecutor_Output(t *testing.T) {
	m := NewMockExecutor()
	m.OutputResult = "container-id-123"

	out, err := m.Output("docker", "ps", "-q")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if out != "container-id-123" {
		t.Errorf("expected 'container-id-123', got %q", out)
	}
}

func TestMockExecutor_OutputWithError(t *testing.T) {
	m := NewMockExecutor()
	m.OutputErr = errors.New("no such container")

	_, err := m.Output("docker", "ps", "-q")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestMockExecutor_Reset(t *testing.T) {
	m := NewMockExecutor()
	m.Run("docker", "ps")
	m.OutputResult = "result"
	m.OutputErr = errors.New("err")

	m.Reset()

	if len(m.Calls) != 0 {
		t.Errorf("expected no calls after reset, got %d", len(m.Calls))
	}
	if m.OutputResult != "" {
		t.Errorf("expected empty output result, got %q", m.OutputResult)
	}
	if m.OutputErr != nil {
		t.Errorf("expected nil error, got %v", m.OutputErr)
	}
}

func TestMockExecutor_MultipleCalls(t *testing.T) {
	m := NewMockExecutor()

	_ = m.Run("docker", "ps")
	_ = m.Run("docker", "images")
	_, _ = m.Output("docker", "inspect", "id")

	if len(m.Calls) != 3 {
		t.Fatalf("expected 3 calls, got %d", len(m.Calls))
	}

	expected := []string{"docker ps", "docker images", "docker inspect id"}
	for i, call := range expected {
		if m.Calls[i] != call {
			t.Errorf("call %d: expected %q, got %q", i, call, m.Calls[i])
		}
	}
}

func TestFormatCall(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"docker", []string{"ps"}, "docker ps"},
		{"docker", []string{}, "docker"},
		{"echo", []string{"hello", "world"}, "echo hello world"},
	}

	for _, tt := range tests {
		got := formatCall(tt.name, tt.args...)
		if got != tt.expected {
			t.Errorf("formatCall(%q, %v) = %q, want %q", tt.name, tt.args, got, tt.expected)
		}
	}
}
