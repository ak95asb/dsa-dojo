---
stepsCompleted: [1, 2, 3, 4, 5]
inputDocuments:
  - '/Users/noi03_ajaysingh/Documents/LearnGo/dsa/docs/prd.md'
workflowType: 'architecture'
lastStep: 5
project_name: 'dsa'
user_name: 'Empire'
date: '2025-12-10'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**

The **dsa** platform requires 54 capabilities organized into 8 architectural domains:

1. **Problem Management (7 FRs)** - Problem library with topic/difficulty organization, metadata management, workspace initialization
2. **Test-Driven Workflow (6 FRs)** - Scaffolding generation, native Go test execution, solution validation, test history tracking
3. **Progress Tracking & Analytics (8 FRs)** - Completion tracking, solution history, timestamps, streaks (Phase 2), weak area identification (Phase 3), personalized recommendations (Phase 3)
4. **CLI Configuration & Customization (8 FRs)** - Hybrid config (file + env + flags), editor preferences, output verbosity, celebration toggles, environment convention respect
5. **Output & Reporting (6 FRs)** - Dual output modes (human-friendly colored + machine-parseable JSON), progress visualizations, error messages, UNIX stdout/stderr conventions, shell completion (Phase 2)
6. **Motivation & Celebration System (5 FRs - Phase 2)** - Milestone celebrations with ASCII art, streak visualization, programming humor, context-aware encouragement, topic mastery recognition
7. **Data Management (5 FRs)** - Local embedded database storage, zero network dependency, import/export (Phase 4), data integrity guarantees, global + project-specific config
8. **Scripting & Automation (5 FRs)** - Shell composability (pipes/redirects/logical operators), consistent exit codes, non-interactive mode support, TTY detection, flag-based confirmations
9. **Community & Extensibility (4 FRs - Phase 4)** - Problem contribution workflow, multi-language support, contribution documentation, custom problem sets

**Architectural Implications:**
- Command routing layer supporting current + future commands
- Database schema designed for migration path through Phase 2/3/4
- Plugin-style architecture for language support extensibility
- Output abstraction layer supporting dual formatting modes
- Configuration precedence system (flags > env > file > defaults)

**Non-Functional Requirements:**

Six quality attribute categories drive architectural constraints:

**Performance (7 NFRs):**
- CLI cold start <500ms, warm execution <100ms
- Database queries <100ms (non-blocking)
- Scaffolding generation <200ms
- Status dashboard <300ms regardless of history size
- Problem browsing <100ms
- Test execution matches native `go test` performance (zero overhead)

**Architectural Impact:** Requires efficient database design, lazy loading strategies, command caching, and direct test framework integration without wrappers.

**Reliability & Data Integrity (6 NFRs):**
- 100% data integrity across sessions (zero progress loss)
- Complete offline operation (no network dependencies)
- Transactional database operations
- Graceful crash/kill handling without corruption
- Safe config parsing with clear errors
- Automatic recovery for corrupted database

**Architectural Impact:** Requires ACID-compliant local database, transactional boundaries around state changes, write-ahead logging, and checkpoint/recovery mechanisms.

**Integration & Compatibility (6 NFRs):**
- Native Go testing framework integration (no wrapper overhead)
- Idiomatic Go code generation
- Respect $EDITOR environment variable and platform conventions
- NO_COLOR and TTY detection for output formatting
- UNIX conventions (exit codes, stdin/stdout/stderr, signal handling)
- Shell completion for bash/zsh/fish

**Architectural Impact:** Requires deep understanding of `go test` execution model, platform-specific editor detection, signal handler registration, and shell completion protocol integration.

**Portability & Cross-Platform Support (5 NFRs):**
- Unmodified operation on macOS, Linux, Windows
- Platform-specific path handling
- ANSI color rendering on Windows 10+, macOS, Linux
- Single binary distribution <20MB
- No external dependencies beyond Go runtime

**Architectural Impact:** Requires cross-compilation strategy, `filepath` package usage, conditional platform code for editor/terminal detection, static linking, and embedded assets.

**Maintainability & Extensibility (6 NFRs):**
- Language addition without core refactoring
- Community problem contributions without code changes
- Backward-compatible configuration schema
- New command addition without breaking workflows
- Database schema migrations
- Go best practices (passes golint, go vet, staticcheck)

**Architectural Impact:** Requires plugin architecture for language support, declarative problem format, versioned config schema, command registry pattern, migration framework, and linting CI integration.

**Usability & Developer Experience (6 NFRs):**
- Actionable error messages
- Comprehensive `--help` with examples
- Setup <5 minutes (clone to first problem)
- Defaults work for 80% of users
- Intuitive verb-noun command naming
- Clear visual hierarchy in output

**Architectural Impact:** Requires error classification system, help text generation framework, onboarding flow optimization, sensible default configuration, consistent command naming convention, and terminal formatting library.

### Scale & Complexity

**Project Scale Assessment:**

- **Primary domain:** CLI tool / Developer tooling
- **Complexity level:** Low-Medium
  - Single-user application (no multi-tenancy, no auth, no servers)
  - Local-first architecture (no distributed systems concerns)
  - Straightforward data model (problems, solutions, progress, config)
  - Clear boundaries (CLI commands, database, filesystem, editor integration)

