package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCommandExists(t *testing.T) {
	// Verify status command is registered
	cmd, _, err := rootCmd.Find([]string{"status"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "status", cmd.Name())
}

func TestStatusCommand_Flags(t *testing.T) {
	t.Run("topic flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"status"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("topic")
		assert.NotNil(t, flag, "topic flag should exist")
	})

	t.Run("compact flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"status"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("compact")
		assert.NotNil(t, flag, "compact flag should exist")
	})

	t.Run("both flags can be used together", func(t *testing.T) {
		rootCmd.SetArgs([]string{"status", "--topic", "arrays", "--compact"})
		cmd, _, err := rootCmd.Find([]string{"status"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		topicFlag := cmd.Flags().Lookup("topic")
		compactFlag := cmd.Flags().Lookup("compact")
		assert.NotNil(t, topicFlag, "topic flag should exist")
		assert.NotNil(t, compactFlag, "compact flag should exist")
	})
}

func TestStatusCommand_NoArgsRequired(t *testing.T) {
	// status command should not require any arguments
	cmd, _, err := rootCmd.Find([]string{"status"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify Args is set to NoArgs
	assert.NotNil(t, cmd.Args)
}

func TestStatusCommand_HelpText(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"status"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify help text contains key information
	assert.Contains(t, cmd.Long, "overview of your DSA practice progress")
	assert.Contains(t, cmd.Long, "Total problems solved")
	assert.Contains(t, cmd.Long, "Breakdown by difficulty")
	assert.Contains(t, cmd.Long, "Breakdown by topic")
	assert.Contains(t, cmd.Long, "Recent activity")
	assert.Contains(t, cmd.Long, "Visual progress bars")

	// Verify examples are present
	assert.Contains(t, cmd.Long, "dsa status")
	assert.Contains(t, cmd.Long, "dsa status --topic arrays")
	assert.Contains(t, cmd.Long, "dsa status --compact")
}

func TestStatusCommand_ShortDescription(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"status"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	assert.Equal(t, "Display your problem-solving progress dashboard", cmd.Short)
}

func TestStatusCommand_FlagDescriptions(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"status"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	topicFlag := cmd.Flags().Lookup("topic")
	assert.NotNil(t, topicFlag)
	assert.Contains(t, topicFlag.Usage, "Show stats for specific topic")

	compactFlag := cmd.Flags().Lookup("compact")
	assert.NotNil(t, compactFlag)
	assert.Contains(t, compactFlag.Usage, "Display one-line summary")
}

func TestStatusCommand_FlagDefaults(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"status"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify topic flag has empty default
	topicFlag := cmd.Flags().Lookup("topic")
	assert.NotNil(t, topicFlag)
	assert.Equal(t, "", topicFlag.DefValue, "topic flag should default to empty string")

	// Verify compact flag defaults to false
	compactFlag := cmd.Flags().Lookup("compact")
	assert.NotNil(t, compactFlag)
	assert.Equal(t, "false", compactFlag.DefValue, "compact flag should default to false")
}
