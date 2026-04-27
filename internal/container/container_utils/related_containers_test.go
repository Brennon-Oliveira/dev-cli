package container_utils

import (
	"fmt"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/exec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPather struct {
	mock.Mock
}

func (m *mockPather) GetAbsPath(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

func (m *mockPather) GetPathFromArgs(args []string) string {
	callArgs := m.Called(args)
	return callArgs.String(0)
}

func (m *mockPather) GetRealPath(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

func TestResolvePaths_SinglePathNoRealPathDifference(t *testing.T) {
	r := require.New(t)

	mockPather := new(mockPather)
	mockPather.On("GetRealPath", "/home/user/app").Return("/home/user/app", nil)

	paths := ResolvePaths("/home/user/app", mockPather)

	r.Len(paths, 1)
	assert.Equal(t, "/home/user/app", paths[0])
	mockPather.AssertExpectations(t)
}

func TestResolvePaths_SinglePathWithRealPathDifference(t *testing.T) {
	r := require.New(t)

	mockPather := new(mockPather)
	mockPather.On("GetRealPath", "/home/user/app").Return("/mnt/wsl/app", nil)

	paths := ResolvePaths("/home/user/app", mockPather)

	r.Len(paths, 2)
	assert.Equal(t, "/home/user/app", paths[0])
	assert.Equal(t, "/mnt/wsl/app", paths[1])
	mockPather.AssertExpectations(t)
}

func TestResolvePaths_GetRealPathReturnsError_IgnoresError(t *testing.T) {
	r := require.New(t)

	mockPather := new(mockPather)
	mockPather.On("GetRealPath", "/home/user/app").Return("/home/user/app", fmt.Errorf("some error"))

	paths := ResolvePaths("/home/user/app", mockPather)

	r.Len(paths, 1)
	assert.Equal(t, "/home/user/app", paths[0])
	mockPather.AssertExpectations(t)
}

func TestFindMainContainersForPath_ContainersFound(t *testing.T) {
	r := require.New(t)

	mockOutput := "abc123\ndef456"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	ids, err := FindMainContainersForPath("docker", "/home/user/app", executor)

	r.Nil(err)
	r.Len(ids, 2)
	assert.Equal(t, "abc123", ids[0])
	assert.Equal(t, "def456", ids[1])
	executor.AssertExpectations(t)
}

func TestFindMainContainersForPath_NoContainersFound(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(""), nil)

	ids, err := FindMainContainersForPath("docker", "/home/user/app", executor)

	r.Nil(err)
	r.Nil(ids)
	executor.AssertExpectations(t)
}

func TestFindMainContainersForPath_ExecutorReturnsError(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return(nil, fmt.Errorf("docker error"))

	ids, err := FindMainContainersForPath("docker", "/home/user/app", executor)

	r.NotNil(err)
	assert.ErrorContains(t, err, "docker error")
	r.Nil(ids)
	executor.AssertExpectations(t)
}

func TestFindMainContainersForPath_OutputWithWindowsLineEndings(t *testing.T) {
	r := require.New(t)

	mockOutput := "abc123\r\ndef456\r\n"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	ids, err := FindMainContainersForPath("docker", "/home/user/app", executor)

	r.Nil(err)
	r.Len(ids, 2)
	assert.Equal(t, "abc123", ids[0])
	assert.Equal(t, "def456", ids[1])
}

func TestFindMainContainersForPath_UsesPodmanTool(t *testing.T) {
	r := require.New(t)

	mockOutput := "container123"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "podman"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	ids, err := FindMainContainersForPath("podman", "/path", executor)

	r.Nil(err)
	r.Len(ids, 1)
	executor.AssertExpectations(t)
}

func TestExtractProjectFromContainer_ProjectFound(t *testing.T) {
	r := require.New(t)

	mockOutput := "my-app-compose"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	project, err := ExtractProjectFromContainer("docker", "abc123", executor)

	r.Nil(err)
	assert.Equal(t, "my-app-compose", project)
	executor.AssertExpectations(t)
}

func TestExtractProjectFromContainer_NoProjectLabel(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(""), nil)

	project, err := ExtractProjectFromContainer("docker", "abc123", executor)

	r.Nil(err)
	assert.Equal(t, "", project)
	executor.AssertExpectations(t)
}

func TestExtractProjectFromContainer_NoValueReturned(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte("<no value>"), nil)

	project, err := ExtractProjectFromContainer("docker", "abc123", executor)

	r.Nil(err)
	assert.Equal(t, "", project)
	executor.AssertExpectations(t)
}

func TestExtractProjectFromContainer_ExecutorReturnsError(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return(nil, fmt.Errorf("inspect error"))

	project, err := ExtractProjectFromContainer("docker", "abc123", executor)

	r.NotNil(err)
	assert.ErrorContains(t, err, "inspect error")
	assert.Equal(t, "", project)
	executor.AssertExpectations(t)
}