- **Estimated architectural components:** 8-10 major components
  1. CLI command router and flag parser
  2. Problem library manager (load/browse/metadata)
  3. Scaffolding generator (template engine for Go code)
  4. Test execution engine (native `go test` integration)
  5. Database layer (schema, queries, migrations)
  6. Configuration manager (file + env + flags with precedence)
  7. Output formatter (terminal rendering + JSON serialization)
  8. Progress tracker (streaks, analytics, recommendations - Phase 2+)
  9. Editor integration (detect, launch, platform-specific)
  10. Shell completion generator (dynamic problem names)

**Complexity Indicators:**
- ✅ Local-only operation (simplifies architecture significantly)
- ✅ Single-user model (no authorization/sharing complexity)
- ✅ Embedded database (no network database configuration)
- ⚠️ Native tool integration (`go test` requires deep understanding)
- ⚠️ Cross-platform support (Windows/macOS/Linux path/terminal differences)
- ⚠️ Multi-phase roadmap (architecture must support Phase 2/3/4 evolution)

### Technical Constraints & Dependencies

**Hard Constraints:**

1. **Go as implementation language** - Required for native `go test` integration and Go developer target audience
2. **Offline-first operation** - No network calls for core functionality (blocks cloud databases, external APIs)
3. **Performance targets** - <500ms cold start rules out heavyweight frameworks
4. **Single binary distribution** - Must embed all assets (templates, problem library, database schema)
5. **Cross-platform support** - Architecture must handle macOS/Linux/Windows differences from day one

**Known Dependencies:**

- **CLI framework:** Cobra (de facto standard for Go CLI tools) or urfave/cli or stdlib only
- **Database:** SQLite (embedded) with choice of Go ORM (GORM, sqlx) or stdlib database/sql
- **Terminal output:** fatih/color or gookit/color for ANSI colors, tablewriter for tables
- **Configuration:** Viper (pairs with Cobra) for config file + env var + flag management
- **Testing framework:** Native Go testing (no additional test framework dependencies)

**Technology Decisions Needed:**
- CLI framework selection (impacts command structure, help generation, completion)
- Database abstraction choice (impacts query performance, migration complexity)
- Terminal library selection (impacts output formatting capabilities)
- Configuration library choice (impacts precedence handling, validation)

### Cross-Cutting Concerns Identified

**Concerns affecting multiple architectural components:**

1. **Error Handling Strategy**
   - Affects: All components
   - Need: Consistent error types, actionable messages, exit code mapping
   - Decision: Error wrapping pattern, custom error types, error classification

2. **Configuration Management**
   - Affects: All components reading settings
   - Need: Precedence order (flags > env > file > defaults), validation, hot-reload
   - Decision: Config struct design, Viper integration, validation framework

3. **Database Schema Evolution**
   - Affects: Data layer, all features touching persistence
   - Need: Migration framework, version tracking, backward compatibility
   - Decision: Migration tool (golang-migrate, goose, custom), schema versioning strategy

4. **Output Formatting**
   - Affects: All commands producing output
   - Need: Consistent styling, color scheme, JSON structure, TTY detection
   - Decision: Output abstraction (Writer interface), formatting library, JSON schema

5. **Testing Strategy**
   - Affects: All components
   - Need: Unit tests, integration tests, CLI command tests, database tests
   - Decision: Test organization, mocking strategy, test data management

6. **Platform Abstraction**
   - Affects: Editor integration, path handling, terminal detection
   - Need: Platform-specific code isolation, graceful fallbacks
   - Decision: Interface boundaries for platform-specific code, build tags vs runtime detection

7. **Logging & Observability**
   - Affects: All components
   - Need: Debug logging, error tracking, performance profiling
   - Decision: Logging library, log levels, verbose mode behavior

## Starter Template Evaluation

### Primary Technology Domain

**CLI Tool / Go** - Based on project requirements analysis, **dsa** is a command-line developer tool implemented in Go.

### Technical Stack Established

**Language & Runtime:**
- **Go** (required for native `go test` integration and target audience)
- **Go modules** for dependency management

**Core Frameworks:**
- **Cobra** - Industry-standard CLI framework (used by Kubernetes, Docker, Hugo, GitHub CLI)
- **Viper** - Configuration management (seamlessly integrates with Cobra)

**Rationale for Cobra + Viper:**
- Cobra handles command structure, flag parsing, help generation, and shell completion out-of-box
- Viper provides configuration precedence (flags > env vars > files > defaults) matching your NFR requirements
- Widely adopted (de facto standard), well-maintained, extensive documentation
- Native support for subcommands (`dsa init`, `dsa solve`, `dsa test`, `dsa status`)
- Built-in completion generation for bash/zsh/fish

### Starter Options Considered

**Option 1: Cobra CLI Generator (Recommended)**
- **Tool:** cobra-cli - Official project scaffolding tool
- **Creates:** Barebones project structure with `main.go`, `cmd/root.go`, boilerplate
- **Benefits:** Fast setup, follows Cobra best practices, generates initial structure
- **Installation:** `go install github.com/spf13/cobra-cli@latest`

**Option 2: Manual Setup with Standard Go Project Layout**
- **Reference:** golang-standards/project-layout
- **Structure:** `cmd/`, `internal/`, `pkg/` directories for organized code
- **Benefits:** More control over initial structure, no generator dependency
- **Trade-off:** More manual setup work

**Option 3: Minimal Setup (Single main.go)**
- **Approach:** Start with single file, refactor as needed
- **Benefits:** Simplest possible start
- **Trade-off:** Doesn't scale well for 54 FRs + 8-10 components

