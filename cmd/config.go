package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Valid configuration keys
var validConfigKeys = []string{
	"editor",
	"editor_args",
	"output_format",
	"no_color",
	"database_path",
	"verbose",
	"list_format",
	"status_format",
	"output_style",
	"color_scheme",
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `View and modify DSA CLI configuration.

Configuration precedence (highest to lowest):
  1. Command-line flags
  2. Environment variables (DSA_*)
  3. Project config (./.dsa/config.yaml)
  4. Global config (~/.dsa/config.yaml)
  5. Default values

Editor Integration:
  Configure your preferred editor and arguments with placeholders:
    dsa config set editor code
    dsa config set editor_args "--goto {file}:{line}"

  Supported placeholders: {file}, {line}, {column}

  Popular editor examples:
    Vim:     dsa config set editor_args "+{line}"
    Neovim:  dsa config set editor_args "-c 'normal {line}G'"
    VS Code: dsa config set editor_args "--goto {file}:{line}"
    Emacs:   dsa config set editor_args "+{line}:{column}"`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration value",
	Long: `Display the current value of a configuration setting.

If no key is specified, displays all settings.

Examples:
  dsa config get editor
  dsa config get output_format
  dsa config get`,
	Args: cobra.MaximumNArgs(1),
	Run:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set configuration value",
	Long: `Update a configuration setting in the global config file.

Valid keys: editor, editor_args, output_format, no_color, database_path, verbose,
            list_format, status_format, output_style, color_scheme

Editor Integration:
  editor       - Editor command (e.g., vim, code, nvim, emacs)
  editor_args  - Arguments with placeholders: {file}, {line}, {column}

Examples:
  dsa config set editor vim
  dsa config set editor_args "+{line}"
  dsa config set editor code
  dsa config set editor_args "--goto {file}:{line}"
  dsa config set output_format json
  dsa config set no_color true

Editor Fallback:
  If editor is not configured, the CLI will:
  1. Check $EDITOR environment variable
  2. Use OS-specific default (open on macOS, xdg-open on Linux, start on Windows)`,
	Args: cobra.ExactArgs(2),
	Run:  runConfigSet,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration settings",
	Long: `Display all configuration settings with their current values and sources.

Shows which settings are using defaults vs custom values.

Example:
  dsa config list`,
	Args: cobra.NoArgs,
	Run:  runConfigList,
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Remove configuration setting",
	Long: `Remove a setting from the global config file, reverting to default.

Examples:
  dsa config unset editor
  dsa config unset output_format`,
	Args: cobra.ExactArgs(1),
	Run:  runConfigUnset,
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long: `Check configuration file for invalid values and type errors.

Validates all configuration keys for:
  - Valid key names
  - Correct data types
  - Valid enum values
  - Proper value constraints

If validation passes, displays success message.
If errors found, shows detailed error messages with suggestions.

Examples:
  dsa config validate`,
	Args: cobra.NoArgs,
	Run:  runConfigValidate,
}

var configResetCmd = &cobra.Command{
	Use:   "reset [key]",
	Short: "Reset configuration to defaults",
	Long: `Reset all configuration or a specific key to default values.

When resetting all configuration:
  - Prompts for confirmation
  - Creates a backup with timestamp
  - Restores all settings to factory defaults

When resetting a specific key:
  - No confirmation required
  - Only resets the specified key
  - Other settings remain unchanged

Examples:
  dsa config reset              # Reset all (with confirmation)
  dsa config reset editor       # Reset only editor key
  dsa config reset color_scheme # Reset only color_scheme key`,
	Args: cobra.MaximumNArgs(1),
	Run:  runConfigReset,
}

var configDefaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Display default configuration values",
	Long: `Display all default configuration values in YAML format.

Use this to see what values will be used if no configuration is set.
Output can be piped or redirected to create a config file template.

Examples:
  dsa config defaults                    # View defaults
  dsa config defaults > ~/.dsa/config.yaml # Create config from defaults`,
	Args: cobra.NoArgs,
	Run:  runConfigDefaults,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configUnsetCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configDefaultsCmd)
}

