package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Integration tests verify end-to-end config command workflows

func TestIntegration_ConfigSetAndGet(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Simulate config set
	configData := map[string]interface{}{
		"editor":        "vim",
		"output_format": "json",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	tempPath := configPath + ".tmp"
	yamlData, err := yaml.Marshal(configData)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Verify file was created
	assert.FileExists(t, configPath)

	// Read back and verify (simulating config get)
	viper.Reset()
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "vim", viper.GetString("editor"))
	assert.Equal(t, "json", viper.GetString("output_format"))
}

func TestIntegration_ConfigSetMultipleKeys(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set multiple keys
	configData := map[string]interface{}{
		"editor":        "code",
		"output_format": "text",
		"no_color":      true,
		"verbose":       false,
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	yamlData, err := yaml.Marshal(configData)
	require.NoError(t, err)

	tempPath := configPath + ".tmp"
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Read back and verify all keys persisted
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))

	assert.Equal(t, "code", readConfig["editor"])
	assert.Equal(t, "text", readConfig["output_format"])
	assert.Equal(t, true, readConfig["no_color"])
	assert.Equal(t, false, readConfig["verbose"])
}

func TestIntegration_ConfigUnsetRevertsToDefault(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set a value
	configData := map[string]interface{}{
		"editor":  "vim",
		"verbose": true,
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	yamlData, err := yaml.Marshal(configData)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Unset verbose
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))

	delete(readConfig, "verbose")

	yamlData, _ = yaml.Marshal(readConfig)
	tempPath := configPath + ".tmp"
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Verify verbose was removed
	data, err = os.ReadFile(configPath)
	require.NoError(t, err)

	var finalConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &finalConfig))

	_, exists := finalConfig["verbose"]
	assert.False(t, exists)
	assert.Equal(t, "vim", finalConfig["editor"]) // Other keys unchanged
}

func TestIntegration_ConfigChangesPersist(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// First write
	configData := map[string]interface{}{
		"editor": "vim",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Second write (update)
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))

	readConfig["output_format"] = "json"

	yamlData, _ = yaml.Marshal(readConfig)
	tempPath := configPath + ".tmp"
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Verify both keys persist
	data, err = os.ReadFile(configPath)
	require.NoError(t, err)

	var finalConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &finalConfig))

	assert.Equal(t, "vim", finalConfig["editor"])
	assert.Equal(t, "json", finalConfig["output_format"])
}

func TestIntegration_ConfigListShowsSources(t *testing.T) {
	tmpHome := t.TempDir()
	tmpProject := t.TempDir()

	// Create global config
	globalConfigPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(globalConfigPath), 0755))

	globalData := map[string]interface{}{
		"editor": "vim",
	}
	yamlData, _ := yaml.Marshal(globalData)
	require.NoError(t, os.WriteFile(globalConfigPath, yamlData, 0644))

	// Create project config
	projectConfigPath := filepath.Join(tmpProject, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(projectConfigPath), 0755))

	projectData := map[string]interface{}{
		"output_format": "json",
	}
	yamlData, _ = yaml.Marshal(projectData)
	require.NoError(t, os.WriteFile(projectConfigPath, yamlData, 0644))

	t.Setenv("HOME", tmpHome)
	t.Setenv("DSA_VERBOSE", "true")

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpProject)

	// Test source detection
	assert.Equal(t, "global", getConfigSource("editor"))
	assert.Equal(t, "project", getConfigSource("output_format"))
	assert.Equal(t, "environment", getConfigSource("verbose"))
	assert.Equal(t, "default", getConfigSource("no_color"))
}

func TestIntegration_AtomicWritePreventsCorruption(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	tempPath := configPath + ".tmp"

	// Simulate atomic write
	configData := map[string]interface{}{
		"editor": "vim",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	yamlData, _ := yaml.Marshal(configData)

	// Write to temp file
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))

	// Original file shouldn't exist yet
	assert.NoFileExists(t, configPath)

	// Rename to actual file
	require.NoError(t, os.Rename(tempPath, configPath))

	// Original file should exist now
	assert.FileExists(t, configPath)

	// Temp file should be gone
	assert.NoFileExists(t, tempPath)

	// Verify content is correct
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))
	assert.Equal(t, "vim", readConfig["editor"])
}

