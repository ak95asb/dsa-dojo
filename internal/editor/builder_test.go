package editor

import (
	"os"
	"runtime"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBuildCommand(t *testing.T) {
	t.Run("simple editor with no args", func(t *testing.T) {
		viper.Reset()
		viper.Set("editor", "vim")
		viper.Set("editor_args", "")

		cmd := BuildCommand("/path/file.go", 1, 1)

		assert.Equal(t, []string{"vim", "/path/file.go"}, cmd)
	})

	t.Run("VS Code with goto placeholder", func(t *testing.T) {
		viper.Reset()
		viper.Set("editor", "code")
		viper.Set("editor_args", "--goto {file}:{line}")

		cmd := BuildCommand("/path/file.go", 10, 5)

		assert.Equal(t, []string{"code", "--goto", "/path/file.go:10"}, cmd)
	})

	t.Run("Vim with line placeholder", func(t *testing.T) {
		viper.Reset()
		viper.Set("editor", "vim")
		viper.Set("editor_args", "+{line}")

		cmd := BuildCommand("/path/file.go", 42, 1)

		assert.Equal(t, []string{"vim", "+42", "/path/file.go"}, cmd)
	})

	t.Run("Neovim with line placeholder", func(t *testing.T) {
		viper.Reset()
		viper.Set("editor", "nvim")
		viper.Set("editor_args", "-c normal {line}G")

		cmd := BuildCommand("/path/file.go", 25, 1)

		assert.Equal(t, []string{"nvim", "-c", "normal", "25G", "/path/file.go"}, cmd)
	})

	t.Run("Emacs with line and column placeholders", func(t *testing.T) {
		viper.Reset()
		viper.Set("editor", "emacs")
		viper.Set("editor_args", "+{line}:{column}")

		cmd := BuildCommand("/path/file.go", 15, 8)

		assert.Equal(t, []string{"emacs", "+15:8", "/path/file.go"}, cmd)
	})

	t.Run("all placeholders substituted", func(t *testing.T) {
		viper.Reset()
		viper.Set("editor", "myeditor")
		viper.Set("editor_args", "{file} --line {line} --col {column}")

		cmd := BuildCommand("/src/main.go", 100, 50)

		assert.Equal(t, []string{"myeditor", "/src/main.go", "--line", "100", "--col", "50"}, cmd)
	})

	t.Run("fallback to OS default when editor not configured", func(t *testing.T) {
		viper.Reset()
		viper.Set("editor", "")
		viper.Set("editor_args", "")
		os.Unsetenv("EDITOR")

		cmd := BuildCommand("/path/file.go", 1, 1)

		// Should use OS-specific default
		expectedEditor := GetSystemDefaultEditor()
		assert.Equal(t, []string{expectedEditor, "/path/file.go"}, cmd)
	})
}

func TestGetFallbackEditor(t *testing.T) {
	t.Run("returns EDITOR env var when set", func(t *testing.T) {
		t.Setenv("EDITOR", "nano")

		editor := GetFallbackEditor()

		assert.Equal(t, "nano", editor)
	})

	t.Run("returns OS default when EDITOR not set", func(t *testing.T) {
		os.Unsetenv("EDITOR")

		editor := GetFallbackEditor()

		expectedEditor := GetSystemDefaultEditor()
		assert.Equal(t, expectedEditor, editor)
	})
}

func TestGetSystemDefaultEditor(t *testing.T) {
	t.Run("returns correct default for current OS", func(t *testing.T) {
		editor := GetSystemDefaultEditor()

		switch runtime.GOOS {
		case "darwin":
			assert.Equal(t, "open", editor)
		case "linux":
			assert.Equal(t, "xdg-open", editor)
		case "windows":
			assert.Equal(t, "start", editor)
		default:
			assert.Equal(t, "vi", editor)
		}
	})
}
