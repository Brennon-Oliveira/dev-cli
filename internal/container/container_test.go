package container

import (
	"fmt"
	"os"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/config"
	"github.com/Brennon-Oliveira/dev-cli/internal/container/container_utils"
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

// ============================================================================
// Helpers
// ============================================================================

type mockConfig struct {
	mock.Mock
}

func (m *mockConfig) Load() config.GlobalConfig {
	args := m.Called()
	return args.Get(0).(config.GlobalConfig)
}

func createMockConfigWithTool(tool string) *mockConfig {
	mockCfg := new(mockConfig)
	globalCfg := config.GlobalConfig{
		Core: struct {
			Tool string `json:"tool"`
		}{Tool: tool},
	}
	mockCfg.On("Load").Return(globalCfg)
	return mockCfg
}

// ============================================================================
// Tests for ListContainersOfActiveDevcontainers
// ============================================================================

func TestListContainersOfActiveDevcontainers_SuccessfulExecutionWithContainers(t *testing.T) {
	r := require.New(t)

	// Setup mock output
	mockOutput := "f7bc76eda682\tsimple-financial-app_devcontainer-sfa.app-1\tUp 5 minutes\t/home/brennon/projects/sfa/simple-financial-app\n" +
		"9c88be57f94f\tsimple-financial-app_devcontainer-sfa.smtp-1\tUp 5 minutes\t\n" +
		"10103216832d\tsimple-financial-app_devcontainer-sfa.db-1\tUp 10 minutes\t"

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.Anything, mock.Anything, mock.Anything).Return([]byte(mockOutput), nil)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.ListContainersOfActiveDevcontainers()

	r.Nil(err)
}

// ============================================================================
// Tests for CleanResources
// ============================================================================

func TestCleanResources_BothPrunesSucceed_ReturnsNil(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(nil).Times(2)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestCleanResources_ContainerPruneFailsReturnsNilAndStops(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	// Container prune fails, should return nil and NOT call network prune
	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(fmt.Errorf("container prune error")).Once()

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	// Error in container prune returns nil (per the implementation)
	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestCleanResources_ContainerPruneFailsDoesNotContinue(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	// Container prune fails, should return nil and NOT call network prune
	// So we only expect 1 call
	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(fmt.Errorf("container error")).Once()

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	// Returns nil when container prune fails
	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestCleanResources_NetworkPruneFailsAfterContainerPruneSucceeds(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	// First call (container prune) succeeds, second call (network prune) fails
	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(nil).Once()
	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(fmt.Errorf("network prune failed")).Once()

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	r.NotNil(err)
	assert.ErrorContains(t, err, "network prune failed")
	executor.AssertExpectations(t)
}

func TestCleanResources_NetworkPruneFailsWithSpecificError(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(nil).Once()
	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(fmt.Errorf("network not found")).Once()

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	r.NotNil(err)
	assert.ErrorContains(t, err, "network not found")
}

func TestCleanResources_CallsExecutorWithCorrectArgs_ContainerPrune(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(nil).Times(2)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	r.Nil(err)
	executor.AssertCalled(t, "Run", "docker", []string{"container", "prune", "-f"})
}

func TestCleanResources_CallsExecutorWithCorrectArgs_NetworkPrune(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(nil).Times(2)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	r.Nil(err)
	executor.AssertCalled(t, "Run", "docker", []string{"network", "prune", "-f"})
}

func TestCleanResources_UsesConfigToolCorrectly(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	// Both calls should use "podman" as the tool (from config)
	executor.EXPECT().Run("podman", mock.Anything).Return(nil).Times(2)

	configMock := createMockConfigWithTool("podman")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	r.Nil(err)
	executor.AssertExpectations(t)
}

func TestCleanResources_NetworkPruneReturnsErrorWhenContainerSucceeds(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(nil).Once()
	executor.EXPECT().Run(mock.Anything, mock.Anything).Return(fmt.Errorf("network error")).Once()

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	// Network error is returned
	r.NotNil(err)
	assert.ErrorContains(t, err, "network error")
	executor.AssertExpectations(t)
}

func TestCleanResources_SequenceOfOperations_BothSucceed(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	var callSequence []string

	executor.EXPECT().Run(mock.Anything, mock.Anything).Run(func(name string, args ...string) {
		if len(args) > 0 && args[0] == "container" {
			callSequence = append(callSequence, "container")
		}
	}).Return(nil).Once()

	executor.EXPECT().Run(mock.Anything, mock.Anything).Run(func(name string, args ...string) {
		if len(args) > 0 && args[0] == "network" {
			callSequence = append(callSequence, "network")
		}
	}).Return(nil).Once()

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.CleanResources()

	r.Nil(err)
	// Verify sequence: container first, then network
	assert.Equal(t, []string{"container", "network"}, callSequence)
	executor.AssertExpectations(t)
}

func TestListContainersOfActiveDevcontainers_ExecutorReturnsError(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("docker not running"))

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.ListContainersOfActiveDevcontainers()

	r.NotNil(err)
	assert.ErrorContains(t, err, "docker not running")
}

func TestListContainersOfActiveDevcontainers_EmptyOutput(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.Anything, mock.Anything, mock.Anything).Return([]byte(""), nil)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.ListContainersOfActiveDevcontainers()

	r.Nil(err)
}

func TestListContainersOfActiveDevcontainers_UsesInjectedParseFunction(t *testing.T) {
	r := require.New(t)

	mockOutput := "id1\tname1\tUp\t/path"
	mockParseWasCalled := false

	customParseFunc := func(output string) map[string][]*container_utils.Container {
		mockParseWasCalled = true
		return container_utils.ParseContainerOutput(output)
	}

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.Anything, mock.Anything, mock.Anything).Return([]byte(mockOutput), nil)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
		WithParseContainerOutput(customParseFunc),
	)

	err := containerCLI.ListContainersOfActiveDevcontainers()

	r.Nil(err)
	r.True(mockParseWasCalled, "custom parse function was not called")
}

