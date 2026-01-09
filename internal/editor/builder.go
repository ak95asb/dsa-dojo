package editor

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// BuildCommand constructs the editor command with argument substitution
func BuildCommand(filePath string, line, column int) []string {
	editor := viper.GetString("editor")
	editorArgs := viper.GetString("editor_args")

	// Apply fallback logic if editor not configured
	if editor == "" {
		editor = GetFallbackEditor()
	}

	// If no editor_args configured, use simple "editor file" format
	if editorArgs == "" {
		return []string{editor, filePath}
	}

	// Check if {file} placeholder exists in editor_args
	hasFilePlaceholder := strings.Contains(editorArgs, "{file}")

	// Replace placeholders in editor_args
	args := strings.ReplaceAll(editorArgs, "{file}", filePath)
	args = strings.ReplaceAll(args, "{line}", fmt.Sprintf("%d", line))
	args = strings.ReplaceAll(args, "{column}", fmt.Sprintf("%d", column))

	// Parse arguments (simple space split for now)
	argList := strings.Fields(args)

	// Build final command: editor + parsed args
	cmd := append([]string{editor}, argList...)

	// If {file} placeholder was not in args, append file path at the end
	if !hasFilePlaceholder {
		cmd = append(cmd, filePath)
	}

	return cmd
}

// GetFallbackEditor returns editor from $EDITOR or OS-specific default
func GetFallbackEditor() string {
	// Try $EDITOR environment variable first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	// Fall back to OS-specific defaults
	return GetSystemDefaultEditor()
}

// GetSystemDefaultEditor returns the OS-specific default editor
func GetSystemDefaultEditor() string {
	switch runtime.GOOS {
	case "darwin":
		return "open"
	case "linux":
		return "xdg-open"
	case "windows":
		return "start"
	default:
		return "vi"
	}
}
