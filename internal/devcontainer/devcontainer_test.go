package devcontainer

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestUp_JustReturnWhenWork(t *testing.T) {
	workspace := "/tmp/workspace"
	executor := exec.NewMockExecutor(t)

	var capturedArgs []string
	executor.EXPECT().Run(mock.Anything, mock.Anything).Run(func(name string, args ...string) {
		capturedArgs = args
	}).Return(nil)

	devcontainerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := devcontainerCLI.Up(workspace)

	assert.Nil(t, err)
	assert.Contains(t, capturedArgs, workspace)
}

func TestUp_IfErrorThrowUp(t *testing.T) {
	workspace := "/tmp/workspace"
	executor := exec.NewMockExecutor(t)

	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(fmt.Errorf("generic error"))

	devcontainerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := devcontainerCLI.Up(workspace)

	assert.ErrorContains(t, err, "generic error")
}

func TestReadConfiguration_ReturnUnmarshaledConfig(t *testing.T) {
	r := require.New(t)
	workspace := "/tmp/workspace"

	configMock := &DevContainerConfiguration{
		Workspace: DevContainerConfiguration_Workspace{
			WorkspaceFolder: "my-project",
		},
	}

	configJsonRaw, err := json.Marshal(configMock)
	r.Nil(err)

	executor := exec.NewMockExecutor(t)

	executor.EXPECT().Output(mock.Anything, mock.Anything).RunAndReturn(func(name string, args ...string) ([]byte, error) {
		if assert.Contains(t, args, workspace) {
			return configJsonRaw, nil
		}
		return []byte{}, fmt.Errorf("workspace dir not passed")
	})

	devcontainerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	config, err := devcontainerCLI.ReadConfiguration(workspace)
	r.Nil(err)

	assert.IsType(t, DevContainerConfiguration{}, *config)
	assert.IsType(t, DevContainerConfiguration_Workspace{}, config.Workspace)
	assert.Equal(t, config.Workspace.WorkspaceFolder, "my-project")
}

func TestReadConfiguration_ThrowErrorIfNotAbleToUnmarshal(t *testing.T) {
	r := require.New(t)
	workspace := "/tmp/workspace"

	configMock := &DevContainerConfiguration{
		Workspace: DevContainerConfiguration_Workspace{
			WorkspaceFolder: "my-project",
		},
	}

	configJsonRaw, err := json.Marshal(configMock)
	r.Nil(err)

	executor := exec.NewMockExecutor(t)

	executor.EXPECT().Output(mock.Anything, mock.Anything).Return(append(configJsonRaw, []byte("test-for-break")...), nil)

	devcontainerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	config, err := devcontainerCLI.ReadConfiguration(workspace)
	assert.Empty(t, config)
	var syntaxErr *json.SyntaxError
	assert.ErrorAs(t, err, &syntaxErr)
}