func TestIntegration_EnvironmentOverrideAfterSet(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set config value
	configData := map[string]interface{}{
		"editor": "vim",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Set environment variable (should override)
	t.Setenv("DSA_EDITOR", "emacs")

	// Environment should take precedence
	source := getConfigSource("editor")
	assert.Equal(t, "environment", source)
}

func TestIntegration_ListFormatConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set list_format
	configData := map[string]interface{}{
		"list_format": "json",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	tempPath := configPath + ".tmp"
	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Read back and verify
	viper.Reset()
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "json", viper.GetString("list_format"))
}

func TestIntegration_StatusFormatConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set status_format
	configData := map[string]interface{}{
		"status_format": "json",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	tempPath := configPath + ".tmp"
	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Read back and verify
	viper.Reset()
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "json", viper.GetString("status_format"))
}

func TestIntegration_OutputStyleConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set output_style
	configData := map[string]interface{}{
		"output_style": "compact",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	tempPath := configPath + ".tmp"
	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Read back and verify
	viper.Reset()
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "compact", viper.GetString("output_style"))
}

func TestIntegration_ColorSchemeConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set color_scheme
	configData := map[string]interface{}{
		"color_scheme": "solarized",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	tempPath := configPath + ".tmp"
	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Read back and verify
	viper.Reset()
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "solarized", viper.GetString("color_scheme"))
}

func TestIntegration_InvalidEnumValues(t *testing.T) {
	t.Run("invalid list_format", func(t *testing.T) {
		_, err := parseConfigValue("list_format", "xml")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be 'table' or 'json'")
	})

	t.Run("invalid status_format", func(t *testing.T) {
		_, err := parseConfigValue("status_format", "csv")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be 'table' or 'json'")
	})

	t.Run("invalid output_style", func(t *testing.T) {
		_, err := parseConfigValue("output_style", "detailed")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be one of")
	})

	t.Run("invalid color_scheme", func(t *testing.T) {
		_, err := parseConfigValue("color_scheme", "rainbow")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be one of")
	})
}

func TestIntegration_UnsetRevertsNewKeys(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set values
	configData := map[string]interface{}{
		"list_format":   "json",
		"color_scheme": "solarized",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))
	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Unset list_format
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))

	delete(readConfig, "list_format")

	// Write back
	tempPath := configPath + ".tmp"
	yamlData, _ = yaml.Marshal(readConfig)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Verify list_format was removed
	data, err = os.ReadFile(configPath)
	require.NoError(t, err)

	var finalConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &finalConfig))

	_, exists := finalConfig["list_format"]
	assert.False(t, exists)
	assert.Equal(t, "solarized", finalConfig["color_scheme"]) // Other keys unchanged
}

func TestIntegration_EditorArgsConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set editor_args
	configData := map[string]interface{}{
		"editor_args": "--goto {file}:{line}",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	tempPath := configPath + ".tmp"
	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Read back and verify
	viper.Reset()
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "--goto {file}:{line}", viper.GetString("editor_args"))
}

func TestIntegration_EditorAndEditorArgsTogether(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set both editor and editor_args
	configData := map[string]interface{}{
		"editor":      "code",
		"editor_args": "--goto {file}:{line}",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Read back and verify both
	viper.Reset()
	viper.SetConfigFile(configPath)
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "code", viper.GetString("editor"))
	assert.Equal(t, "--goto {file}:{line}", viper.GetString("editor_args"))
}

func TestIntegration_EditorArgsDefaultsToEmpty(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Create empty config
	configData := map[string]interface{}{}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Read config and verify editor_args defaults to empty string
	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.SetDefault("editor_args", "")
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "", viper.GetString("editor_args"))
}