### Selected Approach: Cobra CLI Generator + Standard Layout

**Rationale:**
- Cobra CLI generator provides immediate productivity (working CLI in minutes)
- Standard Go project layout supports growth through Phase 2/3/4
- Matches NFR requirement: "Setup <5 minutes from clone to first problem"
- Establishes patterns for AI agent consistency

### Initialization Commands

**Step 1: Install Cobra CLI Generator**
```bash
go install github.com/spf13/cobra-cli@latest
```

**Step 2: Initialize Go Module**
```bash
mkdir dsa
cd dsa
go mod init github.com/yourusername/dsa
```

**Step 3: Initialize Cobra Project**
```bash
cobra-cli init
```

**Step 4: Add Commands**
```bash
cobra-cli add init
cobra-cli add solve
cobra-cli add test
cobra-cli add status
```

**Step 5: Test Initial Setup**
```bash
go run main.go
```

### Architectural Decisions Provided by Starter

**Project Structure (Standard Go Layout):**
```
dsa/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command + global flags
│   ├── init.go            # dsa init command
│   ├── solve.go           # dsa solve command
│   ├── test.go            # dsa test command
│   └── status.go          # dsa status command
├── internal/              # Private application code
│   ├── database/          # SQLite/ORM layer
│   ├── problems/          # Problem library manager
│   ├── scaffold/          # Code generation engine
│   ├── config/            # Configuration management
│   ├── output/            # Terminal formatting + JSON
│   └── editor/            # Editor integration
├── pkg/                   # Public libraries (if needed)
├── problems/              # Embedded problem definitions
├── main.go                # Application entry point
├── go.mod                 # Go module dependencies
└── go.sum                 # Dependency checksums
```

**Configuration Management (Viper):**
- **Precedence:** Flags > Environment Variables > Config Files > Defaults
- **Config file locations:** `~/.dsa/config.yaml`, `.dsa/config.yaml` (project-specific)
- **Automatic binding:** Viper auto-binds flags to config keys
- **Format support:** YAML, TOML, JSON (recommend YAML for readability)

**Command Structure Pattern:**
```go
// Each command in cmd/ follows this pattern:
var solveCmd = &cobra.Command{
    Use:   "solve [problem]",
    Short: "Start working on a specific problem",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation delegates to internal packages
    },
}
```

**Flag Management:**
```go
// Global flags (root.go)
rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "JSON output")
rootCmd.PersistentFlags().StringVar(&editor, "editor", "", "code editor")

// Command-specific flags
solveCmd.Flags().StringP("difficulty", "d", "all", "filter by difficulty")
```

**Viper Integration:**
```go
// Automatic config file discovery
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath("$HOME/.dsa")
viper.AddConfigPath(".")

// Environment variable support
viper.SetEnvPrefix("DSA")
viper.AutomaticEnv()

// Bind flags to config
viper.BindPFlag("editor", rootCmd.PersistentFlags().Lookup("editor"))
```

**Shell Completion (Built-in):**
```go
// Cobra generates completion for bash/zsh/fish automatically
rootCmd.CompletionOptions.DisableDefaultCmd = false

// Custom completion for problem names (dynamic from database)
solveCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    // Query database for matching problem names
    return problemNames, cobra.ShellCompDirectiveNoFileComp
}
```

**Help Text Generation:**
- Cobra auto-generates `--help` for all commands
- Example usage embedded in command definitions
- Hierarchical help (root help lists subcommands)

**Error Handling Pattern:**
```go
// Cobra provides SilenceErrors and SilenceUsage for clean error messages
rootCmd.SilenceErrors = true
rootCmd.SilenceUsage = true

// Custom error types with exit codes
type CLIError struct {
    Message string
    Code    int
}
```

### Development Experience Features

**Hot Reloading:**
- Use `air` or `CompileDaemon` for Go file watching during development
- Native `go run` for simple development workflow

**Testing Infrastructure:**
- Native Go testing (`go test ./...`)
- Cobra command testing helpers
- Table-driven tests for command flag combinations

**Linting & Formatting:**
- `golangci-lint` for comprehensive linting
- `gofmt` / `goimports` for code formatting
- CI integration with GitHub Actions

**Debugging:**
- Delve debugger integration
- VS Code Go extension support
- Verbose flag for detailed logging

**Note:** Project initialization using Cobra CLI should be documented in the first implementation story, but the actual code implementation will follow after all architectural decisions are made.

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
1. Database & ORM: GORM with SQLite driver (enables data persistence)
2. Testing Strategy: testify/assert + in-memory SQLite (ensures quality)
3. CI/CD Pipeline: GitHub Actions (enables releases)

**Important Decisions (Shape Architecture):**
1. Migration Strategy: GORM AutoMigrate (simplifies schema evolution)
2. Distribution: GoReleaser (professional release automation)
3. Versioning: SemVer with git tags (clear version communication)

**Deferred Decisions (Post-MVP):**
1. Caching: sync.RWMutex pattern if needed later (premature optimization avoided)
2. Advanced migration tools: Add golang-migrate if complex data migrations needed
3. Additional package managers: Homebrew in Phase 2, others as needed

### Data Architecture

**Database: SQLite with GORM ORM**

**Technology Stack:**
- **ORM:** GORM v1.30.1+ (gorm.io/gorm)
- **SQLite Driver:** gorm.io/driver/sqlite v1.6.1+
- **Installation:** `go get -u gorm.io/gorm gorm.io/driver/sqlite`

