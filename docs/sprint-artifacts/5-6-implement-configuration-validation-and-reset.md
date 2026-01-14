# Story 5.6: Implement Configuration Validation and Reset

Status: review

## Story

As a **user**,
I want **to validate my configuration and reset to defaults if needed**,
So that **I can fix configuration issues easily** (FR28, FR29).

## Acceptance Criteria

### AC1: Validate Configuration

**Given** I have modified my configuration
**When** I run `dsa config validate`
**Then** The CLI checks all config values for validity:
  - Data types (string, int, bool)
  - Value ranges (e.g., test_timeout > 0)
  - Valid enum values (e.g., color_scheme in [default, solarized, ...])
**And** If valid, I see: "✓ Configuration is valid"
**And** If invalid, I see specific errors with suggestions for fixes

### AC2: Validate with Detailed Errors

**Given** My configuration is corrupted
**When** I run `dsa config validate`
**Then** I see detailed errors:
  - "Error: 'test_timeout' must be a positive integer, got: 'invalid'"
  - "Error: 'color_scheme' must be one of: default, solarized, monokai, nord. Got: 'unknown'"
**And** I see suggestion: "Run 'dsa config reset' to restore defaults"

### AC3: Reset All Configuration

**Given** I want to reset all configuration
**When** I run `dsa config reset`
**Then** The CLI asks for confirmation: "Reset all configuration to defaults? [y/N]"
**And** If confirmed, a backup is created at `~/.dsa/config.yaml.backup.<timestamp>`
**And** The config file is replaced with factory defaults
**And** I see: "✓ Configuration reset to defaults. Backup saved to <path>"

### AC4: Reset Specific Key

**Given** I want to reset a specific key
**When** I run `dsa config reset editor`
**Then** Only the "editor" key is reset to its default value
**And** Other configuration values are preserved
**And** I see: "✓ 'editor' reset to default value: vim"

### AC5: View Default Configuration

**Given** I want to view the default configuration
**When** I run `dsa config defaults`
**Then** I see all default values in YAML format
**And** I can use this as a reference for resetting specific values

## Tasks / Subtasks

- [x] **Task 1: Add Config Validation Subcommand** (AC: 1, 2)
  - [x] Create `dsa config validate` subcommand
  - [x] Add help text explaining what validation checks
  - [x] Set up for no arguments (validates entire config)

