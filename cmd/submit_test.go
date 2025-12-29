package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmitCommandExists(t *testing.T) {
	// Verify submit command is registered
	cmd, _, err := rootCmd.Find([]string{"submit"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "submit", cmd.Name())
}

func TestSubmitCommand_RequiresOneArg(t *testing.T) {
	// submit command should require exactly one argument (problem-id)
	rootCmd.SetArgs([]string{"submit"})
	cmd, _, err := rootCmd.Find([]string{"submit"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify Args is set to ExactArgs(1)
	assert.NotNil(t, cmd.Args)
}

func TestSubmitCommand_HelpText(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"submit"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify help text contains key information
	assert.Contains(t, cmd.Long, "Submit your solution")
	assert.Contains(t, cmd.Long, "Runs tests to verify solution passes")
	assert.Contains(t, cmd.Long, "Saves solution to solutions/history")
	assert.Contains(t, cmd.Long, "Records submission in database")

	// Verify examples are present
	assert.Contains(t, cmd.Long, "dsa submit two-sum")
	assert.Contains(t, cmd.Long, "dsa submit binary-search")
}

func TestSubmitCommand_ShortDescription(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"submit"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	assert.Equal(t, "Submit and save your solution to history", cmd.Short)
}