**Rationale:**
- Auto-migrations simplify Phase 2/3/4 schema evolution (streaks, analytics, multi-language)
- Developer ergonomics accelerate MVP delivery (solo/2-person team, 4-6 week timeline)
- Performance overhead negligible for single-user local CLI (<1ms query time vs <0.5ms)
- Large community and excellent documentation reduce debugging time
- Meets all NFRs: <100ms queries (achieves <1ms), ACID compliance, zero data loss

**Schema Design:**
```go
// Phase 1 MVP
type Problem struct {
    ID         uint   `gorm:"primaryKey"`
    Slug       string `gorm:"uniqueIndex"`
    Title      string
    Difficulty string // "easy", "medium", "hard"
    Topic      string // "arrays", "linked-lists", "trees", "sorting"
    Description string
    CreatedAt  time.Time
}

type Solution struct {
    ID         uint `gorm:"primaryKey"`
    ProblemID  uint `gorm:"index"`
    Code       string
    Language   string // "go" (later: "rust", "python", "typescript")
    Passed     bool
    CreatedAt  time.Time
}

type Progress struct {
    ID          uint `gorm:"primaryKey"`
    ProblemID   uint `gorm:"uniqueIndex"`
    Status      string // "not_started", "in_progress", "completed"
    Attempts    int
    LastAttempt time.Time
    // Phase 2 additions: FirstCompletedAt, StreakData
    // Phase 3 additions: TimeSpent, SuccessRate, ReviewSchedule
}
```

**Migration Strategy: GORM AutoMigrate**

**Approach:**
```go
// On application startup
db.AutoMigrate(&Problem{}, &Solution{}, &Progress{})
```

**Rationale:**
- Schema evolution is purely additive (new columns for Phase 2/3/4 features)
- Single-user local database = lower risk than production multi-tenant API
- GORM AutoMigrate handles column additions automatically
- Can add golang-migrate later if data migrations become necessary

**Phase Evolution Plan:**
- **Phase 1:** Base schema (problems, solutions, progress)
- **Phase 2:** Add `FirstCompletedAt time.Time`, `StreakData string` to Progress
- **Phase 3:** Add `TimeSpent int`, `SuccessRate float64`, `ReviewSchedule time.Time`
- **Phase 4:** Add `Language string` support, language-specific test data

**Caching Strategy: Deferred**

**Decision:** No explicit caching in MVP

**Rationale:**
- All data already local in SQLite (database IS the cache)
- SQLite queries <1ms for expected data volumes
- Problem library small (10-15 problems in MVP)
- Viper already caches configuration
- Adding Redis/in-memory cache adds complexity without benefit

**Future Option:** If needed, use simple Go map with `sync.RWMutex` for problem metadata:
```go
type ProblemCache struct {
    sync.RWMutex
    problems map[string]*Problem
}
```

### Testing Strategy

**Multi-Layer Testing Approach**

**Layer 1: Unit Tests (internal packages)**
- Test individual components in isolation
- Database layer, problem manager, scaffolding generator, output formatter
- Table-driven tests for comprehensive coverage

**Layer 2: Integration Tests (component interaction)**
- Test database operations, file I/O, editor integration
- Use in-memory SQLite for isolated, fast tests

**Layer 3: CLI Command Tests (end-to-end)**
- Test actual CLI commands as users invoke them
- Cobra command testing helpers + golden files for output validation

**Layer 4: Cross-Platform Tests**
- GitHub Actions matrix builds (Windows, macOS, Linux)
- Validate platform-specific code (paths, editor detection, terminal)

**Testing Tools & Frameworks:**

**Assertion Library: testify/assert**
```go
import "github.com/stretchr/testify/assert"

func TestProblemCreation(t *testing.T) {
    problem := &Problem{Title: "Two Sum", Difficulty: "easy"}
    assert.Equal(t, "easy", problem.Difficulty)
    assert.NotEmpty(t, problem.Title)
}
```

**Rationale:** Cleaner, more readable tests with better failure messages worth the minimal dependency

**Test Database Strategy:**

**In-memory SQLite for unit/integration tests:**
```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(t, err)
    db.AutoMigrate(&Problem{}, &Solution{}, &Progress{})
    return db
}
```

**Temporary file database for CLI command tests:**
```go
func setupCLITestDB(t *testing.T) string {
    tmpfile, _ := os.CreateTemp("", "dsa-test-*.db")
    t.Cleanup(func() { os.Remove(tmpfile.Name()) })
    return tmpfile.Name()
}
```

**Rationale:** In-memory is fast and isolated for most tests; temp file tests actual file I/O edge cases

**Coverage Targets:**
- **Overall project:** 70%+ coverage
- **Critical paths:** 80%+ coverage (database layer, progress tracking, test execution engine)
- **CLI commands:** Golden file tests for each command
- **Continuous monitoring:** Coverage reports in CI/CD

**Test Organization:**
```
dsa/
├── internal/
│   ├── database/
│   │   ├── models.go
│   │   └── models_test.go        # Unit tests
│   ├── problems/
│   │   ├── manager.go
│   │   └── manager_test.go       # Integration tests
│   └── scaffold/
│       ├── generator.go
│       └── generator_test.go     # Unit tests
├── cmd/
│   └── commands_test.go          # CLI command tests
└── testdata/                     # Golden files, fixtures
    ├── golden/
    │   ├── status_output.txt
    │   └── init_output.txt
    └── fixtures/
        └── sample_problems.json
```

