package completion

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCompletionInstaller(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	installer := NewCompletionInstaller(rootCmd)

	if installer == nil {
		t.Error("expected installer to be created")
	}
}

func TestMockInstaller(t *testing.T) {
	mock := NewMockInstaller()

	if mock.DetectShell() != "bash" {
		t.Errorf("expected bash, got %s", mock.DetectShell())
	}

	err := mock.Install("zsh")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !mock.InstallCalled {
		t.Error("expected Install to be called")
	}
}

func TestMockInstaller_WithError(t *testing.T) {
	mock := NewMockInstaller()
	mock.InstallErr = mockErr("install failed")

	err := mock.Install("bash")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDetectShell(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	installer := NewCompletionInstaller(rootCmd)

	shell := installer.DetectShell()
	if shell != "" && shell != "bash" && shell != "zsh" && shell != "powershell" {
		t.Errorf("unexpected shell: %s", shell)
	}
}

type mockErr string

func (e mockErr) Error() string { return string(e) }
