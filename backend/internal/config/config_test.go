package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o644))

	return path
}

func TestLoad_DefaultsOnly(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port)
	assert.Empty(t, cfg.SessionSecret)
}

func TestLoad_FileOverridesPort(t *testing.T) {
	path := writeFile(t, "port: 9090\n")

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, 9090, cfg.Port)
}

func TestLoad_FileWithUnknownFieldIsIgnored(t *testing.T) {
	path := writeFile(t, "port: 9090\nunknown_field: something\n")

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, 9090, cfg.Port)
}

func TestLoad_MissingFilePath(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	assert.Error(t, err)
}

func TestLoad_MalformedYAML(t *testing.T) {
	path := writeFile(t, "port: [this is not an int\n")

	_, err := Load(path)
	assert.Error(t, err)
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	path := writeFile(t, "port: 9090\n")
	t.Setenv("LTG_PORT", "7070")

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, 7070, cfg.Port)
}

func TestLoad_EnvOverridesDefaults(t *testing.T) {
	t.Setenv("LTG_PORT", "7070")

	cfg, err := Load("")
	require.NoError(t, err)
	assert.Equal(t, 7070, cfg.Port)
}

func TestLoad_InvalidEnvValue(t *testing.T) {
	t.Setenv("LTG_PORT", "not-a-number")

	_, err := Load("")
	assert.Error(t, err)
}

func TestLoad_SessionSecretFromFile(t *testing.T) {
	path := writeFile(t, "port: 9090\nsession_secret: from-file\n")

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "from-file", cfg.SessionSecret)
}

func TestLoad_SessionSecretEnvOverridesFile(t *testing.T) {
	path := writeFile(t, "port: 9090\nsession_secret: from-file\n")
	t.Setenv("LTG_SESSION_SECRET", "from-env")

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "from-env", cfg.SessionSecret)
}

func TestLoad_DatabaseURLFromFile(t *testing.T) {
	path := writeFile(t, "database_url: postgres://from-file\n")

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "postgres://from-file", cfg.DatabaseURL)
}

func TestLoad_DatabaseURLEnvOverridesFile(t *testing.T) {
	path := writeFile(t, "database_url: postgres://from-file\n")
	t.Setenv("LTG_DATABASE_URL", "postgres://from-env")

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "postgres://from-env", cfg.DatabaseURL)
}