// runConfigGet displays configuration values
func runConfigGet(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		// Show all settings
		allSettings := viper.AllSettings()
		keys := make([]string, 0, len(allSettings))
		for key := range allSettings {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			fmt.Printf("%s: %v\n", key, allSettings[key])
		}
		return
	}

	key := args[0]
	if !isValidConfigKey(key) {
		fmt.Fprintf(os.Stderr, "Error: Unknown configuration key '%s'\n", key)
		fmt.Fprintf(os.Stderr, "Valid keys: %s\n", strings.Join(validConfigKeys, ", "))
		os.Exit(2)
	}

	value := viper.Get(key)
	fmt.Printf("%s: %v\n", key, value)
}

// runConfigSet updates a configuration value
func runConfigSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	// Validate key
	if !isValidConfigKey(key) {
		fmt.Fprintf(os.Stderr, "Error: Unknown configuration key '%s'\n", key)
		fmt.Fprintf(os.Stderr, "Valid keys: %s\n", strings.Join(validConfigKeys, ", "))
		os.Exit(2)
	}

	// Validate and convert value based on key type
	parsedValue, err := parseConfigValue(key, value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// Get global config file path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(home, ".dsa", "config.yaml")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create config directory: %v\n", err)
		os.Exit(1)
	}

	// Read existing config or create empty map
	configData := make(map[string]interface{})
	if data, err := os.ReadFile(configPath); err == nil {
		yaml.Unmarshal(data, &configData)
	}

	// Update value
	configData[key] = parsedValue

	// Write atomically (temp file + rename)
	tempPath := configPath + ".tmp"
	yamlData, err := yaml.Marshal(configData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to marshal config: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(tempPath, yamlData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write config: %v\n", err)
		os.Exit(1)
	}

	if err := os.Rename(tempPath, configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to save config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Configuration updated: %s = %v\n", key, parsedValue)
}

// runConfigList displays all configuration settings with sources
func runConfigList(cmd *cobra.Command, args []string) {
	keys := make([]string, len(validConfigKeys))
	copy(keys, validConfigKeys)
	sort.Strings(keys)

	// Print table header
	fmt.Printf("%-20s %-30s %-20s\n", "KEY", "VALUE", "SOURCE")
	fmt.Printf("%-20s %-30s %-20s\n", strings.Repeat("-", 20), strings.Repeat("-", 30), strings.Repeat("-", 20))

	for _, key := range keys {
		value := viper.Get(key)
		source := getConfigSource(key)
		fmt.Printf("%-20s %-30v %-20s\n", key, value, source)
	}
}

// runConfigUnset removes a configuration setting
func runConfigUnset(cmd *cobra.Command, args []string) {
	key := args[0]

	if !isValidConfigKey(key) {
		fmt.Fprintf(os.Stderr, "Error: Unknown configuration key '%s'\n", key)
		fmt.Fprintf(os.Stderr, "Valid keys: %s\n", strings.Join(validConfigKeys, ", "))
		os.Exit(2)
	}

	// Get global config file path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(home, ".dsa", "config.yaml")

	// Read existing config
	configData := make(map[string]interface{})
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("✓ Configuration key '%s' is not set\n", key)
			return
		}
		fmt.Fprintf(os.Stderr, "Error: Failed to read config: %v\n", err)
		os.Exit(1)
	}

	yaml.Unmarshal(data, &configData)

	// Remove key
	if _, exists := configData[key]; !exists {
		fmt.Printf("✓ Configuration key '%s' is not set\n", key)
		return
	}

	delete(configData, key)

	// Write back atomically
	tempPath := configPath + ".tmp"
	yamlData, _ := yaml.Marshal(configData)
	if err := os.WriteFile(tempPath, yamlData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write config: %v\n", err)
		os.Exit(1)
	}

	if err := os.Rename(tempPath, configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to save config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Configuration key '%s' removed (reverted to default)\n", key)
}

