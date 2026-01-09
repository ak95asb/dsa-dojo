package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests verify end-to-end config loading scenarios

func TestIntegration_ConfigPrecedence(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	tmpProject := t.TempDir()

	// Create global config
	globalDsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(globalDsaDir, 0755))
	globalConfig := `
editor: vim
output_format: json
no_color: false
verbose: false
`
	require.NoError(t, os.WriteFile(filepath.Join(globalDsaDir, "config.yaml"), []byte(globalConfig), 0644))

	// Create project config
	projectDsaDir := filepath.Join(tmpProject, ".dsa")
	require.NoError(t, os.MkdirAll(projectDsaDir, 0755))
	projectConfig := `
output_format: text
verbose: true
`
	require.NoError(t, os.WriteFile(filepath.Join(projectDsaDir, "config.yaml"), []byte(projectConfig), 0644))

	t.Setenv("HOME", tmpHome)
	t.Setenv("DSA_NO_COLOR", "true") // Env var override

	// Change to project directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpProject)

	// Bind flag
	viper.Set("editor", "code") // Flag override

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()

	// Verify precedence: flags > env > project > global > defaults
	assert.Equal(t, "code", cfg.Editor)       // From flag (highest)
	assert.True(t, cfg.NoColor)               // From env var
	assert.Equal(t, "text", cfg.OutputFormat) // From project config
	assert.True(t, cfg.Verbose)               // From project config
}

func TestIntegration_RealFileLoading(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	// Create real config file
	configContent := `
# DSA Configuration
editor: code
output_format: json
no_color: true
database_path: /tmp/test.db
verbose: true
`
	configPath := filepath.Join(dsaDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	t.Setenv("HOME", tmpHome)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	// Verify all settings loaded correctly
	cfg := Get()
	assert.Equal(t, "code", cfg.Editor)
	assert.Equal(t, "json", cfg.OutputFormat)
	assert.True(t, cfg.NoColor)
	assert.True(t, cfg.Verbose)
	assert.Equal(t, "/tmp/test.db", cfg.DatabasePath)
}

func TestIntegration_ConfigChanges(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	// Create initial config
	configPath := filepath.Join(dsaDir, "config.yaml")
	initialConfig := `
editor: vim
output_format: text
`
	require.NoError(t, os.WriteFile(configPath, []byte(initialConfig), 0644))

	t.Setenv("HOME", tmpHome)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	// First load
	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, "vim", cfg.Editor)
	assert.Equal(t, "text", cfg.OutputFormat)

	// Simulate config file change
	resetViper()
	updatedConfig := `
editor: code
output_format: json
`
	require.NoError(t, os.WriteFile(configPath, []byte(updatedConfig), 0644))

	// Reload config
	err = InitConfig()
	require.NoError(t, err)

	cfg = Get()
	assert.Equal(t, "code", cfg.Editor)
	assert.Equal(t, "json", cfg.OutputFormat)
}

func TestIntegration_InvalidConfigHandling(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	// Create config with invalid value
	invalidConfig := `
editor: vim
output_format: invalid_format
`
	configPath := filepath.Join(dsaDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(invalidConfig), 0644))

	t.Setenv("HOME", tmpHome)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "output_format")
	assert.Contains(t, err.Error(), "invalid_format")
}

func TestIntegration_EnvironmentVariableOverrides(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	// Create config file
	configContent := `
editor: vim
output_format: text
no_color: false
verbose: false
`
	configPath := filepath.Join(dsaDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	t.Setenv("HOME", tmpHome)

	// Set all environment variables
	t.Setenv("DSA_EDITOR", "emacs")
	t.Setenv("DSA_OUTPUT_FORMAT", "json")
	t.Setenv("DSA_NO_COLOR", "true")
	t.Setenv("DSA_VERBOSE", "true")

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	// All values should come from environment variables
	cfg := Get()
	assert.Equal(t, "emacs", cfg.Editor)
	assert.Equal(t, "json", cfg.OutputFormat)
	assert.True(t, cfg.NoColor)
	assert.True(t, cfg.Verbose)
}

func TestIntegration_NoConfigFiles(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("EDITOR", "nano")

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	// Should not error when no config files exist
	err := InitConfig()
	require.NoError(t, err)

	// Should use defaults
	cfg := Get()
	assert.Equal(t, "nano", cfg.Editor) // From EDITOR env var
	assert.Equal(t, "text", cfg.OutputFormat)
	assert.False(t, cfg.NoColor)
	assert.False(t, cfg.Verbose)
}

func TestIntegration_DatabasePathCreation(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	customDBPath := filepath.Join(t.TempDir(), "nested", "path", "dsa.db")

	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	configContent := "database_path: " + customDBPath + "\n"
	configPath := filepath.Join(dsaDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	t.Setenv("HOME", tmpHome)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, customDBPath, cfg.DatabasePath)

	// Verify directory was created
	dbDir := filepath.Dir(customDBPath)
	_, err = os.Stat(dbDir)
	assert.NoError(t, err, "Database directory should have been created")
}
