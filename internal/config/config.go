package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration settings
type Config struct {
	Editor       string
	OutputFormat string
	NoColor      bool
	DatabasePath string
	Verbose      bool
}

var globalConfig *Config

// InitConfig initializes configuration from files, env vars, and defaults
func InitConfig() error {
	home, err := os.UserHomeDir()

	// Set defaults first (lowest priority)
	defaultEditor := os.Getenv("EDITOR")
	if defaultEditor == "" {
		defaultEditor = "vim"
	}

	viper.SetDefault("editor", defaultEditor)
	viper.SetDefault("editor_args", "")
	viper.SetDefault("output_format", "text")
	viper.SetDefault("no_color", false)
	if home != "" && err == nil {
		viper.SetDefault("database_path", filepath.Join(home, ".dsa", "dsa.db"))
	}
	viper.SetDefault("verbose", false)
	viper.SetDefault("list_format", "table")
	viper.SetDefault("status_format", "table")
	viper.SetDefault("output_style", "normal")
	viper.SetDefault("color_scheme", "default")

	// Check for active profile
	var activeConfigFile string
	if home != "" && err == nil {
		activeProfilePath := filepath.Join(home, ".dsa", "active-profile")
		if activeProfileData, err := os.ReadFile(activeProfilePath); err == nil {
			// Active profile exists - use it
			profileName := strings.TrimSpace(string(activeProfileData))
			activeConfigFile = filepath.Join(home, ".dsa", "profiles", profileName+".yaml")
		} else {
			// No active profile - use default config.yaml
			activeConfigFile = filepath.Join(home, ".dsa", "config.yaml")
		}

		// Read the active config file (if it exists)
		if _, err := os.Stat(activeConfigFile); err == nil {
			viper.SetConfigFile(activeConfigFile)
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	// Read project config and merge (project config overrides global/profile)
	projectConfigFile := filepath.Join(".dsa", "config.yaml")
	if _, err := os.Stat(projectConfigFile); err == nil {
		viper.SetConfigFile(projectConfigFile)
		if err := viper.MergeInConfig(); err != nil {
			return fmt.Errorf("failed to read project config file: %w", err)
		}
	}

	// Bind environment variables (override config files)
	viper.SetEnvPrefix("DSA")
	viper.AutomaticEnv()

	// Validate configuration
	if err := validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Load into global config
	globalConfig = &Config{
		Editor:       viper.GetString("editor"),
		OutputFormat: viper.GetString("output_format"),
		NoColor:      viper.GetBool("no_color"),
		DatabasePath: viper.GetString("database_path"),
		Verbose:      viper.GetBool("verbose"),
	}

	return nil
}

// Get returns the global config instance
func Get() *Config {
	if globalConfig == nil {
		// Initialize with defaults if not already loaded
		if err := InitConfig(); err != nil {
			// Return defaults if init fails
			return &Config{
				Editor:       "vim",
				OutputFormat: "text",
				NoColor:      false,
				DatabasePath: filepath.Join(os.Getenv("HOME"), ".dsa", "dsa.db"),
				Verbose:      false,
			}
		}
	}
	return globalConfig
}

// validateConfig validates configuration values
func validateConfig() error {
	// Validate output_format
	format := viper.GetString("output_format")
	if format != "text" && format != "json" {
		return fmt.Errorf("output_format must be 'text' or 'json', got '%s'", format)
	}

	// Validate database_path is writable
	dbPath := viper.GetString("database_path")
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		// Directory doesn't exist - try to create it
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("database directory is not writable: %s", dbDir)
		}
	}

	return nil
}

// GetString returns a config string value with flag override
func GetString(key string, flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return viper.GetString(key)
}

// GetBool returns a config bool value with flag override
func GetBool(key string, flagValue bool, flagSet bool) bool {
	if flagSet {
		return flagValue
	}
	return viper.GetBool(key)
}