### Infrastructure & Deployment

**CI/CD Pipeline: GitHub Actions**

**Workflow Configuration:**
```yaml
# .github/workflows/ci.yml
name: CI/CD
on: [push, pull_request]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.23']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Run tests
        run: go test -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4

      - name: Run linters
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  release:
    needs: test
    if: startsWith(github.ref, 'refs/tags/v')
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Rationale:**
- Free for open source projects
- Matrix builds ensure cross-platform compatibility
- Integrated with GitHub (where code lives)
- Standard for Go open source projects

**Binary Distribution: GoReleaser**

**Configuration:**
```yaml
# .goreleaser.yml
project_name: dsa

builds:
  - env: [CGO_ENABLED=0]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    files:
      - LICENSE
      - README.md
      - CHANGELOG.md

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^chore:'

brews:
  - repository:
      owner: yourusername
      name: homebrew-tap
    homepage: https://github.com/yourusername/dsa
    description: CLI-based DSA practice platform for Go developers
    license: MIT
```

**Installation Methods:**

**Phase 1 (MVP): Direct Binary Download**
```bash
# Linux/macOS
curl -L https://github.com/user/dsa/releases/latest/download/dsa_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv dsa /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/user/dsa/releases/latest/download/dsa_Windows_x86_64.zip" -OutFile "dsa.zip"
```

**Phase 2: Homebrew Tap** (auto-generated by GoReleaser)
```bash
brew install yourusername/tap/dsa
```

**Future: Additional Methods**
- `go install github.com/yourusername/dsa@latest`
- Package managers: apt, yum, chocolatey, scoop
- Docker image (if containerized workflow requested)

**Versioning Strategy: Semantic Versioning (SemVer)**

**Version Scheme:**
- **v1.0.0** - Phase 1 MVP (core CLI workflow: init, solve, test, status)
- **v1.1.0** - Phase 2 (celebration system, streaks, enhanced status)
- **v1.2.0** - Phase 3 (smart learning, spaced repetition, analytics)
- **v2.0.0** - Phase 4 (multi-language support - breaking change to config schema)

**Release Process:**
```bash
# 1. Update CHANGELOG.md
# 2. Update version in code if needed
# 3. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0: MVP with core CLI workflow"
git push origin v1.0.0