// runConfigValidate validates the configuration file
func runConfigValidate(cmd *cobra.Command, args []string) {
	errors, err := validateConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	if len(errors) == 0 {
		fmt.Println("✓ Configuration is valid")
		return
	}

	// Print all errors
	fmt.Fprintf(os.Stderr, "Configuration validation failed:\n\n")
	for _, e := range errors {
		fmt.Fprintf(os.Stderr, "  Error in '%s' = %v\n", e.Key, e.Value)
		fmt.Fprintf(os.Stderr, "    %s\n\n", e.Message)
	}

	fmt.Fprintf(os.Stderr, "Hint: Run 'dsa config reset' to restore defaults\n")
	os.Exit(2)
}

// runConfigReset resets configuration to defaults
func runConfigReset(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(home, ".dsa", "config.yaml")

	if len(args) == 0 {
		// Full reset with confirmation
		fmt.Print("Reset all configuration to defaults? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Reset cancelled")
			return
		}

		// Create backup with timestamp
		timestamp := time.Now().Format("20060102-150405")
		backupPath := fmt.Sprintf("%s.backup.%s", configPath, timestamp)

		if _, err := os.Stat(configPath); err == nil {
			data, _ := os.ReadFile(configPath)
			os.WriteFile(backupPath, data, 0644)
		}

		// Write defaults
		defaults := getDefaultConfig()
		yamlData, _ := yaml.Marshal(defaults)

		// Ensure directory exists
		os.MkdirAll(filepath.Dir(configPath), 0755)

		tempPath := configPath + ".tmp"
		if err := os.WriteFile(tempPath, yamlData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to write config: %v\n", err)
			os.Exit(1)
		}

		if err := os.Rename(tempPath, configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save config: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(backupPath); err == nil {
			fmt.Printf("✓ Configuration reset to defaults. Backup saved to %s\n", backupPath)
		} else {
			fmt.Println("✓ Configuration reset to defaults")
		}
	} else {
		// Single key reset
		key := args[0]

		if !isValidConfigKey(key) {
			fmt.Fprintf(os.Stderr, "Error: Unknown configuration key '%s'\n", key)
			fmt.Fprintf(os.Stderr, "Valid keys: %s\n", strings.Join(validConfigKeys, ", "))
			os.Exit(2)
		}

		// Get default value for this key
		defaults := getDefaultConfig()
		defaultValue, exists := defaults[key]
		if !exists {
			fmt.Fprintf(os.Stderr, "Error: No default value found for '%s'\n", key)
			os.Exit(1)
		}

		// Load current config
		configData := make(map[string]interface{})
		if data, err := os.ReadFile(configPath); err == nil {
			yaml.Unmarshal(data, &configData)
		}

		// Set to default
		configData[key] = defaultValue

		// Ensure directory exists
		os.MkdirAll(filepath.Dir(configPath), 0755)

		// Write back atomically
		yamlData, _ := yaml.Marshal(configData)
		tempPath := configPath + ".tmp"
		if err := os.WriteFile(tempPath, yamlData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to write config: %v\n", err)
			os.Exit(1)
		}

		if err := os.Rename(tempPath, configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ '%s' reset to default value: %v\n", key, defaultValue)
	}
}

// runConfigDefaults displays default configuration values
func runConfigDefaults(cmd *cobra.Command, args []string) {
	defaults := getDefaultConfig()

	yamlData, err := yaml.Marshal(defaults)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(yamlData))
}

// Helper functions

func isValidConfigKey(key string) bool {
	for _, valid := range validConfigKeys {
		if key == valid {
			return true
		}
	}
	return false
}

func parseConfigValue(key, value string) (interface{}, error) {
	switch key {
	case "no_color", "verbose":
		// Boolean values
		if value == "true" || value == "false" {
			return value == "true", nil
		}
		return nil, fmt.Errorf("'%s' must be 'true' or 'false'", key)
	case "output_format":
		// Enum values
		if value != "text" && value != "json" {
			return nil, fmt.Errorf("output_format must be 'text' or 'json'")
		}
		return value, nil
	case "list_format", "status_format":
		// Enum values for format options
		if value != "table" && value != "json" {
			return nil, fmt.Errorf("%s must be 'table' or 'json'", key)
		}
		return value, nil
	case "output_style":
		// Enum values for output style
		validStyles := []string{"normal", "compact", "verbose"}
		for _, valid := range validStyles {
			if value == valid {
				return value, nil
			}
		}
		return nil, fmt.Errorf("output_style must be one of: %s", strings.Join(validStyles, ", "))
	case "color_scheme":
		// Enum values for color schemes
		validSchemes := []string{"default", "solarized", "monokai", "nord"}
		for _, valid := range validSchemes {
			if value == valid {
				return value, nil
			}
		}
		return nil, fmt.Errorf("color_scheme must be one of: %s", strings.Join(validSchemes, ", "))
	case "editor", "editor_args", "database_path":
		// String values
		return value, nil
	default:
		return nil, fmt.Errorf("unknown key: %s", key)
	}
}

