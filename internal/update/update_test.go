package update

import (
	"strings"
	"testing"
)

func TestNewGitHubUpdater(t *testing.T) {
	updater := NewGitHubUpdater()
	if updater == nil {
		t.Error("expected updater to be created")
	}
}

func TestMockUpdater(t *testing.T) {
	mock := NewMockUpdater()

	mock.CheckResult = &ReleaseInfo{TagName: "v1.0.0"}
	mock.DownloadPath = "/tmp/test.tar.gz"
	mock.ExtractPath = "/tmp/dev_new"

	info, err := mock.CheckForUpdate("v0.9.0")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if info.TagName != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", info.TagName)
	}

	path, err := mock.Download("http://example.com/test.tar.gz")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if path != "/tmp/test.tar.gz" {
		t.Errorf("unexpected path: %s", path)
	}

	extractPath, err := mock.Extract("/tmp/test.tar.gz")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if extractPath != "/tmp/dev_new" {
		t.Errorf("unexpected extract path: %s", extractPath)
	}
}

func TestGetAssetName(t *testing.T) {
	name := getAssetName("v1.0.0")
	if name == "" {
		t.Error("expected non-empty asset name")
	}
	if !strings.HasPrefix(name, "dev-") {
		t.Errorf("expected asset name to start with dev-, got %s", name)
	}
}

func TestIsWindowsArchive(t *testing.T) {
	if isWindowsArchive("test.tar.gz") {
		t.Error("expected tar.gz not to be windows archive")
	}
	if !isWindowsArchive("test.zip") {
		t.Error("expected zip to be windows archive")
	}
}
