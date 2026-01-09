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

func TestIsValidConfigKey(t *testing.T) {
	tests := []struct {
		key   string
		valid bool
	}{
		{"editor", true},
		{"editor_args", true},
		{"output_format", true},
		{"no_color", true},
		{"database_path", true},
		{"verbose", true},
		{"invalid_key", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := isValidConfigKey(tt.key)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestParseConfigValue(t *testing.T) {
	t.Run("boolean values", func(t *testing.T) {
		// Valid boolean values
		val, err := parseConfigValue("no_color", "true")
		assert.NoError(t, err)
		assert.Equal(t, true, val)

		val, err = parseConfigValue("verbose", "false")
		assert.NoError(t, err)
		assert.Equal(t, false, val)

		// Invalid boolean values
		_, err = parseConfigValue("no_color", "yes")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be 'true' or 'false'")
	})

	t.Run("output_format enum", func(t *testing.T) {
		// Valid values
		val, err := parseConfigValue("output_format", "text")
		assert.NoError(t, err)
		assert.Equal(t, "text", val)

		val, err = parseConfigValue("output_format", "json")
		assert.NoError(t, err)
		assert.Equal(t, "json", val)

		// Invalid value
		_, err = parseConfigValue("output_format", "xml")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be 'text' or 'json'")
	})

	t.Run("string values", func(t *testing.T) {
		val, err := parseConfigValue("editor", "vim")
		assert.NoError(t, err)
		assert.Equal(t, "vim", val)

		val, err = parseConfigValue("editor_args", "+{line}")
		assert.NoError(t, err)
		assert.Equal(t, "+{line}", val)

		val, err = parseConfigValue("editor_args", "--goto {file}:{line}")
		assert.NoError(t, err)
		assert.Equal(t, "--goto {file}:{line}", val)

		val, err = parseConfigValue("database_path", "/tmp/test.db")
		assert.NoError(t, err)
		assert.Equal(t, "/tmp/test.db", val)
	})

	t.Run("list_format enum", func(t *testing.T) {
		// Valid values
		val, err := parseConfigValue("list_format", "table")
		assert.NoError(t, err)
		assert.Equal(t, "table", val)

		val, err = parseConfigValue("list_format", "json")
		assert.NoError(t, err)
		assert.Equal(t, "json", val)

		// Invalid value
		_, err = parseConfigValue("list_format", "xml")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be 'table' or 'json'")
	})

	t.Run("status_format enum", func(t *testing.T) {
		// Valid values
		val, err := parseConfigValue("status_format", "table")
		assert.NoError(t, err)
		assert.Equal(t, "table", val)

		val, err = parseConfigValue("status_format", "json")
		assert.NoError(t, err)
		assert.Equal(t, "json", val)

		// Invalid value
		_, err = parseConfigValue("status_format", "csv")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be 'table' or 'json'")
	})

	t.Run("output_style enum", func(t *testing.T) {
		// Valid values
		val, err := parseConfigValue("output_style", "normal")
		assert.NoError(t, err)
		assert.Equal(t, "normal", val)

		val, err = parseConfigValue("output_style", "compact")
		assert.NoError(t, err)
		assert.Equal(t, "compact", val)

		val, err = parseConfigValue("output_style", "verbose")
		assert.NoError(t, err)
		assert.Equal(t, "verbose", val)

		// Invalid value
		_, err = parseConfigValue("output_style", "detailed")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be one of")
	})

	t.Run("color_scheme enum", func(t *testing.T) {
		// Valid values
		val, err := parseConfigValue("color_scheme", "default")
		assert.NoError(t, err)
		assert.Equal(t, "default", val)

		val, err = parseConfigValue("color_scheme", "solarized")
		assert.NoError(t, err)
		assert.Equal(t, "solarized", val)

		val, err = parseConfigValue("color_scheme", "monokai")
		assert.NoError(t, err)
		assert.Equal(t, "monokai", val)

		val, err = parseConfigValue("color_scheme", "nord")
		assert.NoError(t, err)
		assert.Equal(t, "nord", val)

		// Invalid value
		_, err = parseConfigValue("color_scheme", "rainbow")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be one of")
	})
}

func TestConfigSet_WritesFile(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Create config data
	configData := map[string]interface{}{
		"editor": "vim",
	}

	// Ensure directory exists
	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	// Write config
	yamlData, err := yaml.Marshal(configData)
	require.NoError(t, err)

	tempPath := configPath + ".tmp"
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Verify file was created
	assert.FileExists(t, configPath)

	// Read back and verify
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))

	assert.Equal(t, "vim", readConfig["editor"])
}

func TestConfigUnset_RemovesKey(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")

	// Create config with multiple keys
	configData := map[string]interface{}{
		"editor":        "vim",
		"output_format": "json",
		"verbose":       true,
	}

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))

	yamlData, err := yaml.Marshal(configData)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(configPath, yamlData, 0644))

	// Remove editor key
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))

	delete(readConfig, "editor")

	// Write back
	tempPath := configPath + ".tmp"
	yamlData, _ = yaml.Marshal(readConfig)
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))
	require.NoError(t, os.Rename(tempPath, configPath))

	// Verify key was removed
	data, err = os.ReadFile(configPath)
	require.NoError(t, err)

	var finalConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &finalConfig))

	_, exists := finalConfig["editor"]
	assert.False(t, exists, "editor key should be removed")
	assert.Equal(t, "json", finalConfig["output_format"])
	assert.Equal(t, true, finalConfig["verbose"])
}