func getConfigSource(key string) string {
	// Check if set via environment variable
	envKey := "DSA_" + strings.ToUpper(key)
	if os.Getenv(envKey) != "" {
		return "environment"
	}

	// Check if in project config
	projectConfigPath := filepath.Join(".dsa", "config.yaml")
	if data, err := os.ReadFile(projectConfigPath); err == nil {
		var projectConfig map[string]interface{}
		if yaml.Unmarshal(data, &projectConfig) == nil {
			if _, exists := projectConfig[key]; exists {
				return "project"
			}
		}
	}

	// Check if in global config
	home, _ := os.UserHomeDir()
	if home != "" {
		globalConfigPath := filepath.Join(home, ".dsa", "config.yaml")
		if data, err := os.ReadFile(globalConfigPath); err == nil {
			var globalConfig map[string]interface{}
			if yaml.Unmarshal(data, &globalConfig) == nil {
				if _, exists := globalConfig[key]; exists {
					return "global"
				}
			}
		}
	}

	return "default"
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Key     string
	Value   interface{}
	Message string
}

// validateConfig validates the configuration file
func validateConfig() ([]ValidationError, error) {
	errors := []ValidationError{}

	// Read config file
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".dsa", "config.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No config file is valid (uses defaults)
		}
		return nil, err
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("config file is not valid YAML: %v", err)
	}

	// Validate each key
	for key, value := range config {
		// Check if key is valid
		if !isValidConfigKey(key) {
			errors = append(errors, ValidationError{
				Key:     key,
				Value:   value,
				Message: fmt.Sprintf("unknown configuration key (valid keys: %s)", strings.Join(validConfigKeys, ", ")),
			})
			continue
		}

		// Validate value using parseConfigValue
		if _, err := parseConfigValue(key, fmt.Sprintf("%v", value)); err != nil {
			errors = append(errors, ValidationError{
				Key:     key,
				Value:   value,
				Message: err.Error(),
			})
		}
	}

	return errors, nil
}

// getDefaultConfig returns all default configuration values
func getDefaultConfig() map[string]interface{} {
	defaults := make(map[string]interface{})

	// Set default values (mirroring internal/config/config.go InitConfig)
	defaultEditor := os.Getenv("EDITOR")
	if defaultEditor == "" {
		defaultEditor = "vim"
	}

	defaults["editor"] = defaultEditor
	defaults["editor_args"] = ""
	defaults["output_format"] = "text"
	defaults["no_color"] = false
	defaults["verbose"] = false
	defaults["list_format"] = "table"
	defaults["status_format"] = "table"
	defaults["output_style"] = "normal"
	defaults["color_scheme"] = "default"

	// database_path default
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		defaults["database_path"] = filepath.Join(home, ".dsa", "dsa.db")
	} else {
		defaults["database_path"] = ""
	}

	return defaults
}

// Profile commands
var configProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long: `Create, list, switch, and manage configuration profiles.

Profiles allow you to save and switch between different configuration sets.
Each profile is a separate YAML file stored in ~/.dsa/profiles/.

Examples:
  dsa config profile create work
  dsa config profile list
  dsa config profile switch work
  dsa config profile delete old-config`,
}

var configProfileCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create new configuration profile",
	Long: `Create a new configuration profile as a copy of the current config.

The profile name must be alphanumeric with hyphens or underscores.
The name "default" is reserved and cannot be used.

Examples:
  dsa config profile create work
  dsa config profile create personal-dev`,
	Args: cobra.ExactArgs(1),
	Run:  runConfigProfileCreate,
}

var configProfileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration profiles",
	Long: `Display all available configuration profiles.

The currently active profile is marked with an asterisk (*).

Example:
  dsa config profile list`,
	Args: cobra.NoArgs,
	Run:  runConfigProfileList,
}

var configProfileSwitchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to a different profile",
	Long: `Switch to a different configuration profile.

All subsequent commands will use the settings from the selected profile.
Use "default" to switch back to the default configuration.

Examples:
  dsa config profile switch work
  dsa config profile switch default`,
	Args: cobra.ExactArgs(1),
	Run:  runConfigProfileSwitch,
}

var configProfileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a configuration profile",
	Long: `Delete a configuration profile with confirmation.

The "default" profile cannot be deleted.
If you delete the active profile, it will switch to "default" automatically.

Examples:
  dsa config profile delete work
  dsa config profile delete old-config`,
	Args: cobra.ExactArgs(1),
	Run:  runConfigProfileDelete,
}

var configProfileExportCmd = &cobra.Command{
	Use:   "export <name>",
	Short: "Export profile to file",
	Long: `Export a configuration profile to a file.

The exported file can be shared with other users or machines.

Examples:
  dsa config profile export work --output work-config.yaml
  dsa config profile export personal -o personal.yaml`,
	Args: cobra.ExactArgs(1),
	Run:  runConfigProfileExport,
}

var configProfileImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import profile from file",
	Long: `Import a configuration profile from a file.

Creates a new profile with the specified name from the import file.

Examples:
  dsa config profile import --file work-config.yaml --name work
  dsa config profile import -f config.yaml -n imported`,
	Args: cobra.NoArgs,
	Run:  runConfigProfileImport,
}

func init() {
	// Add profile command to config
	configCmd.AddCommand(configProfileCmd)
	
	// Add all profile subcommands
	configProfileCmd.AddCommand(configProfileCreateCmd)
	configProfileCmd.AddCommand(configProfileListCmd)
	configProfileCmd.AddCommand(configProfileSwitchCmd)
	configProfileCmd.AddCommand(configProfileDeleteCmd)
	configProfileCmd.AddCommand(configProfileExportCmd)
	configProfileCmd.AddCommand(configProfileImportCmd)
	
	// Flags for export
	configProfileExportCmd.Flags().StringP("output", "o", "", "Output file path")
	configProfileExportCmd.MarkFlagRequired("output")
	
	// Flags for import
	configProfileImportCmd.Flags().StringP("file", "f", "", "Input file path")
	configProfileImportCmd.Flags().StringP("name", "n", "", "Profile name")
	configProfileImportCmd.MarkFlagRequired("file")
	configProfileImportCmd.MarkFlagRequired("name")
}

// Profile command implementations
func runConfigProfileCreate(cmd *cobra.Command, args []string) {
	profileName := args[0]

	// Validate profile name
	if !isValidProfileName(profileName) {
		if profileName == "default" {
			fmt.Fprintf(os.Stderr, "Error: 'default' is a reserved profile name\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error: Invalid profile name '%s'\n", profileName)
			fmt.Fprintf(os.Stderr, "Profile name must be 1-50 alphanumeric characters, hyphens, or underscores\n")
		}
		os.Exit(2)
	}

	// Check if profile already exists
	if profileExists(profileName) {
		fmt.Fprintf(os.Stderr, "Error: Profile '%s' already exists\n", profileName)
		os.Exit(2)
	}

	// Get current config file path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	currentConfigPath := filepath.Join(home, ".dsa", "config.yaml")

	// Read current config
	configData, err := os.ReadFile(currentConfigPath)
	if err != nil {
		// If config doesn't exist, create empty profile
		if os.IsNotExist(err) {
			configData = []byte("{}\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error: Failed to read current config: %v\n", err)
			os.Exit(1)
		}
	}

	// Ensure profiles directory exists
	profilesDir, err := getProfilesDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get profiles directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create profiles directory: %v\n", err)
		os.Exit(1)
	}

	// Write to new profile file (atomic write)
	profilePath, _ := getProfilePath(profileName)
	tempPath := profilePath + ".tmp"

	if err := os.WriteFile(tempPath, configData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write profile: %v\n", err)
		os.Exit(1)
	}

	if err := os.Rename(tempPath, profilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to save profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Profile '%s' created\n", profileName)
}

