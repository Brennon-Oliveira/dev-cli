package update

type ReleaseInfo struct {
	TagName string
	Assets  []Asset
}

type Asset struct {
	Name string
	URL  string
}

type Updater interface {
	CheckForUpdate(currentVersion string) (*ReleaseInfo, error)
	Download(assetURL string) (string, error)
	Extract(archivePath string) (string, error)
}

type GitHubUpdater struct{}

func NewGitHubUpdater() *GitHubUpdater {
	return &GitHubUpdater{}
}

type MockUpdater struct {
	CheckResult  *ReleaseInfo
	CheckErr     error
	DownloadPath string
	DownloadErr  error
	ExtractPath  string
	ExtractErr   error
}

func NewMockUpdater() *MockUpdater {
	return &MockUpdater{}
}

func (m *MockUpdater) CheckForUpdate(currentVersion string) (*ReleaseInfo, error) {
	return m.CheckResult, m.CheckErr
}

func (m *MockUpdater) Download(assetURL string) (string, error) {
	return m.DownloadPath, m.DownloadErr
}

func (m *MockUpdater) Extract(archivePath string) (string, error) {
	return m.ExtractPath, m.ExtractErr
}