func TestListContainersOfActiveDevcontainers_UsesInjectedFormatFunction(t *testing.T) {
	r := require.New(t)

	mockOutput := "id1\tname1\tUp\t/path"
	mockFormatWasCalled := false

	customFormatFunc := func(grouped map[string][]*container_utils.Container) string {
		mockFormatWasCalled = true
		return container_utils.FormatGroupedContainers(grouped)
	}

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.Anything, mock.Anything, mock.Anything).Return([]byte(mockOutput), nil)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
		WithFormatGroupedContainers(customFormatFunc),
	)

	err := containerCLI.ListContainersOfActiveDevcontainers()

	r.Nil(err)
	r.True(mockFormatWasCalled, "custom format function was not called")
}

// ============================================================================
// Tests for Container Type
// ============================================================================

func TestContainerStructure_AllFieldsExportedAndAccessible(t *testing.T) {
	r := require.New(t)

	container := &container_utils.Container{
		ID:          "abc123",
		Names:       "myapp_devcontainer-sfa.app-1",
		Status:      "Up 5 minutes",
		LocalFolder: "/home/user/myapp",
	}

	r.Equal("abc123", container.ID)
	r.Equal("myapp_devcontainer-sfa.app-1", container.Names)
	r.Equal("Up 5 minutes", container.Status)
	r.Equal("/home/user/myapp", container.LocalFolder)
}

func TestContainerStructure_CanCreateAndModifyFields(t *testing.T) {
	r := require.New(t)

	container := &container_utils.Container{}

	container.ID = "test123"
	container.Names = "test_devcontainer-sfa.test-1"
	container.Status = "Up"
	container.LocalFolder = "/test/path"

	r.Equal("test123", container.ID)
	r.Equal("test_devcontainer-sfa.test-1", container.Names)
	r.Equal("Up", container.Status)
	r.Equal("/test/path", container.LocalFolder)
}

// ============================================================================
// Tests for Builder Pattern
// ============================================================================

func TestNewContainerCLI_DefaultFunctionsAreSet(t *testing.T) {
	r := require.New(t)

	containerCLI := NewContainerCLI()

	r.NotNil(containerCLI)
	r.NotNil(containerCLI.parseContainerOutput)
	r.NotNil(containerCLI.formatGroupedContainers)
}

func TestNewContainerCLI_WithExecutor_OptionApplied(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)

	containerCLI := NewContainerCLI(WithExecutor(executor))

	r.NotNil(containerCLI)
	r.Equal(executor, containerCLI.executor)
}

func TestNewContainerCLI_WithConfig_OptionApplied(t *testing.T) {
	r := require.New(t)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(WithConfig(configMock))

	r.NotNil(containerCLI)
	r.Equal(configMock, containerCLI.config)
}

func TestNewContainerCLI_WithParseContainerOutput_OptionApplied(t *testing.T) {
	r := require.New(t)

	customParseFunc := func(output string) map[string][]*container_utils.Container {
		return make(map[string][]*container_utils.Container)
	}

	containerCLI := NewContainerCLI(WithParseContainerOutput(customParseFunc))

	r.NotNil(containerCLI)
	r.NotNil(containerCLI.parseContainerOutput)
}

func TestNewContainerCLI_WithFormatGroupedContainers_OptionApplied(t *testing.T) {
	r := require.New(t)

	customFormatFunc := func(grouped map[string][]*container_utils.Container) string {
		return "custom format"
	}

	containerCLI := NewContainerCLI(WithFormatGroupedContainers(customFormatFunc))

	r.NotNil(containerCLI)
	r.NotNil(containerCLI.formatGroupedContainers)
}

func TestNewContainerCLI_MultipleOptions_AllApplied(t *testing.T) {
	r := require.New(t)

	executor := exec.NewMockExecutor(t)
	configMock := createMockConfigWithTool("docker")

	customParseFunc := func(output string) map[string][]*container_utils.Container {
		return make(map[string][]*container_utils.Container)
	}

	customFormatFunc := func(grouped map[string][]*container_utils.Container) string {
		return "custom"
	}

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
		WithParseContainerOutput(customParseFunc),
		WithFormatGroupedContainers(customFormatFunc),
	)

	r.NotNil(containerCLI)
	r.Equal(executor, containerCLI.executor)
	r.Equal(configMock, containerCLI.config)
	r.NotNil(containerCLI.parseContainerOutput)
	r.NotNil(containerCLI.formatGroupedContainers)
}

// ============================================================================
// Tests for Integration Flow
// ============================================================================

func TestListContainersIntegration_FlowWithRealParsing(t *testing.T) {
	r := require.New(t)

	// Test data simulating real docker output
	mockOutput := "f7bc76eda682\tapp_devcontainer-sfa.app-1\tUp 5 minutes\t/home/user/app\n" +
		"9c88be57f94f\tapp_devcontainer-sfa.db-1\tUp 5 minutes\t"

	executor := exec.NewMockExecutor(t)
	executor.EXPECT().Output(mock.Anything, mock.Anything, mock.Anything).Return([]byte(mockOutput), nil)

	configMock := createMockConfigWithTool("docker")

	containerCLI := NewContainerCLI(
		WithExecutor(executor),
		WithConfig(configMock),
	)

	err := containerCLI.ListContainersOfActiveDevcontainers()

	r.Nil(err)
}