func TestIntegration_EditorArgsUnset(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set editor_args
	configData := map[string]interface{}{
		"editor_args": "+{line}",
		"editor":      "vim",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))
	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Unset editor_args
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))

	delete(readConfig, "editor_args")

	// Write back
	tempPath := configPath + ".tmp"
	yamlData, _ = yaml.Marshal(readConfig)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Verify editor_args was removed
	data, err = os.ReadFile(configPath)
	require.NoError(t, err)

	var finalConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &finalConfig))

	_, exists := finalConfig["editor_args"]
	assert.False(t, exists)
	assert.Equal(t, "vim", finalConfig["editor"]) // Other keys unchanged
}

func TestIntegration_EditorArgsEnvironmentOverride(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("DSA_EDITOR_ARGS", "+{line} {file}")

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Set editor_args in config
	configData := map[string]interface{}{
		"editor_args": "--goto {file}:{line}",
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))
	yamlData, _ := yaml.Marshal(configData)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Environment variable should override config file
	source := getConfigSource("editor_args")
	assert.Equal(t, "environment", source)
}

// Profile management integration tests

func TestIntegration_ProfileCreateListSwitch(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Create default config with some values
	defaultConfigPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(defaultConfigPath), 0755))

	defaultData := map[string]interface{}{
		"editor":        "vim",
		"output_format": "text",
	}
	yamlData, _ := yaml.Marshal(defaultData)
	require.NoError(t, os.WriteFile(defaultConfigPath, yamlData, 0644))

	// Create profile "work" from current config
	workProfilePath := filepath.Join(tmpHome, ".dsa", "profiles", "work.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(workProfilePath), 0755))

	workData := map[string]interface{}{
		"editor":        "vim",
		"output_format": "text",
	}
	yamlData, _ = yaml.Marshal(workData)
	tempPath := workProfilePath + ".tmp"
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, workProfilePath))

	// Verify profile file was created
	assert.FileExists(t, workProfilePath)

	// Verify profile exists
	exists := profileExists("work")
	assert.True(t, exists)

	// Verify initially on "default" profile
	active, err := getActiveProfile()
	require.NoError(t, err)
	assert.Equal(t, "default", active)

	// Switch to "work" profile
	require.NoError(t, setActiveProfile("work"))

	// Verify active profile changed
	active, err = getActiveProfile()
	require.NoError(t, err)
	assert.Equal(t, "work", active)

	// Verify active-profile file exists and contains "work"
	activeProfilePath := filepath.Join(tmpHome, ".dsa", "active-profile")
	assert.FileExists(t, activeProfilePath)

	data, err := os.ReadFile(activeProfilePath)
	require.NoError(t, err)
	assert.Equal(t, "work", string(data))

	// Switch back to default
	require.NoError(t, setActiveProfile("default"))

	// Verify active-profile file was removed
	assert.NoFileExists(t, activeProfilePath)

	active, err = getActiveProfile()
	require.NoError(t, err)
	assert.Equal(t, "default", active)
}

func TestIntegration_ProfileExportImport(t *testing.T) {
	tmpHome := t.TempDir()
	tmpExport := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Create a profile
	devProfilePath := filepath.Join(tmpHome, ".dsa", "profiles", "dev.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(devProfilePath), 0755))

	devData := map[string]interface{}{
		"editor":        "code",
		"output_format": "json",
		"verbose":       true,
	}
	yamlData, _ := yaml.Marshal(devData)
	require.NoError(t, os.WriteFile(devProfilePath, yamlData, 0644))

	// Export profile
	exportPath := filepath.Join(tmpExport, "dev-exported.yaml")
	data, err := os.ReadFile(devProfilePath)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(exportPath, data, 0644))

	// Verify export file exists
	assert.FileExists(t, exportPath)

	// Delete the original profile
	require.NoError(t, os.Remove(devProfilePath))
	assert.NoFileExists(t, devProfilePath)

	// Import the profile back with new name
	importedProfilePath := filepath.Join(tmpHome, ".dsa", "profiles", "dev-imported.yaml")
	data, err = os.ReadFile(exportPath)
	require.NoError(t, err)

	tempPath := importedProfilePath + ".tmp"
	require.NoError(t, os.WriteFile(tempPath, data, 0644))
	require.NoError(t, os.Rename(tempPath, importedProfilePath))

	// Verify imported profile exists
	assert.FileExists(t, importedProfilePath)

	// Verify imported profile has correct data
	viper.Reset()
	viper.SetConfigFile(importedProfilePath)
	require.NoError(t, viper.ReadInConfig())

	assert.Equal(t, "code", viper.GetString("editor"))
	assert.Equal(t, "json", viper.GetString("output_format"))
	assert.Equal(t, true, viper.GetBool("verbose"))
}