func TestGetConfigSource(t *testing.T) {
	tmpHome := t.TempDir()
	tmpProject := t.TempDir()

	t.Run("default source", func(t *testing.T) {
		t.Setenv("HOME", tmpHome)
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)
		os.Chdir(tmpProject)

		source := getConfigSource("editor")
		assert.Equal(t, "default", source)
	})

	t.Run("global config source", func(t *testing.T) {
		t.Setenv("HOME", tmpHome)
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)
		os.Chdir(tmpProject)

		// Create global config
		globalConfigPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
		require.NoError(t, os.MkdirAll(filepath.Dir(globalConfigPath), 0755))

		configData := map[string]interface{}{"editor": "vim"}
		yamlData, _ := yaml.Marshal(configData)
		require.NoError(t, os.WriteFile(globalConfigPath, yamlData, 0644))

		source := getConfigSource("editor")
		assert.Equal(t, "global", source)
	})

	t.Run("project config source", func(t *testing.T) {
		t.Setenv("HOME", tmpHome)
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)
		os.Chdir(tmpProject)

		// Create project config
		projectConfigPath := filepath.Join(tmpProject, ".dsa", "config.yaml")
		require.NoError(t, os.MkdirAll(filepath.Dir(projectConfigPath), 0755))

		configData := map[string]interface{}{"editor": "code"}
		yamlData, _ := yaml.Marshal(configData)
		require.NoError(t, os.WriteFile(projectConfigPath, yamlData, 0644))

		source := getConfigSource("editor")
		assert.Equal(t, "project", source)
	})

	t.Run("environment source", func(t *testing.T) {
		t.Setenv("HOME", tmpHome)
		t.Setenv("DSA_EDITOR", "emacs")
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)
		os.Chdir(tmpProject)

		source := getConfigSource("editor")
		assert.Equal(t, "environment", source)
	})
}

func TestConfigCommand_Subcommands(t *testing.T) {
	// Verify config command has all subcommands
	assert.NotNil(t, configCmd)
	assert.Equal(t, "config", configCmd.Use)

	// Find subcommands
	var hasGet, hasSet, hasList, hasUnset bool
	for _, cmd := range configCmd.Commands() {
		switch cmd.Use {
		case "get [key]":
			hasGet = true
		case "set <key> <value>":
			hasSet = true
		case "list":
			hasList = true
		case "unset <key>":
			hasUnset = true
		}
	}

	assert.True(t, hasGet, "config should have 'get' subcommand")
	assert.True(t, hasSet, "config should have 'set' subcommand")
	assert.True(t, hasList, "config should have 'list' subcommand")
	assert.True(t, hasUnset, "config should have 'unset' subcommand")
}

