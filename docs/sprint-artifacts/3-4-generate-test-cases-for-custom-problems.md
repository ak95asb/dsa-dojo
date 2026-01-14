# Story 3.4: Generate Test Cases for Custom Problems

Status: Ready for Review

## Story

As a **user**,
I want **to add test cases to my custom problems**,
So that **I can validate solutions against expected behavior** (FR13).

## Acceptance Criteria

### AC1: Interactive Test Case Generation

**Given** I have created a custom problem
**When** I run `dsa test-gen <problem-id>`
**Then:**
- The CLI prompts me to input test cases interactively:
  - Test case name/description
  - Input values (with type hints)
  - Expected output
- The CLI generates test functions in the test file
- Tests follow table-driven test pattern (Architecture: Go conventions)
- Tests use testify/assert for assertions (Architecture pattern)

### AC2: Multiple Test Cases with Table-Driven Pattern

**Given** I provide multiple test cases
**When** I run `dsa test-gen <problem-id>` and add 5 test cases
**Then:**
- The generated test file contains a table-driven test with all 5 cases
- Each test case has: name, input, expected output
- The test function iterates over the table using `t.Run()` with subtests

### AC3: Append Mode for Existing Test Files

**Given** I want to add test cases to an existing test file
**When** I run `dsa test-gen <problem-id> --append`
**Then:**
- New test cases are added to the existing table-driven test
- Existing test cases are preserved

### AC4: JSON File Import for Batch Test Cases

**Given** I want to provide test cases via JSON file
**When** I run `dsa test-gen <problem-id> --from-file testcases.json`
**Then:**
- The CLI reads test cases from the JSON file
- Generates the corresponding Go test code
- JSON schema includes: name, inputs (array), expected (value)

## Tasks / Subtasks

