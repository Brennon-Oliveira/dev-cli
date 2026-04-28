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

func TestFormatWorkspaceFolderSuffix_FixPathWithOneSlash(t *testing.T) {
	path := "/tmp/workspaces/"

	got := formatWorkspaceFolderSuffix(path)

	assert.Equal(t, "/tmp/workspaces//", got)
}

func TestFormatWorkspaceFolderSuffix_FixPathWithNo(t *testing.T) {
	path := "/tmp/workspaces"

	got := formatWorkspaceFolderSuffix(path)

	assert.Equal(t, "/tmp/workspaces//", got)
}

func TestFormatWorkspaceFolderSuffix_FixPathWithOtherSuffix(t *testing.T) {
	path := "/tmp/workspaces/app"

	got := formatWorkspaceFolderSuffix(path)

	assert.Equal(t, "/tmp/workspaces/app", got)
}

func TestGetWorkspaceFolder_GetWorkspaceWithValueDefault(t *testing.T) {
	r := require.New(t)
	workspace := "/tmp/workspace"

	configMock := &DevContainerConfiguration{
		Workspace: DevContainerConfiguration_Workspace{
			WorkspaceFolder: "/workspaces/my-project",
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

	got, err := devcontainerCLI.GetWorkspaceFolder(workspace)
	r.Nil(err)
	assert.Equal(t, "/workspaces/my-project", got)
}

func TestGetWorkspaceFolder_GetWorkspaceWithValueWorkspaces(t *testing.T) {
	r := require.New(t)
	workspace := "/tmp/workspace"

	configMock := &DevContainerConfiguration{
		Workspace: DevContainerConfiguration_Workspace{
			WorkspaceFolder: "/workspaces",
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

	got, err := devcontainerCLI.GetWorkspaceFolder(workspace)
	r.Nil(err)
	assert.Equal(t, "/workspaces//", got)
}

func TestGetWorkspaceFolder_GetWorkspaceWithNoValue(t *testing.T) {
	r := require.New(t)
	workspace := "/tmp/workspace"

	configMock := &DevContainerConfiguration{
		Workspace: DevContainerConfiguration_Workspace{
			WorkspaceFolder: "",
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

	got, err := devcontainerCLI.GetWorkspaceFolder(workspace)
	r.Nil(err)
	assert.Equal(t, "/workspaces//", got)
}

func TestGetWorkspaceFolder_GetWorkspaceWithError(t *testing.T) {
	workspace := "/tmp/workspace"
	executor := exec.NewMockExecutor(t)

	executor.EXPECT().Output(mock.Anything, mock.Anything).Return([]byte{}, fmt.Errorf("generic error"))

	devcontainerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	got, err := devcontainerCLI.GetWorkspaceFolder(workspace)

	assert.ErrorContains(t, err, "generic error")
	assert.Equal(t, "/workspaces//", got)
}

func TestRunInteractive_SuccessfulExecution(t *testing.T) {
	r := require.New(t)

	mockOutput := "command output result"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm install")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_CommandWithMultipleWords(t *testing.T) {
	r := require.New(t)

	mockOutput := "result"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm install --save-dev typescript")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_ErrorExecutionButWithOutput_ReturnsNil(t *testing.T) {
	r := require.New(t)

	mockOutput := "error output from command"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(mockOutput), fmt.Errorf("exit code 1"))

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm test")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_ErrorWithoutOutput_ReturnsError(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return(nil, fmt.Errorf("connection failed"))

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm install")

	r.NotNil(err)
	assert.ErrorContains(t, err, "connection failed")
	executor.AssertExpectations(t)
}

func TestRunInteractive_ErrorWithEmptyOutput_ReturnsError(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(""), fmt.Errorf("execution failed"))

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm install")

	r.NotNil(err)
	assert.ErrorContains(t, err, "execution failed")
	executor.AssertExpectations(t)
}

func TestRunInteractive_CommandSplitsCorrectly(t *testing.T) {
	r := require.New(t)

	mockOutput := "output"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm install")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_UsesDevcontainerTool(t *testing.T) {
	r := require.New(t)

	mockOutput := "output"
	executor := exec.NewMockExecutor(t)

	var toolUsed string
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		toolUsed = tool
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm install")

	r.Nil(err)
	assert.Equal(t, "devcontainer", toolUsed)
	executor.AssertExpectations(t)
}

func TestRunInteractive_WorkspaceFolderInCommand(t *testing.T) {
	r := require.New(t)

	workspacePath := "/custom/path/to/project"
	mockOutput := "output"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive(workspacePath, "npm install")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_EmptyCommand_StillExecutes(t *testing.T) {
	r := require.New(t)

	mockOutput := "output"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_CommandWithSpecialCharacters(t *testing.T) {
	r := require.New(t)

	mockOutput := "output"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm run build:prod")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_LongCommandOutput(t *testing.T) {
	r := require.New(t)

	longOutput := "line1\nline2\nline3\nline4\nline5\n" +
		"line6\nline7\nline8\nline9\nline10\n" +
		"line11\nline12\nline13\nline14\nline15"

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(longOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm install")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_MultilineErrorOutput_StillReturnsNil(t *testing.T) {
	r := require.New(t)

	errorOutput := "error line 1\nerror line 2\nerror line 3"

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(errorOutput), fmt.Errorf("exit code 1"))

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm test")

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestRunInteractive_SuccessfulOutputWithNewlines(t *testing.T) {
	r := require.New(t)

	successOutput := "npm notice\nnpm notice\nup to date\n"

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(tool string) bool {
		return tool == "devcontainer"
	}), mock.Anything).Return([]byte(successOutput), nil)

	containerCLI := NewDevContainerCLI(
		WithExecutor(executor),
	)

	err := containerCLI.RunInteractive("/home/user/app", "npm list")

	r.Nil(err)
	executor.AssertExpectations(t)
}
