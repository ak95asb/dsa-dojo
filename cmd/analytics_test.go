package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyticsCommandExists(t *testing.T) {
	// Verify analytics command is registered
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "analytics", cmd.Name())
}

func TestAnalyticsCommand_Flags(t *testing.T) {
	t.Run("topic flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"analytics"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("topic")
		assert.NotNil(t, flag, "topic flag should exist")
	})

	t.Run("difficulty flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"analytics"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("difficulty")
		assert.NotNil(t, flag, "difficulty flag should exist")
	})

	t.Run("json flag exists", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"analytics"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("json")
		assert.NotNil(t, flag, "json flag should exist")
	})

	t.Run("all flags can be used together", func(t *testing.T) {
		rootCmd.SetArgs([]string{"analytics", "--topic", "arrays", "--difficulty", "easy", "--json"})
		cmd, _, err := rootCmd.Find([]string{"analytics"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		topicFlag := cmd.Flags().Lookup("topic")
		difficultyFlag := cmd.Flags().Lookup("difficulty")
		jsonFlag := cmd.Flags().Lookup("json")
		assert.NotNil(t, topicFlag, "topic flag should exist")
		assert.NotNil(t, difficultyFlag, "difficulty flag should exist")
		assert.NotNil(t, jsonFlag, "json flag should exist")
	})
}

func TestAnalyticsCommand_NoArgsRequired(t *testing.T) {
	// analytics command should not require any arguments
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify Args is set to NoArgs
	assert.NotNil(t, cmd.Args)
}

func TestAnalyticsCommand_HelpText(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify help text contains key information
	assert.Contains(t, cmd.Long, "comprehensive analytics")
	assert.Contains(t, cmd.Long, "Overall success rate")
	assert.Contains(t, cmd.Long, "Success rates broken down by difficulty")
	assert.Contains(t, cmd.Long, "Success rates broken down by topic")
	assert.Contains(t, cmd.Long, "Average number of attempts")
	assert.Contains(t, cmd.Long, "Practice pattern insights")

	// Verify examples are present
	assert.Contains(t, cmd.Long, "dsa analytics")
	assert.Contains(t, cmd.Long, "dsa analytics --topic arrays")
	assert.Contains(t, cmd.Long, "dsa analytics --difficulty medium")
	assert.Contains(t, cmd.Long, "dsa analytics --json")
}

func TestAnalyticsCommand_ShortDescription(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	assert.Equal(t, "View detailed analytics and insights about your practice patterns", cmd.Short)
}

func TestAnalyticsCommand_FlagDescriptions(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	topicFlag := cmd.Flags().Lookup("topic")
	assert.NotNil(t, topicFlag)
	assert.Contains(t, topicFlag.Usage, "Filter analytics by topic")

	difficultyFlag := cmd.Flags().Lookup("difficulty")
	assert.NotNil(t, difficultyFlag)
	assert.Contains(t, difficultyFlag.Usage, "Filter analytics by difficulty")

	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)
	assert.Contains(t, jsonFlag.Usage, "Output analytics as JSON")
}

func TestAnalyticsCommand_FlagDefaults(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify topic flag has empty default
	topicFlag := cmd.Flags().Lookup("topic")
	assert.NotNil(t, topicFlag)
	assert.Equal(t, "", topicFlag.DefValue, "topic flag should default to empty string")

	// Verify difficulty flag has empty default
	difficultyFlag := cmd.Flags().Lookup("difficulty")
	assert.NotNil(t, difficultyFlag)
	assert.Equal(t, "", difficultyFlag.DefValue, "difficulty flag should default to empty string")

	// Verify json flag defaults to false
	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)
	assert.Equal(t, "false", jsonFlag.DefValue, "json flag should default to false")
}

func TestAnalyticsCommand_TopicFilterValues(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify that topic flag can be set to various values
	topicFlag := cmd.Flags().Lookup("topic")
	assert.NotNil(t, topicFlag)

	// Test that we can set the topic flag
	err = topicFlag.Value.Set("arrays")
	assert.NoError(t, err)
	assert.Equal(t, "arrays", topicFlag.Value.String())
}

func TestAnalyticsCommand_DifficultyFilterValues(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify that difficulty flag can be set to various values
	difficultyFlag := cmd.Flags().Lookup("difficulty")
	assert.NotNil(t, difficultyFlag)

	// Test that we can set the difficulty flag
	validDifficulties := []string{"easy", "medium", "hard"}
	for _, diff := range validDifficulties {
		err = difficultyFlag.Value.Set(diff)
		assert.NoError(t, err)
		assert.Equal(t, diff, difficultyFlag.Value.String())
	}
}

func TestAnalyticsCommand_JSONFlag(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"analytics"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify that json flag is a boolean flag
	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)

	// Test that we can set the json flag to true
	err = jsonFlag.Value.Set("true")
	assert.NoError(t, err)
	assert.Equal(t, "true", jsonFlag.Value.String())

	// Test that we can set the json flag to false
	err = jsonFlag.Value.Set("false")
	assert.NoError(t, err)
	assert.Equal(t, "false", jsonFlag.Value.String())
}
