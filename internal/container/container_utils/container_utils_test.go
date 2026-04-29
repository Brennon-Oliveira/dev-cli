package container_utils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Brennon-Oliveira/dev-cli/internal/logger"
	"github.com/stretchr/testify/assert"
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

func createMockContainer(id, names, status, localFolder string) *Container {
	return &Container{
		ID:          id,
		Names:       names,
		Status:      status,
		LocalFolder: localFolder,
	}
}

func createMockDockerOutput(containers []*Container) string {
	var lines []string
	lines = append(lines, "CONTAINER ID\tNAMES\tSTATUS\tLOCAL FOLDER")

	for _, c := range containers {
		line := fmt.Sprintf("%s\t%s\t%s\t%s", c.ID, c.Names, c.Status, c.LocalFolder)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func assertGroupContainsContainer(t *testing.T, grouped map[string][]*Container, folder string, expectedNames []string) {
	r := require.New(t)
	containers, exists := grouped[folder]
	r.True(exists, fmt.Sprintf("folder %s not found in grouped containers", folder))

	actualNames := make([]string, len(containers))
	for i, c := range containers {
		actualNames[i] = c.Names
	}

	assert.ElementsMatch(t, expectedNames, actualNames)
}

// ============================================================================
// Tests for ParseContainerOutput
// ============================================================================

func TestParseContainerOutput_SingleContainerWithLocalFolder_GroupedCorrectly(t *testing.T) {
	r := require.New(t)
	containers := []*Container{
		createMockContainer("abc123", "myapp_devcontainer-sfa.app-1", "Up 5 minutes", "/home/user/myapp"),
	}
	output := createMockDockerOutput(containers)

	got := ParseContainerOutput(output)

	r.Len(got, 1)
	assertGroupContainsContainer(t, got, "/home/user/myapp", []string{"myapp_devcontainer-sfa.app-1"})
}

func TestParseContainerOutput_ThreeContainersWithCommonPrefix_AllGroupedTogether(t *testing.T) {
	r := require.New(t)
	containers := []*Container{
		createMockContainer("f7bc76eda682", "simple-financial-app_devcontainer-sfa.app-1", "Up 5 minutes", "/home/brennon/projects/sfa/simple-financial-app"),
		createMockContainer("9c88be57f94f", "simple-financial-app_devcontainer-sfa.smtp-1", "Up 5 minutes", ""),
		createMockContainer("10103216832d", "simple-financial-app_devcontainer-sfa.db-1", "Up 10 minutes", ""),
	}
	output := createMockDockerOutput(containers)

	got := ParseContainerOutput(output)

	r.Len(got, 1)
	r.Len(got["/home/brennon/projects/sfa/simple-financial-app"], 3)
	assertGroupContainsContainer(t, got, "/home/brennon/projects/sfa/simple-financial-app", []string{
		"simple-financial-app_devcontainer-sfa.app-1",
		"simple-financial-app_devcontainer-sfa.smtp-1",
		"simple-financial-app_devcontainer-sfa.db-1",
	})
}

func TestParseContainerOutput_TwoDevcontainersWithDifferentLocalFolders_CreatedMultipleGroups(t *testing.T) {
	r := require.New(t)
	containers := []*Container{
		createMockContainer("id1", "app1_devcontainer-sfa.app-1", "Up 5 minutes", "/home/user/app1"),
		createMockContainer("id2", "app1_devcontainer-sfa.db-1", "Up 5 minutes", ""),
		createMockContainer("id3", "app2_devcontainer-sfa.app-1", "Up 10 minutes", "/home/user/app2"),
		createMockContainer("id4", "app2_devcontainer-sfa.db-1", "Up 10 minutes", ""),
	}
	output := createMockDockerOutput(containers)

	got := ParseContainerOutput(output)

	r.Len(got, 2)
	r.Len(got["/home/user/app1"], 2)
	r.Len(got["/home/user/app2"], 2)
	assertGroupContainsContainer(t, got, "/home/user/app1", []string{
		"app1_devcontainer-sfa.app-1",
		"app1_devcontainer-sfa.db-1",
	})
	assertGroupContainsContainer(t, got, "/home/user/app2", []string{
		"app2_devcontainer-sfa.app-1",
		"app2_devcontainer-sfa.db-1",
	})
}

func TestParseContainerOutput_ContainerWithMultiWordStatus_ParsedCorrectly(t *testing.T) {
	r := require.New(t)
	containers := []*Container{
		createMockContainer("abc123", "myapp_devcontainer-sfa.app-1", "Up About a minute", "/home/user/myapp"),
	}
	output := createMockDockerOutput(containers)

	got := ParseContainerOutput(output)

	r.Len(got, 1)
	parsedContainer := got["/home/user/myapp"][0]
	assert.Equal(t, "Up About a minute", parsedContainer.Status)
}

func TestParseContainerOutput_ContainerWithVeryLongStatus_HandlesProperly(t *testing.T) {
	r := require.New(t)
	longStatus := "Up 29 days, 13 hours and 45 minutes"
	containers := []*Container{
		createMockContainer("abc123", "myapp_devcontainer-sfa.app-1", longStatus, "/home/user/myapp"),
	}
	output := createMockDockerOutput(containers)

	got := ParseContainerOutput(output)

	r.Len(got, 1)
	parsedContainer := got["/home/user/myapp"][0]
	assert.Equal(t, longStatus, parsedContainer.Status)
}

func TestParseContainerOutput_MalformedLineWithLessThanThreeFields_IgnoredSilently(t *testing.T) {
	r := require.New(t)
	output := "CONTAINER ID\tNAMES\tSTATUS\tLOCAL FOLDER\nabc def"

	got := ParseContainerOutput(output)

	r.Len(got, 0)
}

func TestParseContainerOutput_EmptyOutput_ReturnsEmptyMap(t *testing.T) {
	r := require.New(t)
	output := ""

	got := ParseContainerOutput(output)

	r.Len(got, 0)
}

func TestParseContainerOutput_OutputWithOnlyHeader_ReturnsEmptyMap(t *testing.T) {
	r := require.New(t)
	output := "CONTAINER ID\tNAMES\tSTATUS\tLOCAL FOLDER"

	got := ParseContainerOutput(output)

	r.Len(got, 0)
}

func TestParseContainerOutput_ContainerWithoutLocalFolderNoMainContainer_GroupedInDefaultBucket(t *testing.T) {
	r := require.New(t)
	containers := []*Container{
		createMockContainer("orphan123", "orphan_devcontainer-sfa.app-1", "Up 5 minutes", ""),
	}
	output := createMockDockerOutput(containers)

	got := ParseContainerOutput(output)

	r.Len(got, 1)
	_, exists := got["[sem pasta local mapeada]"]
	r.True(exists)
	assert.Len(t, got["[sem pasta local mapeada]"], 1)
}

func TestParseContainerOutput_MultipleOrphanContainers_AllInDefaultBucket(t *testing.T) {
	r := require.New(t)
	containers := []*Container{
		createMockContainer("id1", "orphan1_devcontainer-sfa.app-1", "Up 5 minutes", ""),
		createMockContainer("id2", "orphan1_devcontainer-sfa.db-1", "Up 5 minutes", ""),
		createMockContainer("id3", "orphan2_devcontainer-sfa.app-1", "Up 5 minutes", ""),
	}
	output := createMockDockerOutput(containers)

	got := ParseContainerOutput(output)

	r.Len(got, 1)
	assert.Len(t, got["[sem pasta local mapeada]"], 3)
}

func TestParseContainerOutput_PrefixWithSpecialCharactersAndNumbers_ExtractedCorrectly(t *testing.T) {
	r := require.New(t)
	containers := []*Container{
		createMockContainer("abc123", "my-app-v2.1_devcontainer-sfa.app-1", "Up 5 minutes", "/home/user/my-app-v2.1"),
		createMockContainer("def456", "my-app-v2.1_devcontainer-sfa.db-1", "Up 5 minutes", ""),
	}
	output := createMockDockerOutput(containers)

	got := ParseContainerOutput(output)

	r.Len(got, 1)
	assert.Len(t, got["/home/user/my-app-v2.1"], 2)
	assertGroupContainsContainer(t, got, "/home/user/my-app-v2.1", []string{
		"my-app-v2.1_devcontainer-sfa.app-1",
		"my-app-v2.1_devcontainer-sfa.db-1",
	})
}

// ============================================================================
// Tests for FormatGroupedContainers
// ============================================================================

func TestFormatGroupedContainers_SingleGroupSingleContainer_FormattedCorrectly(t *testing.T) {
	r := require.New(t)
	grouped := map[string][]*Container{
		"/home/user/myapp": {
			createMockContainer("abc123", "myapp_devcontainer-sfa.app-1", "Up 5 minutes", "/home/user/myapp"),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, "/home/user/myapp")
	assert.Contains(t, got, "CONTAINER ID")
	assert.Contains(t, got, "myapp_devcontainer-sfa.app-1")
	assert.Contains(t, got, "---")
}

func TestFormatGroupedContainers_SingleGroupThreeContainers_AllContainersFormatted(t *testing.T) {
	r := require.New(t)
	grouped := map[string][]*Container{
		"/home/user/myapp": {
			createMockContainer("id1", "myapp_devcontainer-sfa.app-1", "Up 5 minutes", "/home/user/myapp"),
			createMockContainer("id2", "myapp_devcontainer-sfa.db-1", "Up 5 minutes", ""),
			createMockContainer("id3", "myapp_devcontainer-sfa.cache-1", "Up 5 minutes", ""),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, "myapp_devcontainer-sfa.app-1")
	assert.Contains(t, got, "myapp_devcontainer-sfa.db-1")
	assert.Contains(t, got, "myapp_devcontainer-sfa.cache-1")
}

func TestFormatGroupedContainers_TwoGroupsWithMultipleContainers_SeparatedProperly(t *testing.T) {
	r := require.New(t)
	grouped := map[string][]*Container{
		"/home/user/app1": {
			createMockContainer("id1", "app1_devcontainer-sfa.app-1", "Up 5 minutes", "/home/user/app1"),
			createMockContainer("id2", "app1_devcontainer-sfa.db-1", "Up 5 minutes", ""),
		},
		"/home/user/app2": {
			createMockContainer("id3", "app2_devcontainer-sfa.app-1", "Up 10 minutes", "/home/user/app2"),
			createMockContainer("id4", "app2_devcontainer-sfa.db-1", "Up 10 minutes", ""),
			createMockContainer("id5", "app2_devcontainer-sfa.cache-1", "Up 10 minutes", ""),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, "/home/user/app1")
	assert.Contains(t, got, "/home/user/app2")
	assert.Contains(t, got, "app1_devcontainer-sfa.app-1")
	assert.Contains(t, got, "app2_devcontainer-sfa.app-1")
}

func TestFormatGroupedContainers_EmptyGroupsMap_ReturnsNoContainersMessage(t *testing.T) {
	grouped := map[string][]*Container{}

	got := FormatGroupedContainers(grouped)

	assert.Equal(t, "Nenhum DevContainer ativo encontrado.", got)
}

func TestFormatGroupedContainers_GroupWithDefaultBucketName_DisplayedCorrectly(t *testing.T) {
	r := require.New(t)
	grouped := map[string][]*Container{
		"[sem pasta local mapeada]": {
			createMockContainer("orphan123", "orphan_devcontainer-sfa.app-1", "Up 5 minutes", ""),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, "[sem pasta local mapeada]")
	assert.Contains(t, got, "orphan_devcontainer-sfa.app-1")
}

func TestFormatGroupedContainers_OnlyDefaultBucketContainers_FormattedWithCorrectHeader(t *testing.T) {
	r := require.New(t)
	grouped := map[string][]*Container{
		"[sem pasta local mapeada]": {
			createMockContainer("id1", "orphan1_devcontainer-sfa.app-1", "Up 5 minutes", ""),
			createMockContainer("id2", "orphan2_devcontainer-sfa.app-1", "Up 5 minutes", ""),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, "[sem pasta local mapeada]")
	// Verify contains the expected content rather than exact line count
	assert.Contains(t, got, "orphan1_devcontainer-sfa.app-1")
	assert.Contains(t, got, "orphan2_devcontainer-sfa.app-1")
}

func TestFormatGroupedContainers_VeryLongContainerNames_FormattedWithoutBreakage(t *testing.T) {
	r := require.New(t)
	longName := "my-extremely-long-application-name-with-many-parts_devcontainer-sfa.application-instance-1"
	grouped := map[string][]*Container{
		"/home/user/app": {
			createMockContainer("abc123", longName, "Up 5 minutes", "/home/user/app"),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, longName)
	assert.Contains(t, got, "---")
}

func TestFormatGroupedContainers_VeryLongLocalFolderPath_DisplayedCompletely(t *testing.T) {
	r := require.New(t)
	longPath := "/home/user/projects/company/department/team/project-name/subproject/component/src"
	grouped := map[string][]*Container{
		longPath: {
			createMockContainer("abc123", "app_devcontainer-sfa.app-1", "Up 5 minutes", longPath),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, longPath)
}

func TestFormatGroupedContainers_ContainerWithSpecialCharactersInName_PreservedInOutput(t *testing.T) {
	r := require.New(t)
	specialName := "my-app_v2.0_devcontainer-sfa.component@v1-test_1"
	grouped := map[string][]*Container{
		"/home/user/app": {
			createMockContainer("abc123", specialName, "Up 5 minutes", "/home/user/app"),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, specialName)
}

func TestFormatGroupedContainers_LocalFolderWithSpecialPathCharacters_DisplayedAsIs(t *testing.T) {
	r := require.New(t)
	specialPath := "/home/user/my-projects/app_v2.0-beta/src-folder"
	grouped := map[string][]*Container{
		specialPath: {
			createMockContainer("abc123", "app_devcontainer-sfa.app-1", "Up 5 minutes", specialPath),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, specialPath)
}

func TestFormatGroupedContainers_OutputAlwaysContainsAllHeaders(t *testing.T) {
	r := require.New(t)
	grouped := map[string][]*Container{
		"/home/user/app": {
			createMockContainer("abc123", "app_devcontainer-sfa.app-1", "Up 5 minutes", "/home/user/app"),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	assert.Contains(t, got, "CONTAINER ID")
	assert.Contains(t, got, "NAMES")
	assert.Contains(t, got, "STATUS")
	assert.Contains(t, got, "FOLDER")
}

func TestFormatGroupedContainers_OutputAlwaysContainsSeparators(t *testing.T) {
	r := require.New(t)
	grouped := map[string][]*Container{
		"/home/user/app": {
			createMockContainer("abc123", "app_devcontainer-sfa.app-1", "Up 5 minutes", "/home/user/app"),
		},
	}

	got := FormatGroupedContainers(grouped)

	r.NotEmpty(got)
	// Should have at least 2 separators (one after folder header, one after containers)
	separatorCount := strings.Count(got, "---")
	assert.GreaterOrEqual(t, separatorCount, 2)
}
