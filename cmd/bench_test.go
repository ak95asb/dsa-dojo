package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBenchCommandExists(t *testing.T) {
	// Verify bench command is registered
	cmd, _, err := rootCmd.Find([]string{"bench"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "bench", cmd.Name())
}

func TestBenchCommand_Flags(t *testing.T) {
	t.Run("save flag exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"bench", "problem", "--save"})
		cmd, _, err := rootCmd.Find([]string{"bench"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("save")
		assert.NotNil(t, flag, "save flag should exist")
	})

	t.Run("mem flag exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"bench", "problem", "--mem"})
		cmd, _, err := rootCmd.Find([]string{"bench"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("mem")
		assert.NotNil(t, flag, "mem flag should exist")
	})

	t.Run("cpuprofile flag exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"bench", "problem", "--cpuprofile=test.prof"})
		cmd, _, err := rootCmd.Find([]string{"bench"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("cpuprofile")
		assert.NotNil(t, flag, "cpuprofile flag should exist")
	})

	t.Run("memprofile flag exists", func(t *testing.T) {
		rootCmd.SetArgs([]string{"bench", "problem", "--memprofile=mem.prof"})
		cmd, _, err := rootCmd.Find([]string{"bench"})
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		flag := cmd.Flags().Lookup("memprofile")
		assert.NotNil(t, flag, "memprofile flag should exist")
	})
}

func TestBenchCommand_RequiresOneArg(t *testing.T) {
	// bench command should require exactly one argument (problem-id)
	rootCmd.SetArgs([]string{"bench"})
	cmd, _, err := rootCmd.Find([]string{"bench"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify Args is set to ExactArgs(1)
	assert.NotNil(t, cmd.Args)
}

func TestBenchCommand_HelpText(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"bench"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	// Verify help text contains key information
	assert.Contains(t, cmd.Long, "Execute Go benchmarks")
	assert.Contains(t, cmd.Long, "go test -bench")
	assert.Contains(t, cmd.Long, "iterations")
	assert.Contains(t, cmd.Long, "profiling")

	// Verify examples are present
	assert.Contains(t, cmd.Long, "dsa bench two-sum")
	assert.Contains(t, cmd.Long, "dsa bench two-sum --save")
	assert.Contains(t, cmd.Long, "dsa bench two-sum --mem")
	assert.Contains(t, cmd.Long, "dsa bench two-sum --cpuprofile")
}

func TestBenchCommand_ShortDescription(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"bench"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	assert.Equal(t, "Run performance benchmarks on your solution", cmd.Short)
}

func TestBenchCommand_FlagDescriptions(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"bench"})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	saveFlag := cmd.Flags().Lookup("save")
	assert.NotNil(t, saveFlag)
	assert.Contains(t, saveFlag.Usage, "Save benchmark results")

	memFlag := cmd.Flags().Lookup("mem")
	assert.NotNil(t, memFlag)
	assert.Contains(t, memFlag.Usage, "memory profiling")

	cpuFlag := cmd.Flags().Lookup("cpuprofile")
	assert.NotNil(t, cpuFlag)
	assert.Contains(t, cpuFlag.Usage, "CPU profile")

	memProfileFlag := cmd.Flags().Lookup("memprofile")
	assert.NotNil(t, memProfileFlag)
	assert.Contains(t, memProfileFlag.Usage, "memory profile")
}
