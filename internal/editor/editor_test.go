package editor

import (
	"os"
	"runtime"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDetect(t *testing.T) {
	// Save original state
	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	// Reset viper for clean test state
	viper.Reset()

	t.Run("detects editor from viper config", func(t *testing.T) {
		viper.Set("editor", "code")
		defer viper.Reset()

		editor := Detect()
		assert.Equal(t, "code", editor)
	})

	t.Run("falls back to EDITOR env var when config not set", func(t *testing.T) {
		viper.Reset()
		os.Setenv("EDITOR", "nano")
		defer os.Unsetenv("EDITOR")

		editor := Detect()
		assert.Equal(t, "nano", editor)
	})

	t.Run("uses platform default when no config or env var", func(t *testing.T) {
		viper.Reset()
		os.Unsetenv("EDITOR")

		editor := Detect()

		switch runtime.GOOS {
		case "windows":
			assert.Equal(t, "notepad", editor)
		default:
			assert.Equal(t, "vi", editor)
		}
	})

	t.Run("prioritizes config over env var", func(t *testing.T) {
		viper.Set("editor", "emacs")
		os.Setenv("EDITOR", "vim")
		defer viper.Reset()
		defer os.Unsetenv("EDITOR")

		editor := Detect()
		assert.Equal(t, "emacs", editor)
	})
}

func TestLaunch(t *testing.T) {
	t.Run("returns error for non-existent editor", func(t *testing.T) {
		tmpFile := t.TempDir() + "/test.txt"
		os.WriteFile(tmpFile, []byte("test"), 0644)

		err := Launch("non-existent-editor-12345", tmpFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to launch")
	})

	// Note: Testing successful editor launch is difficult without mocking exec.Command
	// These tests verify the error cases and basic functionality
	// Integration tests will verify the full workflow
}
