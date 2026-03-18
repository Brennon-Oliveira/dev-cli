package completion

import (
	"github.com/spf13/cobra"
)

type Installer interface {
	Install(shell string) error
	DetectShell() string
}

type CompletionInstaller struct {
	rootCmd *cobra.Command
}

func NewCompletionInstaller(rootCmd *cobra.Command) *CompletionInstaller {
	return &CompletionInstaller{rootCmd: rootCmd}
}

type MockInstaller struct {
	InstallErr    error
	DetectResult  string
	InstallCalled bool
}

func NewMockInstaller() *MockInstaller {
	return &MockInstaller{DetectResult: "bash"}
}

func (m *MockInstaller) Install(shell string) error {
	m.InstallCalled = true
	return m.InstallErr
}

func (m *MockInstaller) DetectShell() string {
	return m.DetectResult
}