# 4. GitHub Actions + GoReleaser automatically:
#    - Runs all tests and linting
#    - Builds binaries for all platforms
#    - Creates GitHub release with artifacts
#    - Updates Homebrew tap
```

**Pre-Release Quality Gates:**

**Automated Checks (CI/CD enforced):**
1. ✅ All tests pass (`go test ./...` across all platforms)
2. ✅ Race detector clean (`go test -race`)
3. ✅ Coverage meets threshold (70%+ overall, 80%+ critical)
4. ✅ Linting passes (`golangci-lint run`)
5. ✅ Code formatting (`gofmt -s`, `goimports`)
6. ✅ Static analysis (`go vet`, `staticcheck`)
7. ✅ Build succeeds for all target platforms
8. ✅ Binary size <20MB (NFR23)

**Manual Checks (pre-tag):**
- CHANGELOG.md updated with release notes
- README.md updated if CLI commands changed
- Version number follows SemVer conventions
- Breaking changes documented if major version bump

### Decision Impact Analysis

**Implementation Sequence:**

**Phase 1 - Foundation:**
1. Initialize Cobra CLI project structure
2. Set up GORM with SQLite, implement base models
3. Implement core commands (init, solve, test, status)
4. Set up GitHub Actions CI/CD
5. Configure GoReleaser for releases

**Phase 2 - Quality & Distribution:**
1. Achieve 70%+ test coverage
2. Set up golangci-lint configuration
3. Create first release (v1.0.0)
4. Set up Homebrew tap

**Phase 3 - Feature Evolution:**
1. Phase 2 features: Add streak/celebration models (AutoMigrate handles schema)
2. Phase 3 features: Add analytics models (AutoMigrate handles schema)
3. Phase 4 features: Multi-language support (major version bump)

**Cross-Component Dependencies:**

1. **GORM AutoMigrate → All Data Access**
   - All components using database depend on GORM being initialized first
   - Migration runs on app startup before any commands execute

2. **Cobra Commands → Testing Strategy**
   - CLI command tests depend on Cobra's testing helpers
   - Golden files validate command output format

3. **GitHub Actions → GoReleaser**
   - Release workflow triggers only after tests pass
   - Tag creation requires clean CI/CD pipeline

4. **Viper Config → All Commands**
   - All commands read config through Viper singleton
   - Config precedence (flags > env > file) affects all command behavior

5. **testify/assert → All Tests**
   - Assertion library used across all test layers
   - Consistent test style across codebase

## Implementation Patterns & Consistency Rules

### Pattern Categories Defined

**Critical Conflict Points Identified:** 5 major categories with 25+ specific patterns to prevent AI agent conflicts

### Naming Patterns

**Go Code Naming Conventions (Follow Go standards):**

- **Exported functions:** `PascalCase` - `GetProblem()`, `CreateSolution()`, `ValidateInput()`
- **Unexported functions:** `camelCase` - `validateInput()`, `parseConfig()`, `formatOutput()`
- **Variables:** `camelCase` - `problemID`, `userName`, `configPath`
- **Constants:** `PascalCase` for exported, `camelCase` for unexported
  ```go
  const (
      DefaultTimeout = 30 * time.Second  // Exported
      maxRetries     = 3                 // Unexported
  )
  ```
- **Interfaces:** `-er` suffix when applicable - `ProblemManager`, `ConfigLoader`, `OutputFormatter`
- **Type names:** `PascalCase` - `Problem`, `Solution`, `Progress`, `CLIError`

**File & Package Naming:**

- **Files:** `snake_case.go` - `problem_manager.go`, `config_loader.go`, `output_formatter.go`
- **Test files:** Co-located `*_test.go` - `problem_manager_test.go`, `database_test.go`
- **Packages:** Short, singular nouns, lowercase - `database`, `problem`, `config`, `output`, `editor`
- **No underscores in package names:** Use `scaffold` not `code_scaffold`, `testdata` not `test_data`

**Database Naming Conventions (GORM):**

- **Tables:** Plural, snake_case - `problems`, `solutions`, `progress`
- **Columns:** snake_case - `user_id`, `created_at`, `problem_slug`, `difficulty_level`
- **Foreign keys:** Explicit naming - `problem_id` (references `problems.id`)
- **Indexes:** Prefix with `idx_` - `idx_problems_slug`, `idx_progress_problem_id`
- **Struct tags:** Match column names exactly
  ```go
  type Problem struct {
      ID         uint      `gorm:"primaryKey" json:"id"`
      Slug       string    `gorm:"uniqueIndex:idx_problems_slug" json:"slug"`
      Title      string    `gorm:"not null" json:"title"`
      Difficulty string    `gorm:"type:varchar(20)" json:"difficulty"`
      CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
  }
  ```

### Structure Patterns

**Project Organization (Standard Go Layout):**

```
dsa/
├── cmd/                    # Cobra commands
│   ├── root.go            # Root command + global flags
│   ├── init.go            # dsa init
│   ├── solve.go           # dsa solve
│   ├── test.go            # dsa test
│   └── status.go          # dsa status
├── internal/              # Private application packages
│   ├── database/          # GORM models, migrations, connection
│   ├── problem/           # Problem library management
│   ├── scaffold/          # Code generation engine
│   ├── config/            # Viper configuration wrapper
│   ├── output/            # Terminal formatting (colors, tables, JSON)
│   └── editor/            # Editor detection and launching
├── pkg/                   # Public libraries (if needed)
├── problems/              # Embedded problem definitions (YAML/JSON)
├── testdata/              # Test fixtures and golden files
│   ├── golden/            # Expected CLI outputs
│   │   ├── status.txt
│   │   └── init.txt
│   └── fixtures/          # Test data
│       └── problems.json
└── main.go                # Application entry point
```

**Package Organization Rules:**

1. **By feature:** Each `internal/` package represents distinct capability
2. **Single responsibility:** One clear purpose per package
3. **Minimize cross-package deps:** Avoid circular imports (use interfaces if needed)
4. **Public interfaces in pkg/:** Only if code needs external import

**Test Organization:**

- **Co-located tests:** `problem_manager.go` → `problem_manager_test.go` in same directory
- **Table-driven tests:** Use `t.Run()` with subtests
  ```go
  func TestDifficultyValidation(t *testing.T) {
      tests := []struct {
          name    string
          input   string
          wantErr bool
      }{
          {"valid easy", "easy", false},
          {"valid medium", "medium", false},
          {"invalid", "super-hard", true},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              err := ValidateDifficulty(tt.input)
              if (err != nil) != tt.wantErr {
                  t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
              }
          })
      }
  }
  ```
- **Test fixtures:** Shared data in `testdata/` directory
- **Golden files:** Expected CLI output in `testdata/golden/`

**Error Types Location:**

- **Package-specific errors:** Define in package where used
  ```go
  // internal/problem/errors.go
  var (
      ErrProblemNotFound = errors.New("problem not found")
      ErrInvalidSlug     = errors.New("invalid problem slug")
  )
  ```
- **Common errors:** Define in `internal/errors/` if shared across packages
- **CLI exit codes:** Define in `cmd/root.go` as constants

**Constants & Enums:**

- **Package-level constants:** Top of file where heavily used
- **Shared constants:** `internal/constants/` if needed across packages
- **Enum pattern:** Typed constants with iota
  ```go
  type Difficulty int

  const (
      DifficultyEasy Difficulty = iota
      DifficultyMedium
      DifficultyHard
  )

  func (d Difficulty) String() string {
      return [...]string{"easy", "medium", "hard"}[d]
  }
  ```

### Format Patterns

**JSON Output Format (--json flag):**

**Success Response:**
```json
{
  "success": true,
  "data": {
    "problem_id": 1,
    "slug": "two-sum",
    "title": "Two Sum",
    "difficulty": "easy",
    "status": "completed"
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": {
    "message": "Problem 'invalid-slug' not found",
    "code": "PROBLEM_NOT_FOUND"
  }
}
```

**JSON Field Naming:**

- **Convention:** `snake_case` (matches database, Go JSON standard)
- **All JSON output:** Consistent snake_case
- **Struct tags:** `json:"problem_id"`, `json:"created_at"`, `json:"first_name"`

**Date/Time Formats:**

- **JSON output:** RFC3339 (`2025-12-10T15:04:05Z`)
  ```go
  json:"created_at" // time.Time automatically marshals to RFC3339
  ```
- **Database storage:** GORM `time.Time` (handles automatically)
- **CLI display:** Human-friendly
  ```go
  createdAt.Format("Jan 2, 2006 3:04 PM")  // "Dec 10, 2025 3:04 PM"
  ```

**CLI Exit Codes:**

```go
const (
    ExitSuccess          = 0  // Command succeeded
    ExitGeneralError     = 1  // General/unknown error
    ExitUsageError       = 2  // Invalid command/flags
    ExitDatabaseError    = 3  // Database operation failed
    ExitTestFailure      = 4  // go test execution failed
)
```

**Boolean & Null Handling:**

- **Booleans:** `true`/`false` in Go, JSON (not 1/0 or strings)
- **Prefer zero values:** Use `""` or `0` instead of null when possible
- **Optional fields:** Use pointers `*string`, `*int` only when truly optional
- **JSON omit:** Use `json:"field,omitempty"` to exclude zero values

### Communication Patterns

**Logging Strategy:**

**Log Levels (using log/slog):**
- **ERROR:** Unrecoverable errors stopping execution (always shown)
- **WARN:** Recoverable issues, degraded functionality (normal mode)
- **INFO:** Important events, milestones (normal mode)
- **DEBUG:** Detailed diagnostics (--verbose flag only)

**Logging Format:**
```go
import "log/slog"