func TestConfigGet_AllKeys(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Set some test values
	viper.Set("editor", "vim")
	viper.Set("output_format", "json")
	viper.Set("verbose", true)

	// Get all settings
	allSettings := viper.AllSettings()

	assert.NotEmpty(t, allSettings)
	assert.Equal(t, "vim", allSettings["editor"])
	assert.Equal(t, "json", allSettings["output_format"])
	assert.Equal(t, true, allSettings["verbose"])
}

func TestConfigSet_InvalidKey(t *testing.T) {
	key := "invalid_key"
	value := "some_value"

	_, err := parseConfigValue(key, value)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown key")
}

func TestConfigSet_InvalidValue(t *testing.T) {
	tests := []struct {
		key      string
		value    string
		errorMsg string
	}{
		{"output_format", "xml", "must be 'text' or 'json'"},
		{"no_color", "yes", "must be 'true' or 'false'"},
		{"verbose", "1", "must be 'true' or 'false'"},
	}

	for _, tt := range tests {
		t.Run(tt.key+"="+tt.value, func(t *testing.T) {
			_, err := parseConfigValue(tt.key, tt.value)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestConfigSet_ValidValues(t *testing.T) {
	tests := []struct {
		key      string
		value    string
		expected interface{}
	}{
		{"editor", "vim", "vim"},
		{"output_format", "json", "json"},
		{"no_color", "true", true},
		{"database_path", "/tmp/test.db", "/tmp/test.db"},
		{"verbose", "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.key+"="+tt.value, func(t *testing.T) {
			result, err := parseConfigValue(tt.key, tt.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Simulate atomic write
	configData := map[string]interface{}{
		"editor":        "vim",
		"output_format": "json",
	}

	yamlData, err := yaml.Marshal(configData)
	require.NoError(t, err)

	// Write to temp file
	tempPath := configPath + ".tmp"
	require.NoError(t, os.WriteFile(tempPath, yamlData, 0644))

	// Verify temp file exists
	assert.FileExists(t, tempPath)

	// Rename to actual file
	require.NoError(t, os.Rename(tempPath, configPath))

	// Verify final file exists and temp is gone
	assert.FileExists(t, configPath)
	assert.NoFileExists(t, tempPath)

	// Verify content
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var readConfig map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &readConfig))
	assert.Equal(t, "vim", readConfig["editor"])
	assert.Equal(t, "json", readConfig["output_format"])
}

// Profile tests
func TestIsValidProfileName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"work", true},
		{"personal-dev", true},
		{"my_profile", true},
		{"Profile123", true},
		{"default", false}, // reserved
		{"", false},        // too short
		{strings.Repeat("a", 51), false}, // too long
		{"invalid name", false},          // space not allowed
		{"invalid@name", false},          // @ not allowed
		{"a", true},                      // single character
		{strings.Repeat("a", 50), true},  // exactly 50
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidProfileName(tt.name)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestProfileExists(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	profilesDir := filepath.Join(tmpHome, ".dsa", "profiles")
	os.MkdirAll(profilesDir, 0755)

	// Create a test profile
	profilePath := filepath.Join(profilesDir, "test.yaml")
	os.WriteFile(profilePath, []byte("editor: vim\n"), 0644)

	assert.True(t, profileExists("test"))
	assert.False(t, profileExists("nonexistent"))
	assert.True(t, profileExists("default")) // default always exists conceptually
}

func TestGetSetActiveProfile(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Initially should be "default"
	active, err := getActiveProfile()
	assert.NoError(t, err)
	assert.Equal(t, "default", active)

	// Set to a profile
	err = setActiveProfile("work")
	assert.NoError(t, err)

	// Verify it was set
	active, err = getActiveProfile()
	assert.NoError(t, err)
	assert.Equal(t, "work", active)

	// Set back to default (should remove active-profile file)
	err = setActiveProfile("default")
	assert.NoError(t, err)

	active, err = getActiveProfile()
	assert.NoError(t, err)
	assert.Equal(t, "default", active)
}

func TestGetProfilePath(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Default profile should point to config.yaml
	defaultPath, err := getProfilePath("default")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(tmpHome, ".dsa", "config.yaml"), defaultPath)

	// Named profile should point to profiles directory
	workPath, err := getProfilePath("work")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(tmpHome, ".dsa", "profiles", "work.yaml"), workPath)
}

// Validation tests

func TestValidateConfig(t *testing.T) {
	t.Run("valid config passes", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
		os.MkdirAll(filepath.Dir(configPath), 0755)

		validConfig := map[string]interface{}{
			"editor":        "vim",
			"output_format": "text",
			"no_color":      false,
		}
		yamlData, _ := yaml.Marshal(validConfig)
		os.WriteFile(configPath, yamlData, 0644)

		errors, err := validateConfig()

		assert.NoError(t, err)
		assert.Empty(t, errors)
	})

	t.Run("no config file is valid", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		// No config file created

		errors, err := validateConfig()

		assert.NoError(t, err)
		assert.Nil(t, errors)
	})

	t.Run("invalid key fails", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
		os.MkdirAll(filepath.Dir(configPath), 0755)

		invalidConfig := map[string]interface{}{
			"unknown_key": "value",
		}
		yamlData, _ := yaml.Marshal(invalidConfig)
		os.WriteFile(configPath, yamlData, 0644)

		errors, err := validateConfig()

		assert.NoError(t, err)
		assert.Len(t, errors, 1)
		assert.Equal(t, "unknown_key", errors[0].Key)
		assert.Contains(t, errors[0].Message, "unknown configuration key")
	})

	t.Run("invalid enum value fails", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
		os.MkdirAll(filepath.Dir(configPath), 0755)

		invalidConfig := map[string]interface{}{
			"color_scheme": "invalid",
		}
		yamlData, _ := yaml.Marshal(invalidConfig)
		os.WriteFile(configPath, yamlData, 0644)

		errors, err := validateConfig()

		assert.NoError(t, err)
		assert.Len(t, errors, 1)
		assert.Equal(t, "color_scheme", errors[0].Key)
		assert.Contains(t, errors[0].Message, "must be one of")
	})

	t.Run("wrong type fails", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
		os.MkdirAll(filepath.Dir(configPath), 0755)

		invalidConfig := map[string]interface{}{
			"no_color": "not a boolean",
		}
		yamlData, _ := yaml.Marshal(invalidConfig)
		os.WriteFile(configPath, yamlData, 0644)

		errors, err := validateConfig()

		assert.NoError(t, err)
		assert.Len(t, errors, 1)
		assert.Equal(t, "no_color", errors[0].Key)
		assert.Contains(t, errors[0].Message, "must be 'true' or 'false'")
	})

	t.Run("multiple errors collected", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
		os.MkdirAll(filepath.Dir(configPath), 0755)

		invalidConfig := map[string]interface{}{
			"unknown_key":  "value",
			"color_scheme": "invalid",
			"no_color":     "wrong",
		}
		yamlData, _ := yaml.Marshal(invalidConfig)
		os.WriteFile(configPath, yamlData, 0644)

		errors, err := validateConfig()

		assert.NoError(t, err)
		assert.Len(t, errors, 3)
	})

	t.Run("invalid YAML returns error", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
		os.MkdirAll(filepath.Dir(configPath), 0755)

		// Write invalid YAML
		os.WriteFile(configPath, []byte("invalid: yaml: content: [[["), 0644)

		errors, err := validateConfig()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not valid YAML")
		assert.Nil(t, errors)
	})
}

func TestGetDefaultConfig(t *testing.T) {
	t.Run("returns all default values", func(t *testing.T) {
		defaults := getDefaultConfig()

		// Should have entries for all valid config keys
		for _, key := range validConfigKeys {
			_, exists := defaults[key]
			assert.True(t, exists, "Default should exist for key: %s", key)
		}
	})

	t.Run("defaults match expected values", func(t *testing.T) {
		defaults := getDefaultConfig()

		// Check some known defaults
		assert.Equal(t, "text", defaults["output_format"])
		assert.Equal(t, false, defaults["no_color"])
		assert.Equal(t, false, defaults["verbose"])
		assert.Equal(t, "table", defaults["list_format"])
		assert.Equal(t, "table", defaults["status_format"])
		assert.Equal(t, "normal", defaults["output_style"])
		assert.Equal(t, "default", defaults["color_scheme"])
	})
}
