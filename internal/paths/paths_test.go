package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAbsPath_EmptyReturnsCWD(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	got, err := GetAbsPath("")
	if err != nil {
		t.Fatalf("GetAbsPath('') returned error: %v", err)
	}
	if got != cwd {
		t.Errorf("expected %q, got %q", cwd, got)
	}
}

func TestGetAbsPath_DotReturnsCWD(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	got, err := GetAbsPath(".")
	if err != nil {
		t.Fatalf("GetAbsPath('.') returned error: %v", err)
	}
	if got != cwd {
		t.Errorf("expected %q, got %q", cwd, got)
	}
}

func TestGetAbsPath_RelativePath(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	got, err := GetAbsPath("subdir")
	if err != nil {
		t.Fatalf("GetAbsPath('subdir') returned error: %v", err)
	}
	expected := filepath.Join(cwd, "subdir")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGetAbsPath_AbsolutePathUnchanged(t *testing.T) {
	absPath := "/tmp/test"
	got, err := GetAbsPath(absPath)
	if err != nil {
		t.Fatalf("GetAbsPath returned error: %v", err)
	}
	if got != absPath {
		t.Errorf("expected %q, got %q", absPath, got)
	}
}

func TestGetHostPath_NoWSLReturnsOriginal(t *testing.T) {
	t.Setenv("WSL_DISTRO_NAME", "")
	path := "/tmp/my-project"
	if got := GetHostPath(path); got != path {
		t.Errorf("expected same path %q, got %q", path, got)
	}
}

func TestGetHostPath_WithWSL(t *testing.T) {
	bin := t.TempDir()
	t.Setenv("WSL_DISTRO_NAME", "Ubuntu")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	wslpathScript := `#!/bin/sh
printf 'C:\\Repo\\project\n'
`
	wslpathPath := filepath.Join(bin, "wslpath")
	if err := os.WriteFile(wslpathPath, []byte(wslpathScript), 0755); err != nil {
		t.Fatalf("failed to write wslpath: %v", err)
	}

	got := GetHostPath("/home/user/project")
	want := "C:\\Repo\\project"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
