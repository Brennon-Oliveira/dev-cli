package vscode

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/devcontainer"
	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/Brennon-Oliveira/dev-cli/internal/pather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetContainerWorkspaceURI_ShouldGetContainerURI(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	pather := pather.NewMockPather(t)
	devcontainerCLI := devcontainer.NewMockDevContainerCLI(t)

	absPath := "/tmp/project"

	pather.EXPECT().GetRealPath("/tmp/project").Return("/meu-wsl?tmp/project", nil)

	devcontainerCLI.EXPECT().GetWorkspaceFolder("/tmp/project").Return("/meu-wsl?tmp/project/workspaces//", nil)

	vscode := NewVSCode(
		WithPather(pather),
		WithDevcontainerCLI(devcontainerCLI),
	)

	got, err := vscode.GetContainerWorkspaceURI(absPath)
	r.Nil(err)

	hexBase := hex.EncodeToString([]byte("/meu-wsl?tmp/project"))
	expected := fmt.Sprintf("vscode-remote://dev-container+%s/meu-wsl?tmp/project/workspaces//", hexBase)

	a.Equal(expected, got)
}

func TestOpenWorkspaceByURI_ShouldTryOpenVscode(t *testing.T) {
	a := assert.New(t)

	executor := exec.NewMockExecutor(t)

	hexBase := hex.EncodeToString([]byte("/meu-wsl?tmp/project"))
	workspaceURI := fmt.Sprintf("vscode-remote://dev-container+%s/meu-wsl?tmp/project/workspaces//", hexBase)

	executor.EXPECT().RunDetached(mock.AnythingOfType("string"), mock.Anything).RunAndReturn(func(name string, args ...string) error {
		if !a.Equal("code", name) {
			return fmt.Errorf("invalid cli for code editor")
		}
		if !a.Contains(args, workspaceURI) {
			return fmt.Errorf("the workspace URI are not used")
		}
		return nil
	})

	vscode := NewVSCode(
		WithExecutor(executor),
	)

	err := vscode.OpenWorkspaceByURI(workspaceURI)
	a.Nil(err)
}
