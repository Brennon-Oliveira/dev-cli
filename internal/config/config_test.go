package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultWhenConfigFileMissing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := Load()
	if cfg.Core.Tool != "docker" {
		t.Fatalf("expected default tool docker, got %q", cfg.Core.Tool)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var cfg GlobalConfig
	cfg.Core.Tool = "podman"

	if err := Save(cfg); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded := Load()
	if loaded.Core.Tool != "podman" {
		t.Fatalf("expected tool podman, got %q", loaded.Core.Tool)
	}

	path := filepath.Join(home, ".dev-cli", "config.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to exist at %s: %v", path, err)
	}
}

func TestLoad_EmptyToolFallsBackToDocker(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".dev-cli")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configFile := filepath.Join(path, "config.json")
	if err := os.WriteFile(configFile, []byte(`{"core":{}}`), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	loaded := Load()
	if loaded.Core.Tool != "docker" {
		t.Fatalf("expected fallback tool docker, got %q", loaded.Core.Tool)
	}
}