// Structured logging
slog.Info("problem loaded",
    "problem_id", problemID,
    "difficulty", difficulty,
    "topic", topic)

// Error logging with context
slog.Error("database query failed",
    "error", err,
    "query", "SELECT * FROM problems")
```

**Logging Configuration:**
```go
// Set level based on --verbose flag
var logLevel slog.Level = slog.LevelInfo
if verbose {
    logLevel = slog.LevelDebug
}
logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: logLevel,
}))
slog.SetDefault(logger)
```

**Output Destinations:**

- **Command output:** Stdout (pipeable/redirectable)
- **Logs/diagnostics:** Stderr (doesn't pollute command output)
- **Interactive prompts:** Stderr (doesn't interfere with scripting)
- **Rule:** Never mix command results and diagnostics on same stream

**Progress Indicators:**

- **TTY detection:** Show spinners/bars only if stdout is TTY
  ```go
  if isatty.IsTerminal(os.Stdout.Fd()) {
      // Show spinner
  }
  ```
- **Non-TTY mode:** Silent or simple text updates
- **JSON mode:** No progress indicators, only final JSON output

**User-Facing Messages:**

- **Errors:** Clear, actionable
  - ❌ "Error: file not found"
  - ✅ "Problem 'two-sum' not found. Run 'dsa list' to see available problems."
- **Success:** Concise confirmation
  - ✅ "Workspace initialized at /Users/you/dsa-workspace"
- **Warnings:** Helpful context
  - ⚠️  "No active problem. Use 'dsa solve <problem>' to start practicing."

### Process Patterns

**Error Handling Strategy:**

**Core Principles:**
1. **Return errors, don't panic** (panic only for programmer errors)
2. **Wrap errors with context** using `fmt.Errorf("context: %w", err)`
3. **Define sentinel errors** for expected conditions

**Sentinel Errors:**
```go
// internal/problem/errors.go
var (
    ErrProblemNotFound   = errors.New("problem not found")
    ErrInvalidDifficulty = errors.New("invalid difficulty level")
    ErrSlugTaken         = errors.New("problem slug already exists")
)
```

**Error Wrapping Pattern:**
```go
// Good: Wrap with context, preserve original
if err := db.Create(&problem).Error; err != nil {
    return fmt.Errorf("failed to create problem '%s': %w", problem.Slug, err)
}

// Check for specific errors
if errors.Is(err, problem.ErrNotFound) {
    return fmt.Errorf("problem '%s' not found", slug)
}

