package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Operator.Handle != "smford" {
		t.Fatalf("expected handle smford, got %s", cfg.Operator.Handle)
	}
	if len(cfg.Repositories) != 27 {
		t.Fatalf("expected 27 repositories, got %d", len(cfg.Repositories))
	}
	if len(cfg.BuildLogs) != 5 {
		t.Fatalf("expected 5 build logs, got %d", len(cfg.BuildLogs))
	}
}

func TestYAMLRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "config.yaml")

	cfg := DefaultConfig()
	if err := Save(yamlPath, &cfg); err != nil {
		t.Fatalf("failed to save YAML config: %v", err)
	}

	loaded, err := Load(yamlPath)
	if err != nil {
		t.Fatalf("failed to load YAML config: %v", err)
	}

	if loaded.Operator.Title != cfg.Operator.Title {
		t.Errorf("expected title %s, got %s", cfg.Operator.Title, loaded.Operator.Title)
	}
	if len(loaded.Repositories) != len(cfg.Repositories) {
		t.Errorf("expected %d repos, got %d", len(cfg.Repositories), len(loaded.Repositories))
	}
}

func TestJSONRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "config.json")

	cfg := DefaultConfig()
	if err := Save(jsonPath, &cfg); err != nil {
		t.Fatalf("failed to save JSON config: %v", err)
	}

	loaded, err := Load(jsonPath)
	if err != nil {
		t.Fatalf("failed to load JSON config: %v", err)
	}

	if loaded.Operator.Title != cfg.Operator.Title {
		t.Errorf("expected title %s, got %s", cfg.Operator.Title, loaded.Operator.Title)
	}
}

func TestLoadNonExistent(t *testing.T) {
	_, err := Load(filepath.Join(os.TempDir(), "does-not-exist.yaml"))
	if err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}
}
