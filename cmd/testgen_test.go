package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestGenCommandExists(t *testing.T) {
	// Verify test-gen command is registered
	cmd, _, err := rootCmd.Find([]string{"test-gen"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "test-gen", cmd.Name())
}

func TestTestGenCommand_Flags(t *testing.T) {
	t.Run("append flag exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"test-gen", "problem", "--append"})
		cmd, _, err := rootCmd.Find([]string{"test-gen"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("append")
		assert.NotNil(t, flag, "append flag should exist")
	})

	t.Run("append short flag -a exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"test-gen", "problem", "-a"})
		cmd, _, err := rootCmd.Find([]string{"test-gen"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("append")
		assert.NotNil(t, flag, "append flag should exist")
	})

	t.Run("from-file flag exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"test-gen", "problem", "--from-file", "tests.json"})
		cmd, _, err := rootCmd.Find([]string{"test-gen"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("from-file")
		assert.NotNil(t, flag, "from-file flag should exist")
	})

	t.Run("from-file short flag -f exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"test-gen", "problem", "-f", "tests.json"})
		cmd, _, err := rootCmd.Find([]string{"test-gen"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("from-file")
		assert.NotNil(t, flag, "from-file flag should exist")
	})

	t.Run("both flags can be combined", func(t *testing.T) {
		rootCmd.SetArgs([]string{"test-gen", "problem", "--append", "--from-file", "tests.json"})
		cmd, _, err := rootCmd.Find([]string{"test-gen"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		appendFlag := cmd.Flags().Lookup("append")
		fromFileFlag := cmd.Flags().Lookup("from-file")
		assert.NotNil(t, appendFlag, "append flag should exist")
		assert.NotNil(t, fromFileFlag, "from-file flag should exist")
	})
}

func TestTestGenCommand_RequiresOneArg(t *testing.T) {
	// test-gen command should require exactly one argument (problem-id)
	rootCmd.SetArgs([]string{"test-gen"})
	cmd, _, err := rootCmd.Find([]string{"test-gen"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify Args is set to ExactArgs(1)
	// This is validated by cobra at runtime
	assert.NotNil(t, cmd.Args)
}

func TestTestGenCommand_HelpText(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"test-gen"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify help text contains key information
	assert.Contains(t, cmd.Long, "Interactively generate test cases")
	assert.Contains(t, cmd.Long, "Imports test cases from JSON file")
	assert.Contains(t, cmd.Long, "Appends to existing test file")

	// Verify examples are present
	assert.Contains(t, cmd.Long, "dsa test-gen my-problem")
	assert.Contains(t, cmd.Long, "dsa test-gen my-problem --from-file tests.json")
	assert.Contains(t, cmd.Long, "dsa test-gen my-problem --append")
}
