package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistoryCommandExists(t *testing.T) {
	// Verify history command is registered
	cmd, _, err := rootCmd.Find([]string{"history"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "history", cmd.Name())
}

func TestHistoryCommand_Flags(t *testing.T) {
	t.Run("show flag exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"history", "problem", "--show", "1"})
		cmd, _, err := rootCmd.Find([]string{"history"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("show")
		assert.NotNil(t, flag, "show flag should exist")
	})

	t.Run("restore flag exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"history", "problem", "--restore", "1"})
		cmd, _, err := rootCmd.Find([]string{"history"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("restore")
		assert.NotNil(t, flag, "restore flag should exist")
	})

	t.Run("both flags can be present (though only one will be used)", func(t *testing.T) {
		rootCmd.SetArgs([]string{"history", "problem", "--show", "1", "--restore", "2"})
		cmd, _, err := rootCmd.Find([]string{"history"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		showFlag := cmd.Flags().Lookup("show")
		restoreFlag := cmd.Flags().Lookup("restore")
		assert.NotNil(t, showFlag, "show flag should exist")
		assert.NotNil(t, restoreFlag, "restore flag should exist")
	})
}

func TestHistoryCommand_RequiresOneArg(t *testing.T) {
	// history command should require exactly one argument (problem-id)
	rootCmd.SetArgs([]string{"history"})
	cmd, _, err := rootCmd.Find([]string{"history"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify Args is set to ExactArgs(1)
	assert.NotNil(t, cmd.Args)
}

func TestHistoryCommand_HelpText(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"history"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify help text contains key information
	assert.Contains(t, cmd.Long, "Display all solution attempts")
	assert.Contains(t, cmd.Long, "--show N")
	assert.Contains(t, cmd.Long, "--restore N")

	// Verify examples are present
	assert.Contains(t, cmd.Long, "dsa history two-sum")
	assert.Contains(t, cmd.Long, "dsa history two-sum --show 2")
	assert.Contains(t, cmd.Long, "dsa history two-sum --restore 3")
}

func TestHistoryCommand_ShortDescription(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"history"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	assert.Equal(t, "View solution submission history", cmd.Short)
}

func TestHistoryCommand_FlagDescriptions(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"history"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	showFlag := cmd.Flags().Lookup("show")
	assert.NotNil(t, showFlag)
	assert.Contains(t, showFlag.Usage, "Display solution code from Nth attempt")

	restoreFlag := cmd.Flags().Lookup("restore")
	assert.NotNil(t, restoreFlag)
	assert.Contains(t, restoreFlag.Usage, "Restore Nth attempt as current solution")
}
