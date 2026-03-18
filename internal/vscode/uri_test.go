package vscode

import (
	"encoding/hex"
	"os"
	"testing"
)

func TestGetContainerURI_DefaultWorkspace(t *testing.T) {
	os.Unsetenv("WSL_DISTRO_NAME")
	absPath := "/tmp/project"
	uri := GetContainerURI(absPath, "")

	hexPath := hex.EncodeToString([]byte(absPath))
	want := "vscode-remote://dev-container+" + hexPath + "/workspaces"

	if uri != want {
		t.Errorf("expected %q, got %q", want, uri)
	}
}

func TestGetContainerURI_CustomWorkspace(t *testing.T) {
	os.Unsetenv("WSL_DISTRO_NAME")
	absPath := "/tmp/project"
	uri := GetContainerURI(absPath, "/workspace/custom")

	hexPath := hex.EncodeToString([]byte(absPath))
	want := "vscode-remote://dev-container+" + hexPath + "/workspace/custom//"

	if uri != want {
		t.Errorf("expected %q, got %q", want, uri)
	}
}

func TestGetContainerURI_WorkspaceWithTrailingSlash(t *testing.T) {
	os.Unsetenv("WSL_DISTRO_NAME")
	absPath := "/tmp/project"
	uri := GetContainerURI(absPath, "/workspaces/")

	hexPath := hex.EncodeToString([]byte(absPath))
	want := "vscode-remote://dev-container+" + hexPath + "/workspaces/"

	if uri != want {
		t.Errorf("expected %q, got %q", want, uri)
	}
}

func TestGetContainerURI_WithWSL(t *testing.T) {
	bin := t.TempDir()
	t.Setenv("WSL_DISTRO_NAME", "Ubuntu")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	wslpathScript := `#!/bin/sh
printf 'C:\\Repo\\project\n'
`
	wslpathPath := bin + "/wslpath"
	if err := os.WriteFile(wslpathPath, []byte(wslpathScript), 0755); err != nil {
		t.Fatalf("failed to write wslpath: %v", err)
	}

	absPath := "/home/user/project"
	uri := GetContainerURI(absPath, "")

	hexPath := hex.EncodeToString([]byte("C:\\Repo\\project"))
	want := "vscode-remote://dev-container+" + hexPath + "/workspaces"

	if uri != want {
		t.Errorf("expected %q, got %q", want, uri)
	}
}
