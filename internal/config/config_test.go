package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetViper resets viper state between tests
func resetViper() {
	viper.Reset()
	globalConfig = nil
}

func TestInitConfig_DefaultValues(t *testing.T) {
	resetViper()

	// Use temp directory as home
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("EDITOR", "nano")

	// Change to temp directory (no config files)
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, "nano", cfg.Editor) // From EDITOR env var
	assert.Equal(t, "text", cfg.OutputFormat)
	assert.False(t, cfg.NoColor)
	assert.False(t, cfg.Verbose)
	assert.Contains(t, cfg.DatabasePath, ".dsa")
}

func TestInitConfig_GlobalConfig(t *testing.T) {
	resetViper()

	// Create temp home directory with config
	tmpHome := t.TempDir()
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	configContent := `
editor: vim
output_format: json
no_color: true
verbose: true
`
	configPath := filepath.Join(dsaDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	t.Setenv("HOME", tmpHome)

	// Change to temp directory (no project config)
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, "vim", cfg.Editor)
	assert.Equal(t, "json", cfg.OutputFormat)
	assert.True(t, cfg.NoColor)
	assert.True(t, cfg.Verbose)
}

func TestInitConfig_ProjectConfig(t *testing.T) {
	resetViper()

	// Create temp project directory with config
	tmpProject := t.TempDir()
	dsaDir := filepath.Join(tmpProject, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	configContent := `
editor: code
output_format: text
no_color: false
`
	configPath := filepath.Join(dsaDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Change to project directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpProject)

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, "code", cfg.Editor)
	assert.Equal(t, "text", cfg.OutputFormat)
	assert.False(t, cfg.NoColor)
}

func TestInitConfig_ProjectOverridesGlobal(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	tmpProject := t.TempDir()

	// Create global config
	globalDsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(globalDsaDir, 0755))
	globalConfig := `
editor: vim
output_format: json
no_color: true
`
	require.NoError(t, os.WriteFile(filepath.Join(globalDsaDir, "config.yaml"), []byte(globalConfig), 0644))

	// Create project config (overrides editor and output_format)
	projectDsaDir := filepath.Join(tmpProject, ".dsa")
	require.NoError(t, os.MkdirAll(projectDsaDir, 0755))
	projectConfig := `
editor: code
output_format: text
`
	require.NoError(t, os.WriteFile(filepath.Join(projectDsaDir, "config.yaml"), []byte(projectConfig), 0644))

	t.Setenv("HOME", tmpHome)

	// Change to project directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpProject)

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, "code", cfg.Editor)        // Project overrides global
	assert.Equal(t, "text", cfg.OutputFormat)  // Project overrides global
	assert.True(t, cfg.NoColor)                // From global (not overridden)
}

func TestInitConfig_EnvironmentVariables(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Set environment variables
	t.Setenv("DSA_EDITOR", "emacs")
	t.Setenv("DSA_OUTPUT_FORMAT", "json")
	t.Setenv("DSA_NO_COLOR", "true")
	t.Setenv("DSA_VERBOSE", "true")

	// Change to temp directory (no config files)
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, "emacs", cfg.Editor)
	assert.Equal(t, "json", cfg.OutputFormat)
	assert.True(t, cfg.NoColor)
	assert.True(t, cfg.Verbose)
}

func TestInitConfig_EnvOverridesConfig(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	// Create config file
	configContent := `
editor: vim
output_format: text
`
	require.NoError(t, os.WriteFile(filepath.Join(dsaDir, "config.yaml"), []byte(configContent), 0644))

	t.Setenv("HOME", tmpHome)
	t.Setenv("DSA_EDITOR", "emacs") // Env var should override config

	// Change to temp directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, "emacs", cfg.Editor)      // Env var wins
	assert.Equal(t, "text", cfg.OutputFormat) // From config (not overridden)
}