func runConfigProfileList(cmd *cobra.Command, args []string) {
	// Get active profile
	activeProfile, err := getActiveProfile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get active profile: %v\n", err)
		os.Exit(1)
	}

	// Get profiles directory
	profilesDir, err := getProfilesDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get profiles directory: %v\n", err)
		os.Exit(1)
	}

	// Collect profile names
	profiles := []string{"default"}

	// Scan profiles directory if it exists
	entries, err := os.ReadDir(profilesDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
				profileName := strings.TrimSuffix(entry.Name(), ".yaml")
				profiles = append(profiles, profileName)
			}
		}
	}

	// Sort profiles alphabetically (but keep default first)
	sort.Slice(profiles[1:], func(i, j int) bool {
		return profiles[1+i] < profiles[1+j]
	})

	// Display profiles
	if len(profiles) == 1 {
		fmt.Println("Available profiles:")
		if activeProfile == "default" {
			fmt.Println("  default (active) *")
		} else {
			fmt.Println("  default")
		}
	} else {
		fmt.Println("Available profiles:")
		for _, profile := range profiles {
			if profile == activeProfile {
				fmt.Printf("  %s (active) *\n", profile)
			} else {
				fmt.Printf("  %s\n", profile)
			}
		}
	}
}

func runConfigProfileSwitch(cmd *cobra.Command, args []string) {
	profileName := args[0]

	// Validate profile exists (or is "default")
	if profileName != "default" && !profileExists(profileName) {
		fmt.Fprintf(os.Stderr, "Error: Profile '%s' does not exist\n", profileName)
		fmt.Fprintf(os.Stderr, "Run 'dsa config profile list' to see available profiles\n")
		os.Exit(2)
	}

	// Set active profile
	if err := setActiveProfile(profileName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to switch profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Switched to profile '%s'\n", profileName)
}

func runConfigProfileDelete(cmd *cobra.Command, args []string) {
	profileName := args[0]

	// Cannot delete "default" profile
	if profileName == "default" {
		fmt.Fprintf(os.Stderr, "Error: Cannot delete 'default' profile\n")
		os.Exit(2)
	}

	// Validate profile exists
	if !profileExists(profileName) {
		fmt.Fprintf(os.Stderr, "Error: Profile '%s' does not exist\n", profileName)
		os.Exit(2)
	}

	// Confirmation prompt
	fmt.Printf("Delete profile '%s'? [y/N]: ", profileName)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read input: %v\n", err)
		os.Exit(1)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	// Get profile path
	profilePath, err := getProfilePath(profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get profile path: %v\n", err)
		os.Exit(1)
	}

	// Delete profile file
	if err := os.Remove(profilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to delete profile: %v\n", err)
		os.Exit(1)
	}

	// If this was the active profile, switch to default
	activeProfile, err := getActiveProfile()
	if err == nil && activeProfile == profileName {
		setActiveProfile("default")
		fmt.Printf("✓ Profile '%s' deleted (switched to 'default')\n", profileName)
	} else {
		fmt.Printf("✓ Profile '%s' deleted\n", profileName)
	}
}

func runConfigProfileExport(cmd *cobra.Command, args []string) {
	profileName := args[0]
	outputPath, _ := cmd.Flags().GetString("output")

	// Validate profile exists
	if !profileExists(profileName) {
		fmt.Fprintf(os.Stderr, "Error: Profile '%s' does not exist\n", profileName)
		os.Exit(2)
	}

	// Get profile path
	profilePath, err := getProfilePath(profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get profile path: %v\n", err)
		os.Exit(1)
	}

	// Read profile data
	profileData, err := os.ReadFile(profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read profile: %v\n", err)
		os.Exit(1)
	}

	// Write to output file (atomic write)
	tempPath := outputPath + ".tmp"
	if err := os.WriteFile(tempPath, profileData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write export file: %v\n", err)
		os.Exit(1)
	}

	if err := os.Rename(tempPath, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to save export file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Profile '%s' exported to %s\n", profileName, outputPath)
}

func runConfigProfileImport(cmd *cobra.Command, args []string) {
	inputPath, _ := cmd.Flags().GetString("file")
	profileName, _ := cmd.Flags().GetString("name")

	// Validate profile name
	if !isValidProfileName(profileName) {
		if profileName == "default" {
			fmt.Fprintf(os.Stderr, "Error: 'default' is a reserved profile name\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error: Invalid profile name '%s'\n", profileName)
			fmt.Fprintf(os.Stderr, "Profile name must be 1-50 alphanumeric characters, hyphens, or underscores\n")
		}
		os.Exit(2)
	}

	// Check if profile already exists
	if profileExists(profileName) {
		fmt.Fprintf(os.Stderr, "Error: Profile '%s' already exists\n", profileName)
		os.Exit(2)
	}

	// Read input file
	profileData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read input file: %v\n", err)
		os.Exit(1)
	}

	// Validate it's valid YAML
	var testParse map[string]interface{}
	if err := yaml.Unmarshal(profileData, &testParse); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Input file is not valid YAML: %v\n", err)
		os.Exit(2)
	}

	// Ensure profiles directory exists
	profilesDir, err := getProfilesDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get profiles directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create profiles directory: %v\n", err)
		os.Exit(1)
	}

	// Write to profile file (atomic write)
	profilePath, _ := getProfilePath(profileName)
	tempPath := profilePath + ".tmp"

	if err := os.WriteFile(tempPath, profileData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write profile: %v\n", err)
		os.Exit(1)
	}

	if err := os.Rename(tempPath, profilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to save profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Profile '%s' imported from %s\n", profileName, inputPath)
}