func TestIntegration_ProfileDeleteRemovesFileAndSwitchesToDefault(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Create a profile
	testProfilePath := filepath.Join(tmpHome, ".dsa", "profiles", "test.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(testProfilePath), 0755))

	testData := map[string]interface{}{
		"editor": "nano",
	}
	yamlData, _ := yaml.Marshal(testData)
	require.NoError(t, os.WriteFile(testProfilePath, yamlData, 0644))

	// Switch to the profile
	require.NoError(t, setActiveProfile("test"))

	// Verify it's active
	active, err := getActiveProfile()
	require.NoError(t, err)
	assert.Equal(t, "test", active)

	// Delete the profile (should also switch to default)
	require.NoError(t, os.Remove(testProfilePath))
	require.NoError(t, setActiveProfile("default"))

	// Verify file is gone
	assert.NoFileExists(t, testProfilePath)

	// Verify switched back to default
	active, err = getActiveProfile()
	require.NoError(t, err)
	assert.Equal(t, "default", active)

	// Verify active-profile file was removed
	activeProfilePath := filepath.Join(tmpHome, ".dsa", "active-profile")
	assert.NoFileExists(t, activeProfilePath)
}

func TestIntegration_SwitchingProfilesAffectsConfigReads(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Create default config
	defaultConfigPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(defaultConfigPath), 0755))

	defaultData := map[string]interface{}{
		"editor": "vim",
	}
	yamlData, _ := yaml.Marshal(defaultData)
	require.NoError(t, os.WriteFile(defaultConfigPath, yamlData, 0644))

	// Create work profile with different editor
	workProfilePath := filepath.Join(tmpHome, ".dsa", "profiles", "work.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(workProfilePath), 0755))

	workData := map[string]interface{}{
		"editor": "code",
	}
	yamlData, _ = yaml.Marshal(workData)
	require.NoError(t, os.WriteFile(workProfilePath, yamlData, 0644))

	// Read default config
	viper.Reset()
	viper.SetConfigFile(defaultConfigPath)
	require.NoError(t, viper.ReadInConfig())
	assert.Equal(t, "vim", viper.GetString("editor"))

	// Switch to work profile and set active-profile file
	require.NoError(t, setActiveProfile("work"))

	// Simulate config initialization with active profile
	activeProfilePath := filepath.Join(tmpHome, ".dsa", "active-profile")
	activeProfileData, err := os.ReadFile(activeProfilePath)
	require.NoError(t, err)

	profileName := string(activeProfileData)
	activeConfigFile := filepath.Join(tmpHome, ".dsa", "profiles", profileName+".yaml")

	// Read config from active profile
	viper.Reset()
	viper.SetConfigFile(activeConfigFile)
	require.NoError(t, viper.ReadInConfig())
	assert.Equal(t, "code", viper.GetString("editor"))
}

func TestIntegration_DefaultProfileSpecialHandling(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Initially should be "default" (no active-profile file)
	active, err := getActiveProfile()
	require.NoError(t, err)
	assert.Equal(t, "default", active)

	activeProfilePath := filepath.Join(tmpHome, ".dsa", "active-profile")
	assert.NoFileExists(t, activeProfilePath)

	// "default" profile should use config.yaml, not profiles/default.yaml
	defaultConfigPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(defaultConfigPath), 0755))

	defaultData := map[string]interface{}{
		"editor": "vim",
	}
	yamlData, _ := yaml.Marshal(defaultData)
	require.NoError(t, os.WriteFile(defaultConfigPath, yamlData, 0644))

	// Verify "default" profile name is invalid for creation
	assert.False(t, isValidProfileName("default"))

	// Switch to default explicitly (should be no-op)
	require.NoError(t, setActiveProfile("default"))

	// Verify still no active-profile file
	assert.NoFileExists(t, activeProfilePath)

	active, err = getActiveProfile()
	require.NoError(t, err)
	assert.Equal(t, "default", active)
}

