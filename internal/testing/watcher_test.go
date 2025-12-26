package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTestState_DetectTransition_FirstRun tests initial state detection
func TestTestState_DetectTransition_FirstRun(t *testing.T) {
	state := &TestState{}
	result := &TestResult{AllPassed: true}

	transition := state.DetectTransition(result)

	assert.Equal(t, "", transition, "First run should return empty string")
	assert.NotNil(t, state.LastResult, "Last result should be set after first run")
	assert.True(t, state.LastResult.AllPassed, "Last result should match current")
}

// TestTestState_DetectTransition_FailToPass tests fail→pass transition
func TestTestState_DetectTransition_FailToPass(t *testing.T) {
	state := &TestState{
		LastResult: &TestResult{AllPassed: false},
	}
	currentResult := &TestResult{AllPassed: true}

	transition := state.DetectTransition(currentResult)

	assert.Equal(t, "fail_to_pass", transition, "Should detect fail→pass transition")
	assert.True(t, state.LastResult.AllPassed, "Last result should be updated to passing")
}

// TestTestState_DetectTransition_PassToFail tests pass→fail transition
func TestTestState_DetectTransition_PassToFail(t *testing.T) {
	state := &TestState{
		LastResult: &TestResult{AllPassed: true},
	}
	currentResult := &TestResult{AllPassed: false}

	transition := state.DetectTransition(currentResult)

	assert.Equal(t, "pass_to_fail", transition, "Should detect pass→fail transition")
	assert.False(t, state.LastResult.AllPassed, "Last result should be updated to failing")
}

// TestTestState_DetectTransition_NoChange_AllPassing tests no change when passing
func TestTestState_DetectTransition_NoChange_AllPassing(t *testing.T) {
	state := &TestState{
		LastResult: &TestResult{AllPassed: true},
	}
	currentResult := &TestResult{AllPassed: true}

	transition := state.DetectTransition(currentResult)

	assert.Equal(t, "no_change", transition, "Should detect no change")
	assert.True(t, state.LastResult.AllPassed, "Last result should remain passing")
}

// TestTestState_DetectTransition_NoChange_AllFailing tests no change when failing
func TestTestState_DetectTransition_NoChange_AllFailing(t *testing.T) {
	state := &TestState{
		LastResult: &TestResult{AllPassed: false},
	}
	currentResult := &TestResult{AllPassed: false}

	transition := state.DetectTransition(currentResult)

	assert.Equal(t, "no_change", transition, "Should detect no change")
	assert.False(t, state.LastResult.AllPassed, "Last result should remain failing")
}

// TestTestState_DetectTransition_MultipleTransitions tests sequence of transitions
func TestTestState_DetectTransition_MultipleTransitions(t *testing.T) {
	state := &TestState{}

	// First run - failing
	result1 := &TestResult{AllPassed: false}
	transition1 := state.DetectTransition(result1)
	assert.Equal(t, "", transition1, "First run")

	// Fail → Pass
	result2 := &TestResult{AllPassed: true}
	transition2 := state.DetectTransition(result2)
	assert.Equal(t, "fail_to_pass", transition2)

	// Pass → Pass (no change)
	result3 := &TestResult{AllPassed: true}
	transition3 := state.DetectTransition(result3)
	assert.Equal(t, "no_change", transition3)

	// Pass → Fail
	result4 := &TestResult{AllPassed: false}
	transition4 := state.DetectTransition(result4)
	assert.Equal(t, "pass_to_fail", transition4)

	// Fail → Fail (no change)
	result5 := &TestResult{AllPassed: false}
	transition5 := state.DetectTransition(result5)
	assert.Equal(t, "no_change", transition5)

	// Fail → Pass (again)
	result6 := &TestResult{AllPassed: true}
	transition6 := state.DetectTransition(result6)
	assert.Equal(t, "fail_to_pass", transition6)
}

// TestTestState_DetectTransition_StateUpdate tests that state is updated correctly
func TestTestState_DetectTransition_StateUpdate(t *testing.T) {
	state := &TestState{}

	// First result
	result1 := &TestResult{
		AllPassed:   true,
		PassedCount: 5,
		TotalCount:  5,
	}
	state.DetectTransition(result1)

	assert.Equal(t, 5, state.LastResult.PassedCount, "PassedCount should be stored")
	assert.Equal(t, 5, state.LastResult.TotalCount, "TotalCount should be stored")

	// Second result
	result2 := &TestResult{
		AllPassed:   false,
		PassedCount: 3,
		TotalCount:  5,
	}
	state.DetectTransition(result2)

	assert.Equal(t, 3, state.LastResult.PassedCount, "PassedCount should be updated")
	assert.Equal(t, 5, state.LastResult.TotalCount, "TotalCount should remain same")
	assert.False(t, state.LastResult.AllPassed, "AllPassed should be updated")
}

// TestClearTerminal_DoesNotPanic tests that clearTerminal doesn't crash
func TestClearTerminal_DoesNotPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		clearTerminal()
	}, "clearTerminal should not panic")
}