// Profile storage and management functions

// isValidProfileName validates profile name
func isValidProfileName(name string) bool {
	// "default" is reserved
	if name == "default" {
		return false
	}
	
	// Length: 1-50 characters
	if len(name) == 0 || len(name) > 50 {
		return false
	}
	
	// Alphanumeric, hyphens, underscores only
	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || 
			(ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return false
		}
	}
	
	return true
}

// getProfilesDir returns the profiles directory path
func getProfilesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".dsa", "profiles"), nil
}

// getActiveProfilePath returns the active profile tracker file path
func getActiveProfilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".dsa", "active-profile"), nil
}

// getActiveProfile returns the currently active profile name
func getActiveProfile() (string, error) {
	activeProfilePath, err := getActiveProfilePath()
	if err != nil {
		return "", err
	}
	
	data, err := os.ReadFile(activeProfilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "default", nil
		}
		return "", err
	}
	
	return strings.TrimSpace(string(data)), nil
}

// setActiveProfile sets the active profile
func setActiveProfile(profileName string) error {
	activeProfilePath, err := getActiveProfilePath()
	if err != nil {
		return err
	}
	
	// If switching to default, remove active-profile file
	if profileName == "default" {
		os.Remove(activeProfilePath)
		return nil
	}
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(activeProfilePath), 0755); err != nil {
		return err
	}
	
	// Write profile name to active-profile file (atomic write)
	tempPath := activeProfilePath + ".tmp"
	if err := os.WriteFile(tempPath, []byte(profileName), 0644); err != nil {
		return err
	}
	return os.Rename(tempPath, activeProfilePath)
}

// getProfilePath returns the file path for a named profile
func getProfilePath(profileName string) (string, error) {
	if profileName == "default" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".dsa", "config.yaml"), nil
	}
	
	profilesDir, err := getProfilesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(profilesDir, profileName+".yaml"), nil
}

// profileExists checks if a profile exists
func profileExists(profileName string) bool {
	// "default" profile always exists conceptually (uses config.yaml)
	if profileName == "default" {
		return true
	}

	profilePath, err := getProfilePath(profileName)
	if err != nil {
		return false
	}

	_, err = os.Stat(profilePath)
	return err == nil
}