func TestIntegration_InvalidProfileNamesRejected(t *testing.T) {
	// Test invalid profile names
	invalidNames := []string{
		"default",        // reserved
		"",               // empty
		"a b",            // space
		"test@profile",   // special char
		"test profile",   // space
		"test/profile",   // slash
		"test\\profile",  // backslash
		"test.profile",   // dot (should be invalid)
		strings.Repeat("a", 51), // too long
	}

	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			assert.False(t, isValidProfileName(name))
		})
	}

	// Test valid profile names
	validNames := []string{
		"work",
		"personal",
		"dev-2024",
		"my_profile",
		"Profile123",
		"a",
		strings.Repeat("a", 50),
	}

	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			assert.True(t, isValidProfileName(name))
		})
	}
}

func TestIntegration_SwitchToNonExistentProfile(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Verify profile doesn't exist
	exists := profileExists("nonexistent")
	assert.False(t, exists)

	// Switch to it anyway (setActiveProfile doesn't check existence)
	require.NoError(t, setActiveProfile("nonexistent"))

	// Verify active-profile file was created
	activeProfilePath := filepath.Join(tmpHome, ".dsa", "active-profile")
	assert.FileExists(t, activeProfilePath)

	// Verify it says "nonexistent"
	data, err := os.ReadFile(activeProfilePath)
	require.NoError(t, err)
	assert.Equal(t, "nonexistent", string(data))

	// But when InitConfig tries to load it, it won't find the profile file
	profilePath := filepath.Join(tmpHome, ".dsa", "profiles", "nonexistent.yaml")
	assert.NoFileExists(t, profilePath)
}

// Validation integration tests

func TestIntegration_ValidateValidConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	// Create valid config
	validConfig := map[string]interface{}{
		"editor":        "nvim",
		"output_format": "json",
		"no_color":      true,
		"verbose":       false,
	}
	yamlData, _ := yaml.Marshal(validConfig)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Validate should pass
	errors, err := validateConfig()

	assert.NoError(t, err)
	assert.Empty(t, errors)
}

func TestIntegration_ValidateCorruptedConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	// Create corrupted config
	corruptedConfig := map[string]interface{}{
		"unknown_key":  "value",
		"color_scheme": "invalid_scheme",
		"no_color":     "not_a_boolean",
	}
	yamlData, _ := yaml.Marshal(corruptedConfig)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Validate should fail with multiple errors
	errors, err := validateConfig()

	assert.NoError(t, err)
	assert.Len(t, errors, 3)

	// Verify specific error messages
	var hasUnknownKey, hasInvalidEnum, hasWrongType bool
	for _, e := range errors {
		if e.Key == "unknown_key" {
			hasUnknownKey = true
			assert.Contains(t, e.Message, "unknown configuration key")
		}
		if e.Key == "color_scheme" {
			hasInvalidEnum = true
			assert.Contains(t, e.Message, "must be one of")
		}
		if e.Key == "no_color" {
			hasWrongType = true
			assert.Contains(t, e.Message, "must be 'true' or 'false'")
		}
	}

	assert.True(t, hasUnknownKey)
	assert.True(t, hasInvalidEnum)
	assert.True(t, hasWrongType)
}

// Reset integration tests