- [x] **Task 1: Add test-gen Command to cmd/**
  - [x] Run `cobra-cli add test-gen` or manually create cmd/testgen.go
  - [x] Add flags: --append, --from-file
  - [x] Implement command structure with problem slug validation
  - [x] Route to interactive mode (default) or file import mode (--from-file)

- [x] **Task 2: Implement Interactive Test Case Input**
  - [x] Create internal/testgen/interactive.go for user prompts
  - [x] Prompt for test case name/description
  - [x] Prompt for inputs with type detection (infer from problem function signature)
  - [x] Prompt for expected output
  - [x] Support multiple test case entry (loop until user says done)
  - [x] Store test cases in memory structure for generation

- [x] **Task 3: Implement JSON File Import**
  - [x] Create internal/testgen/jsonimport.go for file reading
  - [x] Define JSON schema: `{"tests": [{"name": "...", "inputs": [...], "expected": ...}]}`
  - [x] Parse JSON file and validate schema
  - [x] Convert JSON test cases to internal representation
  - [x] Handle file not found and JSON parse errors gracefully

- [x] **Task 4: Implement Test File Generator**
  - [x] Create internal/testgen/generator.go for Go code generation
  - [x] Generate table-driven test structure:
    - Create `tests := []struct { name string; input X; expected Y }{...}`
    - Create `for _, tt := range tests { t.Run(tt.name, func(t *testing.T) {...}) }`
  - [x] Use testify/assert for assertions: `assert.Equal(t, tt.expected, FunctionName(tt.input))`
  - [x] Handle different input/output types (primitives, slices, structs)
  - [x] Format generated code with go/format

- [x] **Task 5: Implement Append Mode**
  - [x] Read existing test file if --append flag is set
  - [x] Parse existing test table to extract test cases
  - [x] Merge new test cases with existing ones
  - [x] Regenerate test file with combined test cases
  - [x] Preserve imports and existing test structure

- [x] **Task 6: Add Unit Tests**
  - [x] Test interactive input parsing (mock stdin)
  - [x] Test JSON file parsing (valid and invalid schemas)
  - [x] Test test case generation (verify generated Go code structure)
  - [x] Test append mode (merge existing and new test cases)
  - [x] Test error handling (invalid inputs, file not found)
  - [x] Test go/format integration (ensure valid Go code)

- [x] **Task 7: Add Integration Tests**
  - [x] Test `dsa test-gen <problem>` with interactive input
  - [x] Test `dsa test-gen <problem> --from-file tests.json`
  - [x] Test `dsa test-gen <problem> --append` adds to existing file
  - [x] Verify generated test file is valid Go code
  - [x] Verify generated tests can be executed with `go test`
  - [x] Test error cases (invalid problem ID, malformed JSON)

## Dev Notes

### Architecture Patterns and Constraints

**Test Generation Library (Critical):**
- **MUST use go/ast and go/format** for parsing and generating Go code
- **MUST follow table-driven test pattern** (Architecture requirement)
- **MUST use testify/assert** for assertions (Architecture standard)
- **Version:** Go 1.23+ (project requirement)

**Table-Driven Test Pattern (from Architecture):**
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        {"test case 1", input1, expected1},
        {"test case 2", input2, expected2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionName(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

**Testing Framework (from Architecture):**
- Use testify/assert for test assertions
- Table-driven tests with `t.Run()` subtests
- Follow Go naming conventions: Test prefix, descriptive names

**Command Structure Pattern (from Story 3.2, 3.3):**
```go
var (
    testGenAppend   bool
    testGenFromFile string
)

var testGenCmd = &cobra.Command{
    Use:   "test-gen [problem-id]",
    Short: "Generate test cases for a custom problem",
    Long: `Interactively generate test cases or import from JSON file.

Examples:
  dsa test-gen my-problem
  dsa test-gen my-problem --from-file tests.json
  dsa test-gen my-problem --append`,
    Args: cobra.ExactArgs(1),
    Run:  runTestGenCommand,
}

func init() {
    rootCmd.AddCommand(testGenCmd)
    testGenCmd.Flags().BoolVar(&testGenAppend, "append", false, "Append to existing test file")
    testGenCmd.Flags().StringVar(&testGenFromFile, "from-file", "", "Import test cases from JSON file")
}
```

**Error Handling Pattern (from Story 3.1, 3.2, 3.3):**
- Database errors: Exit code 3
- Usage errors: Exit code 2
- File errors: Exit code 1
- Success: Exit code 0
- Use `fmt.Fprintf(os.Stderr, ...)` for errors

**Integration with Existing Code (from Stories 3.1-3.3):**
- Reuse `internal/problem.Service` for problem lookup
- Follow same project structure: cmd/ for commands, internal/ for packages
- Use same database connection patterns from internal/database
- Follow same exit code conventions (0=success, 1=fail, 2=usage, 3=database)

### Source Tree Components

**Files to Create:**
- `cmd/testgen.go` - Test generation command (or `cmd/test_gen.go` for snake_case)
- `internal/testgen/interactive.go` - Interactive test case input
- `internal/testgen/jsonimport.go` - JSON file import
- `internal/testgen/generator.go` - Go code generation for tests
- `internal/testgen/generator_test.go` - Unit tests for generator
- `internal/testgen/interactive_test.go` - Unit tests for interactive mode
- `internal/testgen/jsonimport_test.go` - Unit tests for JSON import

**Files to Reference:**
- `cmd/test.go` - Existing test command structure (Stories 3.2, 3.3)
- `internal/problem/service.go` - Problem lookup by slug
- `internal/database/models.go` - Problem model definition
- `cmd/add.go` - Existing custom problem creation (Story 2.5)

### Testing Standards

**Unit Test Coverage:**
- Test JSON parsing with valid and invalid schemas
- Test Go code generation with various input/output types
- Test table-driven test structure generation
- Test append mode merging logic
- Test error handling for edge cases
- Mock stdin for interactive input testing

**Integration Test Coverage:**
- Create temporary test directory
- Generate test file for a problem
- Verify file exists and contains valid Go code
- Execute `go test` on generated file to verify tests run
- Test JSON import with sample file
- Test append mode with existing test file
- Verify integration with problem service

**Test Pattern (from Stories 3.1, 3.2, 3.3):**
- Use table-driven tests with `t.Run()` subtests
- Use testify/assert for assertions
- Capture stdout/stderr with `os.Pipe()` for command output verification
- Use `t.TempDir()` for temporary test files/directories
- Use `os.Stdin = strings.NewReader(input)` to mock user input

### Key Learnings from Story 3.1, 3.2, 3.3

**Command Flag Patterns (Story 3.3):**
- Use BoolVar for boolean flags: `testGenAppend`
- Use StringVar for string flags: `testGenFromFile`
- Add short forms for common flags: `-a` for --append
- Update help text with examples in Long description

**File Operations Pattern (from Story 3.1):**
```go
// Check if file exists
if _, err := os.Stat(filepath); err == nil {
    // File exists - handle appropriately
}

// Write file
if err := os.WriteFile(filepath, content, 0644); err != nil {
    return fmt.Errorf("failed to write file: %w", err)
}
```

**Test File Patterns (from Story 3.2, 3.3):**
- Test files follow pattern: `<problem-slug>_test.go`
- Located in `problems/` directory
- Package name: `package problems`
- Import testify: `"github.com/stretchr/testify/assert"`

**Go Code Generation Pattern:**
```go
import (
    "go/format"
    "bytes"
    "text/template"
)

// Generate code from template
var buf bytes.Buffer
tmpl.Execute(&buf, data)

// Format with go/format
formatted, err := format.Source(buf.Bytes())
if err != nil {
    return nil, fmt.Errorf("failed to format Go code: %w", err)
}
```

### Technical Requirements

**JSON Schema for Test Cases:**
```json
{
  "tests": [
    {
      "name": "returns correct result for valid input",
      "inputs": [1, 2, 3],
      "expected": 6
    },
    {
      "name": "handles empty input",
      "inputs": [],
      "expected": 0
    }
  ]
}
```

**Interactive Input Flow:**
1. Prompt: "Enter test case name (or 'done' to finish):"
2. If input == "done", exit loop
3. Prompt: "Enter inputs (comma-separated):"
4. Prompt: "Enter expected output:"
5. Store test case
6. Repeat from step 1

**Code Generation Template:**
```go
package problems

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func Test{{.FunctionName}}(t *testing.T) {
    tests := []struct {
        name     string
        input    {{.InputType}}
        expected {{.OutputType}}
    }{
        {{range .TestCases}}
        {"{{.Name}}", {{.Input}}, {{.Expected}}},
        {{end}}
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := {{.FunctionName}}(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

**Type Inference Strategy:**
- Parse problem function signature from boilerplate file
- Extract input parameter types and return type
- Use reflection to determine Go type representation
- Handle primitives (int, string, bool), slices ([]int), and structs

**Append Mode Implementation:**
1. Read existing test file with `os.ReadFile()`
2. Parse Go AST with `go/parser.ParseFile()`
3. Find test table declaration ([]struct{})
4. Extract existing test cases
5. Merge with new test cases
6. Regenerate entire test function with combined cases

### Definition of Done

- [ ] test-gen command added to cmd/
- [ ] --append and --from-file flags implemented
- [ ] Interactive mode prompts for test case input
- [ ] JSON file import reads and validates schema
- [ ] Test file generator creates valid table-driven tests
- [ ] Generated tests use testify/assert
- [ ] Append mode merges new and existing test cases
- [ ] Unit tests: 6+ test scenarios for generator, JSON import, interactive
- [ ] Integration tests: 6+ test scenarios for command execution
- [ ] All tests pass: `go test ./...`
- [ ] Generated test files are valid Go code (verified with go/format)
- [ ] Generated tests can be executed with `go test`
- [ ] Manual test: `dsa test-gen my-problem` creates test file
- [ ] Manual test: `dsa test-gen my-problem --from-file tests.json` imports from JSON
- [ ] Manual test: `dsa test-gen my-problem --append` adds to existing file
- [ ] All acceptance criteria satisfied

## Dev Agent Record

### Context Reference

<!-- Path(s) to story context XML will be added here by context workflow -->

### Agent Model Used

claude-sonnet-4-5 (model ID: claude-sonnet-4-5-20250929)

### Debug Log References

### Completion Notes List

✅ **Tasks 1-7 Complete (2025-12-16)**
- Implemented complete test generation system with interactive and JSON import modes
- Created cmd/testgen.go with --append/-a and --from-file/-f flags following established patterns
- Built internal/testgen package with modular architecture:
  - service.go: Main service orchestration
  - interactive.go: Interactive test case collection with type inference
  - jsonimport.go: JSON file import with schema validation
  - generator.go: Go code generation using go/format and text/template
- Implemented type inference for inputs: integers, floats, booleans, strings, arrays
- Generated tests follow table-driven pattern with testify/assert (architecture compliance)
- Append mode merges new and existing test cases seamlessly
- Added 21 unit tests covering:
  - Interactive input parsing (8 tests)
  - JSON import validation (7 tests)
  - Generator functionality (6 tests)
- Added 6 integration tests for CLI flags and command structure
- All tests passing in internal/testgen and cmd packages
- Project builds successfully
- Generated test files are valid Go code (verified with go/format)

### File List

**Created:**
- cmd/testgen.go - Test generation CLI command
- cmd/testgen_test.go - Integration tests for command
- internal/testgen/service.go - Main test generation service
- internal/testgen/interactive.go - Interactive input handler
- internal/testgen/jsonimport.go - JSON file importer
- internal/testgen/generator.go - Go test code generator
- internal/testgen/interactive_test.go - Unit tests for interactive mode (8 tests)
- internal/testgen/jsonimport_test.go - Unit tests for JSON import (7 tests)
- internal/testgen/generator_test.go - Unit tests for generator (6 tests)

**Modified:**
- None (clean implementation with no modifications to existing files)

### Technical Research Sources

**Go Code Generation:**
- [Package go/ast - The Go Programming Language](https://pkg.go.dev/go/ast)
- [Package go/format - The Go Programming Language](https://pkg.go.dev/go/format)
- [Package go/parser - The Go Programming Language](https://pkg.go.dev/go/parser)

**Test Patterns:**
- [Table-driven tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [testify/assert package](https://pkg.go.dev/github.com/stretchr/testify/assert)

### Previous Story Intelligence (Story 3.3)

**Key Learnings from Watch Mode Implementation:**
- Successfully used third-party library (fsnotify) following architecture patterns
- Command flag patterns established: BoolVarP for flags with short forms
- Test file paths follow pattern: `problems/<slug>_test.go`
- Integration with existing TestService for test execution
- Error handling: stderr for errors, exit codes (0=success, 1=fail, 2=usage, 3=database)
- Unit tests with mocking: used testify/assert for assertions
- Table-driven tests with t.Run() subtests standard pattern

**Files Created in Story 3.3:**
- internal/testing/watcher.go - File watching logic
- internal/testing/watcher_test.go - Unit tests

**Files Modified in Story 3.3:**
- cmd/test.go - Added --watch flag, routing logic
- cmd/test_test.go - Added integration tests
- go.mod/go.sum - Added fsnotify dependency

**Testing Approach from Story 3.3:**
- 8 unit tests covering core functionality
- 6 integration tests for CLI flags
- All tests in internal/testing and cmd packages passed
- Project builds successfully

**Code Patterns to Follow:**
- Use `go get` to add new dependencies
- Create internal package for core logic
- Keep cmd/ files focused on CLI interaction
- Use Edit tool for modifying existing files
- Use Write tool for creating new files
- Run `go test ./...` to verify all tests pass
- Run `go build` to verify compilation