- [x] **Task 2: Implement Validation Logic** (AC: 1, 2)
  - [x] Load current configuration from file
  - [x] For each key in config file, validate:
    - Key is in validConfigKeys
    - Value matches expected type (string, bool, int)
    - Enum values are in allowed list
    - Numeric values meet constraints (if any)
  - [x] Collect all validation errors (don't stop at first error)
  - [x] Return validation result with detailed error messages

- [x] **Task 3: Implement Validation Output** (AC: 1, 2)
  - [x] If all valid, print success message: "✓ Configuration is valid"
  - [x] If errors found, print each error with context:
    - Key name
    - Current value
    - Expected type/values
    - Suggestion for fix
  - [x] Add hint to run `dsa config reset` if many errors
  - [x] Exit with code 2 if validation fails

- [x] **Task 4: Add Config Reset Subcommand** (AC: 3, 4)
  - [x] Create `dsa config reset` subcommand
  - [x] Accept optional [key] argument for single-key reset
  - [x] Add help text with examples
  - [x] Implement confirmation prompt for full reset

- [x] **Task 5: Implement Full Config Reset** (AC: 3)
  - [x] Show confirmation prompt: "Reset all configuration to defaults? [y/N]"
  - [x] If not confirmed, exit with message
  - [x] If confirmed:
    - Create backup: config.yaml.backup.<timestamp>
    - Generate default config from setDefaults() function
    - Write to config.yaml with atomic pattern
    - Show success message with backup path

- [x] **Task 6: Implement Single Key Reset** (AC: 4)
  - [x] Validate key is in validConfigKeys
  - [x] Get default value for key from Viper defaults
  - [x] Load current config file
  - [x] Update only the specified key to default value
  - [x] Write back to file with atomic pattern
  - [x] Show success message with key and default value

- [x] **Task 7: Add Config Defaults Subcommand** (AC: 5)
  - [x] Create `dsa config defaults` subcommand
  - [x] Collect all default values from Viper defaults
  - [x] Format as YAML output
  - [x] Print to stdout (can be piped/redirected)

- [x] **Task 8: Add Unit Tests** (AC: All)
  - [x] Test validation with valid config
  - [x] Test validation with invalid types
  - [x] Test validation with invalid enum values
  - [x] Test validation error message format
  - [x] Test full reset creates backup
  - [x] Test single key reset preserves other keys
  - [x] Test defaults output format

- [x] **Task 9: Add Integration Tests** (AC: All)
  - [x] Test validate on valid config file
  - [x] Test validate on corrupted config file
  - [x] Test full reset workflow with confirmation
  - [x] Test single key reset updates only one key
  - [x] Test backup file created with timestamp
  - [x] Test defaults output matches known defaults

## Dev Notes

### Architecture Patterns and Constraints

**Configuration Validation System:**
- **Validation Rules:** Based on existing parseConfigValue() logic
- **Error Collection:** Gather all errors before reporting (don't fail fast)
- **Helpful Messages:** Include current value, expected format, and fix suggestion
- **Exit Codes:** 0 for valid, 2 for validation errors

**Reset Patterns:**
- **Backup Format:** `config.yaml.backup.<timestamp>` (e.g., config.yaml.backup.20250129-143052)
- **Atomic Writes:** Use temp + rename pattern for both backup and new config
- **Confirmation:** Required for full reset, not for single-key reset
- **Defaults Source:** Use Viper's default values (already defined in setDefaults())

**Cobra Subcommand Structure:**
```go
// cmd/config.go - Add validation and reset subcommands
var configValidateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate configuration file",
    Long:  `Check configuration file for invalid values and type errors.`,
    Args:  cobra.NoArgs,
    Run:   runConfigValidate,
}

var configResetCmd = &cobra.Command{
    Use:   "reset [key]",
    Short: "Reset configuration to defaults",
    Long: `Reset all configuration or a specific key to default values.

Creates a backup before resetting entire configuration.

Examples:
  dsa config reset              # Reset all (with confirmation)
  dsa config reset editor       # Reset only editor key`,
    Args: cobra.MaximumNArgs(1),
    Run:  runConfigReset,
}

var configDefaultsCmd = &cobra.Command{
    Use:   "defaults",
    Short: "Display default configuration values",
    Args:  cobra.NoArgs,
    Run:   runConfigDefaults,
}

func init() {
    configCmd.AddCommand(configValidateCmd)
    configCmd.AddCommand(configResetCmd)
    configCmd.AddCommand(configDefaultsCmd)
}
```

**Validation Implementation:**
```go
type ValidationError struct {
    Key      string
    Value    interface{}
    Message  string
}

func validateConfig() ([]ValidationError, error) {
    errors := []ValidationError{}

    // Read config file
    home, _ := os.UserHomeDir()
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
```

**Reset Implementation:**
```go
func runConfigReset(cmd *cobra.Command, args []string) {
    home, _ := os.UserHomeDir()
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

        // Create backup
        timestamp := time.Now().Format("20060102-150405")
        backupPath := fmt.Sprintf("%s.backup.%s", configPath, timestamp)

        if _, err := os.Stat(configPath); err == nil {
            data, _ := os.ReadFile(configPath)
            os.WriteFile(backupPath, data, 0644)
        }

        // Write defaults
        defaults := getDefaultConfig()
        yamlData, _ := yaml.Marshal(defaults)

        tempPath := configPath + ".tmp"
        os.WriteFile(tempPath, yamlData, 0644)
        os.Rename(tempPath, configPath)

        fmt.Printf("✓ Configuration reset to defaults. Backup saved to %s\n", backupPath)
    } else {
        // Single key reset
        key := args[0]

        if !isValidConfigKey(key) {
            fmt.Fprintf(os.Stderr, "Error: Unknown configuration key '%s'\n", key)
            os.Exit(2)
        }

        // Get default value
        defaultValue := viper.GetString(key) // This gets the default from Viper

        // Load current config
        configData := make(map[string]interface{})
        if data, err := os.ReadFile(configPath); err == nil {
            yaml.Unmarshal(data, &configData)
        }

        // Set to default
        configData[key] = defaultValue

        // Write back
        yamlData, _ := yaml.Marshal(configData)
        tempPath := configPath + ".tmp"
        os.WriteFile(tempPath, yamlData, 0644)
        os.Rename(tempPath, configPath)

        fmt.Printf("✓ '%s' reset to default value: %v\n", key, defaultValue)
    }
}

func getDefaultConfig() map[string]interface{} {
    // Reset Viper and reload defaults
    v := viper.New()
    setDefaults() // Call same function used in config init

    defaults := make(map[string]interface{})
    for _, key := range validConfigKeys {
        defaults[key] = viper.Get(key)
    }
    return defaults
}
```

**Defaults Display:**
```go
func runConfigDefaults(cmd *cobra.Command, args []string) {
    defaults := getDefaultConfig()

    yamlData, err := yaml.Marshal(defaults)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Print(string(yamlData))
}
```

### Source Tree Components

**Files to Modify:**
- `cmd/config.go` - Add validate, reset, defaults subcommands
- `cmd/config_test.go` - Add unit tests for validation and reset
- `cmd/config_integration_test.go` - Add integration tests

**Files Referenced:**
- `internal/config/config.go` - setDefaults() function for default values

**Files Modified in Previous Stories:**
- Story 5.1: internal/config/config.go (Viper init, setDefaults())
- Story 5.2: cmd/config.go (get/set/list/unset)
- Story 5.3: cmd/config.go (output format keys)
- Story 5.4: cmd/config.go (editor_args)
- Story 5.5: cmd/config.go (profile commands)

### Testing Standards

**Unit Test Coverage:**
- Test validateConfig() with valid config returns no errors
- Test validateConfig() with invalid key returns error
- Test validateConfig() with invalid enum value returns error
- Test validateConfig() with wrong type returns error
- Test validation error messages are helpful
- Test full reset creates backup with timestamp
- Test single key reset preserves other keys
- Test defaults output includes all keys

**Integration Test Coverage:**
- Test validate on valid config file prints success
- Test validate on config with unknown key shows error
- Test validate on config with invalid enum shows error
- Test reset all with confirmation creates backup
- Test reset all without confirmation cancels
- Test reset single key updates only that key
- Test defaults output is valid YAML
- Test defaults can be piped to file

**Test Pattern:**
```go
func TestConfigValidate(t *testing.T) {
    t.Run("valid config passes", func(t *testing.T) {
        tmpHome := t.TempDir()
        t.Setenv("HOME", tmpHome)

        configPath := filepath.Join(tmpHome, ".dsa", "config.yaml")
        os.MkdirAll(filepath.Dir(configPath), 0755)

        validConfig := map[string]interface{}{
            "editor": "vim",
            "output_format": "text",
            "no_color": false,
        }
        yamlData, _ := yaml.Marshal(validConfig)
        os.WriteFile(configPath, yamlData, 0644)

        errors, err := validateConfig()

        assert.NoError(t, err)
        assert.Empty(t, errors)
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
}
```

### Technical Requirements

**Validation Checks:**
1. **Key Validity:** All keys must be in validConfigKeys
2. **Type Validation:** Values must match expected types
3. **Enum Validation:** Enum values must be in allowed list
4. **Range Validation:** Numeric values must meet constraints (if any)

**Valid Configuration Keys (from previous stories):**
- editor (string)
- output_format (enum: "text", "json")
- no_color (boolean)
- database_path (string)
- verbose (boolean)
- list_format (enum: "table", "json")
- status_format (enum: "table", "json")
- output_style (enum: "normal", "compact", "verbose")
- color_scheme (enum: "default", "solarized", "monokai", "nord")
- editor_args (string)

**Backup File Format:**
- Pattern: `config.yaml.backup.<timestamp>`
- Timestamp: `YYYYMMdd-HHmmss` (e.g., 20250129-143052)
- Location: Same directory as config.yaml

**Error Messages:**
- Should include key name, current value, and expected format
- Should provide actionable suggestion
- Should list valid options for enum types

### Previous Story Intelligence (Stories 5.1-5.5)

**Key Learnings:**
- parseConfigValue() already implements validation logic for each key type
- validConfigKeys array is the source of truth for valid keys
- setDefaults() in internal/config/config.go defines all default values
- Atomic write pattern used throughout (temp + rename)
- Confirmation prompts use bufio.Reader for stdin input
- gopkg.in/yaml.v3 for YAML marshaling
- testify for test assertions

**Established Patterns:**
- Cobra subcommand structure well-defined
- Config file operations use atomic writes
- Helpful error messages with suggestions
- Environment variables automatically bound
- Test coverage: 10+ tests per story

**Configuration System Maturity:**
- 10 config keys defined across 5 stories
- Viper integration complete
- Profile system implemented (Story 5.5)
- Get/set/list/unset commands working
- Precedence system functional

### Definition of Done

- [x] Config validate subcommand implemented
- [x] Validation checks all keys and values
- [x] Validation error messages are helpful
- [x] Config reset with confirmation implemented
- [x] Full reset creates timestamped backup
- [x] Single key reset preserves other keys
- [x] Config defaults displays all defaults in YAML
- [x] Unit tests: 8+ test scenarios (9 test scenarios added)
- [x] Integration tests: 8+ test scenarios (8 integration tests added)
- [x] All tests pass: `go test ./...`
- [x] Build succeeds: `go build`
- [x] Manual test: Validate valid config
- [x] Manual test: Validate invalid config shows errors
- [x] Manual test: Full reset creates backup
- [x] Manual test: Single key reset works
- [x] Manual test: Defaults output is correct
- [x] All acceptance criteria satisfied

## Dev Agent Record

### Context Reference

<!-- Path(s) to story context XML will be added here by context workflow -->

### Agent Model Used

claude-sonnet-4-5 (model ID: claude-sonnet-4-5-20250929)

### Debug Log References

### Completion Notes List

**Implementation Summary:**
- Successfully implemented complete configuration validation and reset system with 3 new subcommands
- All validation, reset, and defaults operations working correctly
- Comprehensive error collection and helpful error messages implemented
- Backup system with timestamped files for full reset
- Atomic write pattern maintained throughout
- All tests pass with zero compilation errors on first attempt
- Build succeeds with no errors

**Configuration Validation Features:**
- **Config Validate:** Checks all keys for validity, type correctness, enum values, and constraints
- Collects ALL errors before reporting (doesn't fail fast)
- Provides helpful error messages with current value, expected format, and fix suggestions
- Returns success message for valid config
- Exit code 2 for validation failures

**Configuration Reset Features:**
- **Full Reset:** Interactive confirmation prompt, creates timestamped backup, restores all defaults
- **Single Key Reset:** No confirmation required, resets only specified key, preserves other settings
- Backup format: `config.yaml.backup.YYYYMMdd-HHmmss`
- Default values mirror internal/config/config.go InitConfig() defaults

**Configuration Defaults Features:**
- **Config Defaults:** Displays all default values in YAML format
- Output can be piped or redirected
- Useful for creating config file templates
- Shows all 10 configuration keys with their default values

**Test Coverage:**
- Unit tests: 9 test scenarios covering validation, getDefaultConfig
  - Valid config passes
  - No config file is valid
  - Invalid key fails
  - Invalid enum value fails
  - Wrong type fails
  - Multiple errors collected
  - Invalid YAML returns error
  - Defaults returns all values
  - Defaults match expected values
- Integration tests: 8 comprehensive scenarios
  - Validate valid config
  - Validate corrupted config with detailed errors
  - Full reset creates backup
  - Single key reset preserves others
  - Reset defaults match expected
  - Defaults output is valid YAML
  - Defaults contains all keys
  - Defaults can be written to file
- All tests pass: PASS ok github.com/empire/dsa/cmd

**All Acceptance Criteria Satisfied:**
- AC1: Config validate checks all values ✓
- AC2: Validate shows detailed errors with suggestions ✓
- AC3: Full reset with confirmation and backup ✓
- AC4: Single key reset preserves other keys ✓
- AC5: Config defaults displays all defaults in YAML ✓

### File List

**Modified Files:**
- cmd/config.go: Added validation, reset, and defaults functionality
  - Line 10: Added "time" import for timestamp formatting
  - Lines 124-180: Added 3 new subcommands (validate, reset, defaults) with comprehensive help text
  - Lines 182-190: Updated init() to register new subcommands
  - Lines 360-497: Added 3 handler functions:
    - runConfigValidate: Validates config and displays errors
    - runConfigReset: Resets all or single key with confirmation and backup
    - runConfigDefaults: Displays all default values in YAML
  - Lines 591-677: Added helper functions:
    - ValidationError struct for error reporting
    - validateConfig: Comprehensive validation logic
    - getDefaultConfig: Returns all default values

- cmd/config_test.go: Added unit tests
  - Lines 527-688: Added 2 test functions with 9 test scenarios:
    - TestValidateConfig: 7 validation test scenarios
    - TestGetDefaultConfig: 2 default value test scenarios

- cmd/config_integration_test.go: Added integration tests
  - Lines 891-1114: Added 8 integration test functions:
    - TestIntegration_ValidateValidConfig: Valid config validation
    - TestIntegration_ValidateCorruptedConfig: Multi-error validation
    - TestIntegration_FullResetCreatesBackup: Backup creation testing
    - TestIntegration_SingleKeyResetPreservesOthers: Selective reset testing
    - TestIntegration_ResetToDefaultsMatchesExpected: Default value verification
    - TestIntegration_DefaultsOutputIsValidYAML: YAML format validation
    - TestIntegration_DefaultsContainsAllKeys: Completeness verification
    - TestIntegration_DefaultsCanBeWrittenToFile: File output testing

**No New Files Created:**
All functionality added to existing files.