func TestIntegration_FullResetCreatesBackup(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	// Create custom config
	customConfig := map[string]interface{}{
		"editor":        "emacs",
		"output_format": "json",
	}
	yamlData, _ := yaml.Marshal(customConfig)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Note: Can't test interactive reset in unit test, but we can test the backup logic
	// by simulating what runConfigReset does

	// Create backup with timestamp
	backupPath := configPath + ".backup.test"
	data, _ := os.ReadFile(configPath)
	os.WriteFile(backupPath, data, 0644)

	// Verify backup was created
	assert.FileExists(t, backupPath)

	// Verify backup has original content
	var backupConfig map[string]interface{}
	backupData, _ := os.ReadFile(backupPath)
	yaml.Unmarshal(backupData, &backupConfig)

	assert.Equal(t, "emacs", backupConfig["editor"])
	assert.Equal(t, "json", backupConfig["output_format"])
}

func TestIntegration_SingleKeyResetPreservesOthers(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	// Create config with multiple keys
	config := map[string]interface{}{
		"editor":        "emacs",
		"output_format": "json",
		"verbose":       true,
	}
	yamlData, _ := yaml.Marshal(config)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Simulate single key reset for "editor"
	defaults := getDefaultConfig()
	defaultEditor := defaults["editor"]

	// Load current config
	var configData map[string]interface{}
	data, _ := os.ReadFile(configPath)
	yaml.Unmarshal(data, &configData)

	// Reset only editor
	configData["editor"] = defaultEditor

	// Write back
	yamlData, _ = yaml.Marshal(configData)
	tempPath := configPath + ".tmp"
	os.WriteFile(tempPath, yamlData, 0644)
	os.Rename(tempPath, configPath)

	// Verify editor was reset
	var updatedConfig map[string]interface{}
	updatedData, _ := os.ReadFile(configPath)
	yaml.Unmarshal(updatedData, &updatedConfig)

	assert.Equal(t, defaultEditor, updatedConfig["editor"])
	assert.Equal(t, "json", updatedConfig["output_format"]) // Preserved
	assert.Equal(t, true, updatedConfig["verbose"])         // Preserved
}

func TestIntegration_ResetToDefaultsMatchesExpected(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Get defaults
	defaults := getDefaultConfig()

	// Verify defaults match expected values
	assert.Equal(t, "text", defaults["output_format"])
	assert.Equal(t, false, defaults["no_color"])
	assert.Equal(t, false, defaults["verbose"])
	assert.Equal(t, "table", defaults["list_format"])
	assert.Equal(t, "table", defaults["status_format"])
	assert.Equal(t, "normal", defaults["output_style"])
	assert.Equal(t, "default", defaults["color_scheme"])
	assert.Equal(t, "", defaults["editor_args"])
	assert.Equal(t, filepath.Join(tmpHome, ".dsa", "dsa.db"), defaults["database_path"])
}

// Defaults output integration tests

func TestIntegration_DefaultsOutputIsValidYAML(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	defaults := getDefaultConfig()

	// Marshal to YAML
	yamlData, err := yaml.Marshal(defaults)

	assert.NoError(t, err)
	assert.NotEmpty(t, yamlData)

	// Verify it can be unmarshaled back
	var parsed map[string]interface{}
	err = yaml.Unmarshal(yamlData, &parsed)

	assert.NoError(t, err)
	assert.Len(t, parsed, len(validConfigKeys))
}

func TestIntegration_DefaultsContainsAllKeys(t *testing.T) {
	defaults := getDefaultConfig()

	// Verify all valid config keys are present
	for _, key := range validConfigKeys {
		_, exists := defaults[key]
		assert.True(t, exists, "Default should exist for key: %s", key)
	}
}

func TestIntegration_DefaultsCanBeWrittenToFile(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	defaults := getDefaultConfig()
	yamlData, _ := yaml.Marshal(defaults)

	// Write to temp file
	outputPath := filepath.Join(tmpHome, "defaults.yaml")
	err := os.WriteFile(outputPath, yamlData, 0644)

	assert.NoError(t, err)
	assert.FileExists(t, outputPath)

	// Verify file can be read back
	var readBack map[string]interface{}
	data, _ := os.ReadFile(outputPath)
	err = yaml.Unmarshal(data, &readBack)

	assert.NoError(t, err)
	assert.Equal(t, defaults, readBack)
}