// Type assertions for custom errors
var gormErr *gorm.Error
if errors.As(err, &gormErr) {
    // Handle GORM-specific error
}
```

**CLI Layer Error Handling:**
```go
// In cmd/ files: Convert internal errors to user-friendly messages
func runSolveCommand(cmd *cobra.Command, args []string) error {
    problem, err := problemService.Get(args[0])
    if err != nil {
        if errors.Is(err, problem.ErrNotFound) {
            fmt.Fprintf(os.Stderr, "Problem '%s' not found. Run 'dsa list' to see available problems.\n", args[0])
            os.Exit(ExitUsageError)
        }
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        os.Exit(ExitGeneralError)
    }
    return nil
}
```

**Database Transaction Patterns:**

**When to Use Transactions:**
- Multiple related writes (problem + progress)
- Conditional updates that must be atomic
- Any operation requiring rollback on failure

**Transaction Pattern:**
```go
func (s *Service) CreateProblemWithProgress(problem *Problem) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create problem
        if err := tx.Create(&problem).Error; err != nil {
            return fmt.Errorf("create problem: %w", err)
        }

        // Create initial progress record
        progress := &Progress{
            ProblemID: problem.ID,
            Status:    "not_started",
            Attempts:  0,
        }
        if err := tx.Create(&progress).Error; err != nil {
            return fmt.Errorf("create progress: %w", err)
        }

        return nil // Commit on success, rollback on error
    })
}
```

**Configuration Loading:**

- **Load at startup:** Initialize Viper in `cmd/root.go` `init()` function
- **Precedence enforcement:** Flags > Env vars > Config file > Defaults
- **Validation:** Validate config early, fail fast with clear errors
  ```go
  func initConfig() {
      viper.SetConfigName("config")
      viper.SetConfigType("yaml")
      viper.AddConfigPath("$HOME/.dsa")
      viper.AddConfigPath(".")

      viper.SetEnvPrefix("DSA")
      viper.AutomaticEnv()

      if err := viper.ReadInConfig(); err != nil {
          if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
              // Config file found but error reading
              log.Fatal(err)
          }
          // No config file is okay, use defaults
      }
  }
  ```

**Resource Cleanup:**

- **Defer cleanup:** Always use `defer` for resource cleanup
- **Database connection:** Open once at startup, close in main()
  ```go
  func main() {
      db, err := database.Initialize()
      if err != nil {
          log.Fatalf("failed to initialize database: %v", err)
      }
      defer func() {
          sqlDB, _ := db.DB()
          sqlDB.Close()
      }()

      if err := cmd.Execute(); err != nil {
          os.Exit(ExitGeneralError)
      }
  }
  ```

**Graceful Shutdown:**

- **Signal handling:** Catch SIGINT/SIGTERM for graceful cleanup
  ```go
  ctx, stop := signal.NotifyContext(context.Background(),
      os.Interrupt, syscall.SIGTERM)
  defer stop()

  // Use ctx in long-running operations
  ```

### Enforcement Guidelines

**All AI Agents MUST:**

1. **Follow Go conventions:** Use `gofmt`, `goimports`, pass `golangci-lint`
2. **Use snake_case for:** Database columns, JSON fields, file names
3. **Use PascalCase for:** Exported Go types, functions, constants
4. **Co-locate tests:** `*_test.go` files next to implementation
5. **Wrap errors:** Always add context with `fmt.Errorf("context: %w", err)`
6. **Return errors:** Never panic except for programmer errors
7. **Use GORM AutoMigrate:** For schema changes (defined models automatically)
8. **Table-driven tests:** Use `t.Run()` for multiple test cases
9. **Separate output streams:** Stdout for results, Stderr for logs/errors
10. **Respect exit codes:** Use defined constants for CLI exit status

**Pattern Verification:**

- **Pre-commit checks:** Run `gofmt`, `goimports`, `golangci-lint`
- **CI/CD enforcement:** GitHub Actions runs linting, tests, formatting checks
- **Code review:** Verify patterns during PR review
- **Test coverage:** Maintain 70%+ overall, 80%+ critical paths

**Pattern Documentation:**

- **This document:** Canonical source for implementation patterns
- **Code comments:** Reference patterns in complex code
- **PR descriptions:** Note any pattern deviations and rationale

**Updating Patterns:**

- **Pattern evolution:** Update this doc if patterns need refinement
- **Breaking changes:** Discuss with team before changing established patterns
- **New patterns:** Add to this document when new conflict points discovered

### Pattern Examples

**Good Examples:**

**1. Error Handling with Context:**
```go
// ✅ Good: Wrapped error with context
func (s *ProblemService) GetBySlug(slug string) (*Problem, error) {
    var problem Problem
    err := s.db.Where("slug = ?", slug).First(&problem).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrProblemNotFound
        }
        return nil, fmt.Errorf("query problem by slug '%s': %w", slug, err)
    }
    return &problem, nil
}
```

**2. Table-Driven Test:**
```go
// ✅ Good: Comprehensive table-driven test
func TestValidateDifficulty(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr error
    }{
        {"easy", "easy", nil},
        {"medium", "medium", nil},
        {"hard", "hard", nil},
        {"invalid", "expert", ErrInvalidDifficulty},
        {"empty", "", ErrInvalidDifficulty},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateDifficulty(tt.input)
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("got %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

**3. Consistent Naming:**
```go
// ✅ Good: Follows all naming conventions
type Problem struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    Slug       string    `gorm:"uniqueIndex" json:"slug"`
    Title      string    `json:"title"`
    Difficulty string    `json:"difficulty"`
    CreatedAt  time.Time `json:"created_at"`
}

func (s *ProblemService) CreateProblem(p *Problem) error {
    return s.db.Create(p).Error
}
```

**Anti-Patterns (Avoid These):**

**❌ Bad: Panic instead of returning error**
```go
// ❌ Bad
func GetProblem(id uint) *Problem {
    var p Problem
    if err := db.First(&p, id).Error; err != nil {
        panic(err) // DON'T DO THIS
    }
    return &p
}

// ✅ Good
func GetProblem(id uint) (*Problem, error) {
    var p Problem
    err := db.First(&p, id).Error
    return &p, err
}
```

**❌ Bad: Inconsistent naming**
```go
// ❌ Bad: Mixed conventions
type Problem struct {
    ID         uint   `json:"id"`           // Inconsistent
    problemSlug string `json:"ProblemSlug"` // Wrong on multiple levels
    Title      string `json:"title"`
}
```

**❌ Bad: Lost error context**
```go
// ❌ Bad: No context
if err := doSomething(); err != nil {
    return err // What failed? Where? Why?
}

// ✅ Good: Wrapped with context
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to process problem: %w", err)
}
```

**❌ Bad: Tests without subtests**
```go
// ❌ Bad: One giant test function
func TestAllValidations(t *testing.T) {
    // Test 100 different things sequentially
    // If first fails, rest don't run
}

// ✅ Good: Table-driven with subtests
func TestValidations(t *testing.T) {
    tests := []struct{ ... }{ ... }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) { ... })
    }
}
```
