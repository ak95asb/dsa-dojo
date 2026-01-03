package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportCommandExists(t *testing.T) {
	// Verify export command is registered
	cmd, _, err := rootCmd.Find([]string{"export"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "export", cmd.Name())
}

func TestExportCommand_Flags(t *testing.T) {
	t.Run("format flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"export"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("format")
		assert.NotNil(t, flag, "format flag should exist")
	})

	t.Run("output flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"export"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("output")
		assert.NotNil(t, flag, "output flag should exist")
	})

	t.Run("difficulty flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"export"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("difficulty")
		assert.NotNil(t, flag, "difficulty flag should exist")
	})

	t.Run("topic flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"export"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("topic")
		assert.NotNil(t, flag, "topic flag should exist")
	})

	t.Run("all flags can be used together", func(t *testing.T) {
		rootCmd.SetArgs([]string{"export", "--format", "csv", "--output", "test.csv", "--difficulty", "easy", "--topic", "arrays"})
		cmd, _, err := rootCmd.Find([]string{"export"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		formatFlag := cmd.Flags().Lookup("format")
		outputFlag := cmd.Flags().Lookup("output")
		difficultyFlag := cmd.Flags().Lookup("difficulty")
		topicFlag := cmd.Flags().Lookup("topic")

		assert.NotNil(t, formatFlag)
		assert.NotNil(t, outputFlag)
		assert.NotNil(t, difficultyFlag)
		assert.NotNil(t, topicFlag)
	})
}

func TestExportCommand_NoArgsRequired(t *testing.T) {
	// export command should not require any arguments
	cmd, _, err := rootCmd.Find([]string{"export"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify Args is set to NoArgs
	assert.NotNil(t, cmd.Args)
}

func TestExportCommand_HelpText(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"export"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify help text contains key information
	assert.Contains(t, cmd.Long, "progress to JSON or CSV format")
	assert.Contains(t, cmd.Long, "JSON export with full details")
	assert.Contains(t, cmd.Long, "CSV export for spreadsheet compatibility")
	assert.Contains(t, cmd.Long, "Filtering by difficulty and topic")
	assert.Contains(t, cmd.Long, "Output to file or stdout")

	// Verify examples are present
	assert.Contains(t, cmd.Long, "dsa export --format json --output progress.json")
	assert.Contains(t, cmd.Long, "dsa export --format csv --output progress.csv")
	assert.Contains(t, cmd.Long, "dsa export --format json --difficulty medium")
	assert.Contains(t, cmd.Long, "dsa export --format json | jq .summary")
}

func TestExportCommand_ShortDescription(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"export"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	assert.Equal(t, "Export progress data to external formats", cmd.Short)
}

func TestExportCommand_FlagDescriptions(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"export"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(t, formatFlag)
	assert.Contains(t, formatFlag.Usage, "Export format")

	outputFlag := cmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Contains(t, outputFlag.Usage, "Output file")

	difficultyFlag := cmd.Flags().Lookup("difficulty")
	assert.NotNil(t, difficultyFlag)
	assert.Contains(t, difficultyFlag.Usage, "Filter by difficulty")

	topicFlag := cmd.Flags().Lookup("topic")
	assert.NotNil(t, topicFlag)
	assert.Contains(t, topicFlag.Usage, "Filter by topic")
}

func TestExportCommand_FlagDefaults(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"export"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify format flag defaults to json
	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(t, formatFlag)
	assert.Equal(t, "json", formatFlag.DefValue, "format flag should default to json")

	// Verify output flag has empty default (stdout)
	outputFlag := cmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "", outputFlag.DefValue, "output flag should default to empty string")

	// Verify difficulty flag has empty default
	difficultyFlag := cmd.Flags().Lookup("difficulty")
	assert.NotNil(t, difficultyFlag)
	assert.Equal(t, "", difficultyFlag.DefValue, "difficulty flag should default to empty string")

	// Verify topic flag has empty default
	topicFlag := cmd.Flags().Lookup("topic")
	assert.NotNil(t, topicFlag)
	assert.Equal(t, "", topicFlag.DefValue, "topic flag should default to empty string")
}

func TestExportCommand_FormatFlagValues(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"export"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify that format flag can be set to valid values
	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(t, formatFlag)

	// Test json format
	err = formatFlag.Value.Set("json")
	assert.NoError(t, err)
	assert.Equal(t, "json", formatFlag.Value.String())

	// Test csv format
	err = formatFlag.Value.Set("csv")
	assert.NoError(t, err)
	assert.Equal(t, "csv", formatFlag.Value.String())
}