func TestInitConfig_MissingConfigFile(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Change to temp directory (no config files)
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	// Should not error when config file is missing
	err := InitConfig()
	require.NoError(t, err)

	// Should have defaults
	cfg := Get()
	assert.NotNil(t, cfg)
	assert.Equal(t, "text", cfg.OutputFormat)
}

func TestInitConfig_InvalidOutputFormat(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	// Create config with invalid output_format
	configContent := `
output_format: xml
`
	require.NoError(t, os.WriteFile(filepath.Join(dsaDir, "config.yaml"), []byte(configContent), 0644))

	t.Setenv("HOME", tmpHome)

	// Change to temp directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "output_format")
	assert.Contains(t, err.Error(), "xml")
}

func TestInitConfig_ValidOutputFormats(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"text format", "text"},
		{"json format", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetViper()

			tmpHome := t.TempDir()
			dsaDir := filepath.Join(tmpHome, ".dsa")
			require.NoError(t, os.MkdirAll(dsaDir, 0755))

			configContent := "output_format: " + tt.format + "\n"
			require.NoError(t, os.WriteFile(filepath.Join(dsaDir, "config.yaml"), []byte(configContent), 0644))

			t.Setenv("HOME", tmpHome)

			origDir, _ := os.Getwd()
			defer os.Chdir(origDir)
			os.Chdir(t.TempDir())

			err := InitConfig()
			require.NoError(t, err)

			cfg := Get()
			assert.Equal(t, tt.format, cfg.OutputFormat)
		})
	}
}

func TestInitConfig_DatabasePathValidation(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Change to temp directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	// Should have created .dsa directory
	dbDir := filepath.Dir(cfg.DatabasePath)
	_, err = os.Stat(dbDir)
	assert.NoError(t, err, "Database directory should exist")
}

func TestInitConfig_MalformedYAML(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	// Create malformed YAML
	malformedContent := `
editor: vim
output_format: [invalid yaml
`
	require.NoError(t, os.WriteFile(filepath.Join(dsaDir, "config.yaml"), []byte(malformedContent), 0644))

	t.Setenv("HOME", tmpHome)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.Error(t, err)
	// Error message should indicate config file read failure
	assert.Contains(t, err.Error(), "failed to read")
}

func TestGet_InitializesIfNil(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Change to temp directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	// Get should initialize config if nil
	cfg := Get()
	assert.NotNil(t, cfg)
	assert.Equal(t, "text", cfg.OutputFormat)
}

func TestGetString_FlagOverride(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	viper.Set("editor", "vim")

	// Flag value should override config
	result := GetString("editor", "code")
	assert.Equal(t, "code", result)

	// Empty flag should use config
	result = GetString("editor", "")
	assert.Equal(t, "vim", result)
}

func TestGetBool_FlagOverride(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	viper.Set("verbose", false)

	// Flag set to true should override config
	result := GetBool("verbose", true, true)
	assert.True(t, result)

	// Flag not set should use config
	result = GetBool("verbose", false, false)
	assert.False(t, result)
}

func TestInitConfig_CustomDatabasePath(t *testing.T) {
	resetViper()

	tmpHome := t.TempDir()
	tmpDB := filepath.Join(t.TempDir(), "custom", "path")
	dsaDir := filepath.Join(tmpHome, ".dsa")
	require.NoError(t, os.MkdirAll(dsaDir, 0755))

	configContent := "database_path: " + tmpDB + "/dsa.db\n"
	require.NoError(t, os.WriteFile(filepath.Join(dsaDir, "config.yaml"), []byte(configContent), 0644))

	t.Setenv("HOME", tmpHome)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(t.TempDir())

	err := InitConfig()
	require.NoError(t, err)

	cfg := Get()
	assert.Equal(t, tmpDB+"/dsa.db", cfg.DatabasePath)

	// Verify directory was created
	_, err = os.Stat(tmpDB)
	assert.NoError(t, err, "Custom database directory should exist")
}
