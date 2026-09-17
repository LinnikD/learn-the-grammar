package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("writing test config file: %v", err)
	}

	return path
}

func TestLoad_DefaultsOnly(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Port)
	}
}

func TestLoad_FileOverridesPort(t *testing.T) {
	path := writeFile(t, "port: 9090\n")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090", cfg.Port)
	}
}

func TestLoad_FileWithUnknownFieldIsIgnored(t *testing.T) {
	path := writeFile(t, "port: 9090\nunknown_field: something\n")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090", cfg.Port)
	}
}

func TestLoad_MissingFilePath(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("expected an error for a missing config file, got nil")
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	path := writeFile(t, "port: [this is not an int\n")

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected an error for malformed YAML, got nil")
	}
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	path := writeFile(t, "port: 9090\n")
	t.Setenv("LTG_PORT", "7070")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 7070 {
		t.Errorf("Port = %d, want 7070", cfg.Port)
	}
}

func TestLoad_EnvOverridesDefaults(t *testing.T) {
	t.Setenv("LTG_PORT", "7070")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 7070 {
		t.Errorf("Port = %d, want 7070", cfg.Port)
	}
}

func TestLoad_InvalidEnvValue(t *testing.T) {
	t.Setenv("LTG_PORT", "not-a-number")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected an error for an invalid LTG_PORT value, got nil")
	}
}