func TestExtractProjectFromContainer_ProjectNameWithSpecialCharacters(t *testing.T) {
	r := require.New(t)

	projectName := "my-app_v2.0-beta"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(projectName), nil)

	project, err := ExtractProjectFromContainer("docker", "abc123", executor)

	r.Nil(err)
	assert.Equal(t, projectName, project)
}

func TestFindComposeContainersForProject_ContainersFound(t *testing.T) {
	r := require.New(t)

	mockOutput := "id1\nid2\nid3"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	ids, err := FindComposeContainersForProject("docker", "myapp", executor)

	r.Nil(err)
	r.Len(ids, 3)
	assert.Equal(t, "id1", ids[0])
	assert.Equal(t, "id2", ids[1])
	assert.Equal(t, "id3", ids[2])
	executor.AssertExpectations(t)
}

func TestFindComposeContainersForProject_NoContainersFound(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(""), nil)

	ids, err := FindComposeContainersForProject("docker", "myapp", executor)

	r.Nil(err)
	r.Nil(ids)
	executor.AssertExpectations(t)
}

func TestFindComposeContainersForProject_ExecutorReturnsError(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return(nil, fmt.Errorf("docker error"))

	ids, err := FindComposeContainersForProject("docker", "myapp", executor)

	r.NotNil(err)
	assert.ErrorContains(t, err, "docker error")
	r.Nil(ids)
	executor.AssertExpectations(t)
}

func TestFindComposeContainersForProject_OutputWithWindowsLineEndings(t *testing.T) {
	r := require.New(t)

	mockOutput := "id1\r\nid2\r\nid3\r\n"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "docker"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	ids, err := FindComposeContainersForProject("docker", "myapp", executor)

	r.Nil(err)
	r.Len(ids, 3)
	assert.Equal(t, "id1", ids[0])
	assert.Equal(t, "id2", ids[1])
	assert.Equal(t, "id3", ids[2])
}

func TestFindComposeContainersForProject_UsesPodmanTool(t *testing.T) {
	r := require.New(t)

	mockOutput := "container1\ncontainer2"
	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.MatchedBy(func(name string) bool {
		return name == "podman"
	}), mock.Anything).Return([]byte(mockOutput), nil)

	ids, err := FindComposeContainersForProject("podman", "app", executor)

	r.Nil(err)
	r.Len(ids, 2)
	executor.AssertExpectations(t)
}

func TestDeduplicateAndFilterContainerIDs_SingleID(t *testing.T) {
	idMap := map[string]bool{
		"abc123": true,
	}

	result := DeduplicateAndFilterContainerIDs(idMap)

	assert.Len(t, result, 1)
	assert.Contains(t, result, "abc123")
}

func TestDeduplicateAndFilterContainerIDs_MultipleUniqueIDs(t *testing.T) {
	idMap := map[string]bool{
		"abc123": true,
		"def456": true,
		"ghi789": true,
	}

	result := DeduplicateAndFilterContainerIDs(idMap)

	assert.Len(t, result, 3)
	assert.Contains(t, result, "abc123")
	assert.Contains(t, result, "def456")
	assert.Contains(t, result, "ghi789")
}

func TestDeduplicateAndFilterContainerIDs_FilterEmptyID(t *testing.T) {
	idMap := map[string]bool{
		"abc123": true,
		"":       true,
		"def456": true,
	}

	result := DeduplicateAndFilterContainerIDs(idMap)

	assert.Len(t, result, 2)
	assert.Contains(t, result, "abc123")
	assert.Contains(t, result, "def456")
	assert.NotContains(t, result, "")
}

func TestDeduplicateAndFilterContainerIDs_EmptyMap(t *testing.T) {
	idMap := map[string]bool{}

	result := DeduplicateAndFilterContainerIDs(idMap)

	assert.Len(t, result, 0)
}

func TestDeduplicateAndFilterContainerIDs_OnlyEmptyID(t *testing.T) {
	idMap := map[string]bool{
		"": true,
	}

	result := DeduplicateAndFilterContainerIDs(idMap)

	assert.Len(t, result, 0)
}

func TestDeduplicateAndFilterContainerIDs_ManyIDsWithEmpty(t *testing.T) {
	idMap := map[string]bool{
		"id1": true,
		"":    true,
		"id2": true,
		"id3": true,
	}

	result := DeduplicateAndFilterContainerIDs(idMap)

	assert.Len(t, result, 3)
	assert.Contains(t, result, "id1")
	assert.Contains(t, result, "id2")
	assert.Contains(t, result, "id3")
}
