---
stepsCompleted: [1, 2, 3, 4]
lastStep: 4
workflowCompleted: true
inputDocuments:
  - '/Users/noi03_ajaysingh/Documents/LearnGo/dsa/docs/prd.md'
  - '/Users/noi03_ajaysingh/Documents/LearnGo/dsa/docs/architecture.md'
---

# dsa - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for dsa, decomposing the requirements from the PRD and Architecture into implementable stories.

## Requirements Inventory

### Functional Requirements

**Problem Management (7 FRs):**
- **FR1:** Developer can initialize a DSA practice workspace in their local environment
- **FR2:** Developer can browse available problems by topic (arrays, linked lists, trees, sorting)
- **FR3:** Developer can browse available problems by difficulty (easy, medium, hard)
- **FR4:** Developer can start working on a specific problem by name or identifier
- **FR5:** System provides problem description, constraints, and examples when developer starts a problem
- **FR6:** System provides pre-written test cases for every problem
- **FR7:** Developer can view problem metadata (difficulty, topic, tags, description)

**Test-Driven Workflow (6 FRs):**
- **FR8:** System generates scaffolded solution file with boilerplate code when developer starts a problem
- **FR9:** Developer can run tests against their solution using native Go testing tools
- **FR10:** System opens problem files in developer's preferred code editor automatically
- **FR11:** System validates solution correctness using pre-written test cases
- **FR12:** Developer receives clear test results showing passed/failed test cases
- **FR13:** System tracks test execution history for each problem attempt

**Progress Tracking & Analytics (8 FRs):**
- **FR14:** Developer can view their overall progress (total problems solved, by difficulty, by topic)
- **FR15:** System tracks completion status for each problem (not started, in progress, completed)
- **FR16:** Developer can view their solution history for previously attempted problems
- **FR17:** System records timestamps for problem attempts and completions
- **FR18:** System tracks daily practice streaks (consecutive days with at least one problem solved) [Phase 2]
- **FR19:** System identifies weak areas based on problem-solving history and patterns [Phase 3]
- **FR20:** Developer can view personalized problem recommendations based on practice history [Phase 3]
- **FR21:** System tracks time spent on each problem for performance analytics [Phase 3]

**CLI Configuration & Customization (8 FRs):**
- **FR22:** Developer can configure default code editor via configuration file
- **FR23:** Developer can configure default code editor via environment variables
- **FR24:** Developer can override configuration settings via command-line flags
- **FR25:** Developer can filter problems by difficulty preference in configuration
- **FR26:** Developer can set output verbosity preferences (quiet, normal, verbose)
- **FR27:** Developer can enable or disable celebration features in configuration
- **FR28:** Developer can enable or disable programming humor in output messages
- **FR29:** System respects standard environment conventions (NO_COLOR, TTY detection)

**Output & Reporting (6 FRs):**
- **FR30:** System provides human-friendly terminal output with colors and formatting by default
- **FR31:** Developer can request machine-parseable output (JSON format) for scripting
- **FR32:** System displays progress visualizations (progress bars, ASCII charts) for status commands
- **FR33:** System provides clear error messages when commands fail
- **FR34:** System outputs structured data to stdout and errors to stderr following UNIX conventions
- **FR35:** Developer can use shell completion for commands and problem names [Phase 2]

**Motivation & Celebration System (5 FRs - Phase 2):**
- **FR36:** System celebrates milestone achievements with ASCII art and encouraging messages [Phase 2]
- **FR37:** System displays streak counter with visual indicators in status output [Phase 2]
- **FR38:** System provides encouraging messages with programming humor when problems are solved [Phase 2]
- **FR39:** System recognizes specific milestones (first problem, first hard problem, streak milestones) [Phase 2]
- **FR40:** System displays topic completion progress and celebrates topic mastery [Phase 2]

**Data Management (5 FRs):**
- **FR41:** System stores all progress data locally using embedded database (no network dependency)
- **FR42:** Developer can export their progress data in portable formats (JSON/CSV) [Phase 4]
- **FR43:** Developer can import previously exported progress data [Phase 4]
- **FR44:** System maintains data integrity across sessions with zero progress loss
- **FR45:** System supports both global and project-specific configuration

**Scripting & Automation (5 FRs):**
- **FR46:** Commands can be composed and chained using shell operators (pipes, redirects, logical operators)
- **FR47:** System provides consistent exit codes for success and different error types
- **FR48:** Commands work in non-interactive mode (CI/CD, scripts) without requiring user input
- **FR49:** System detects TTY vs non-TTY contexts and adjusts output accordingly
- **FR50:** Developer can force decisions via flags to skip confirmations in automated contexts

**Community & Extensibility (4 FRs - Phase 4):**
- **FR51:** Community contributor can submit new problem definitions following standardized format [Phase 4]
- **FR52:** Community contributor can add support for additional programming languages [Phase 4]
- **FR53:** System provides clear documentation for contribution workflows [Phase 4]
- **FR54:** Developer can create custom problem sets and themed collections [Phase 4]

### Non-Functional Requirements

**Performance (7 NFRs):**
- **NFR1:** CLI commands execute with cold start time <500ms
- **NFR2:** Warm command execution (after first run) completes in <100ms
- **NFR3:** Database query operations complete without blocking user interaction (<100ms)
- **NFR4:** File scaffolding and generation operations complete in <200ms
- **NFR5:** Status dashboard rendering completes in <300ms regardless of solution history size
- **NFR6:** Problem library browsing operations return results in <100ms
- **NFR7:** Test execution performance matches native `go test` performance (zero overhead)

**Reliability & Data Integrity (6 NFRs):**
- **NFR8:** System maintains 100% data integrity across sessions with zero progress loss
- **NFR9:** System operates completely offline with no network dependencies for core functionality
- **NFR10:** Database operations use transactions to ensure atomic updates
- **NFR11:** System gracefully handles interrupted operations (crash, kill) without data corruption
- **NFR12:** Configuration file parsing fails safely with clear error messages on invalid syntax
- **NFR13:** System provides automatic recovery mechanisms for corrupted local database

**Integration & Compatibility (6 NFRs):**
- **NFR14:** System integrates with native Go testing framework without wrapper overhead
- **NFR15:** Generated code follows idiomatic Go patterns and conventions
- **NFR16:** System respects user's `$EDITOR` environment variable and standard editor detection conventions
- **NFR17:** CLI output respects `NO_COLOR` environment variable and TTY detection for appropriate formatting
- **NFR18:** System follows UNIX conventions for exit codes, stdin/stdout/stderr, and signal handling
- **NFR19:** Shell completion integrates seamlessly with bash, zsh, and fish completion systems

**Portability & Cross-Platform Support (5 NFRs):**
- **NFR20:** System runs without modification on macOS, Linux, and Windows operating systems
- **NFR21:** System handles platform-specific path conventions correctly across all platforms
- **NFR22:** ANSI color output renders correctly on Windows 10+, macOS, and Linux terminals
- **NFR23:** Binary distribution size remains <20MB for single-binary CLI distribution
- **NFR24:** System requires no external dependencies beyond Go runtime for core functionality

**Maintainability & Extensibility (6 NFRs):**
- **NFR25:** Codebase architecture supports addition of new programming languages without core refactoring
- **NFR26:** Problem definition format allows community contributions without code changes
- **NFR27:** Configuration schema supports backward compatibility when adding new settings
- **NFR28:** CLI command structure supports addition of new commands without breaking existing workflows
- **NFR29:** Database schema supports migrations for future feature additions
- **NFR30:** Code follows Go best practices and passes standard linters (golint, go vet, staticcheck)

**Usability & Developer Experience (6 NFRs):**
- **NFR31:** Error messages provide actionable guidance rather than cryptic technical errors
- **NFR32:** Help text (`--help`) for each command is comprehensive and includes usage examples
- **NFR33:** Setup process from clone to first problem takes <5 minutes
- **NFR34:** Configuration defaults work for 80% of users without requiring customization
- **NFR35:** Command naming follows intuitive verb-noun conventions recognizable to CLI users
- **NFR36:** Output formatting provides clear visual hierarchy with appropriate use of color and spacing

### Additional Requirements

**From Architecture - Starter Template:**
- **Project initialization using Cobra CLI generator** (Epic 1, Story 1)
  - Install cobra-cli: `go install github.com/spf13/cobra-cli@latest`
  - Initialize Go module: `go mod init github.com/yourusername/dsa`
  - Initialize Cobra project: `cobra-cli init`
  - Add commands: `cobra-cli add init`, `cobra-cli add solve`, `cobra-cli add test`, `cobra-cli add status`
  - Standard Go project layout (cmd/, internal/, pkg/)

**From Architecture - Database & Migrations:**
- Use GORM v1.30.1+ with SQLite driver (gorm.io/driver/sqlite v1.6.1+)
- Database migrations via GORM AutoMigrate (no separate migration tool needed for Phase 1)
- Schema design with Problem, Solution, Progress models
- Database initialization on application startup

**From Architecture - Testing Framework:**
- Use testify/assert for test assertions
- In-memory SQLite (`:memory:`) for unit and integration tests
- Temporary file database for CLI command tests
- Table-driven tests with `t.Run()` subtests
- Coverage targets: 70%+ overall, 80%+ critical paths
- Golden files for CLI output validation in `testdata/golden/`

**From Architecture - CI/CD & Distribution:**
- GitHub Actions CI/CD with matrix builds (ubuntu-latest, macos-latest, windows-latest)
- Go version: 1.23
- CI pipeline: tests with race detector, coverage upload (codecov), linting (golangci-lint)
- GoReleaser for automated binary distribution
- Release workflow triggers on git tags (v1.0.0, v1.1.0, etc.)
- Homebrew tap auto-generation (Phase 2)
- Semantic versioning: v1.0.0 (MVP), v1.1.0 (Phase 2), v1.2.0 (Phase 3), v2.0.0 (Phase 4)

**From Architecture - Implementation Patterns:**
- Go naming conventions: PascalCase (exported), camelCase (unexported), snake_case (files/DB)
- Package organization: By feature in internal/ (database, problem, scaffold, config, output, editor)
- Error handling: Return errors (not panic), wrap with context using `fmt.Errorf("context: %w", err)`
- Logging: stdlib log/slog with structured logging (ERROR, WARN, INFO, DEBUG levels)
- Output streams: Stdout for results, Stderr for logs/errors
- CLI exit codes: 0 (success), 1 (general error), 2 (usage error), 3 (database error), 4 (test failure)

### FR Coverage Map

**Epic 1: Project Foundation & Setup**
- FR1: Developer can initialize a DSA practice workspace

**Epic 2: Problem Library & Discovery**
- FR2: Developer can browse available problems by topic
- FR3: Developer can browse available problems by difficulty
- FR4: Developer can start working on a specific problem by name or identifier
- FR7: Developer can view problem metadata (difficulty, topic, tags, description)
- FR41: System stores all progress data locally using embedded database

**Epic 3: Test-Driven Practice Workflow**
- FR5: System provides problem description, constraints, and examples
- FR6: System provides pre-written test cases for every problem
- FR8: System generates scaffolded solution file with boilerplate code
- FR9: Developer can run tests against their solution using native Go testing
- FR10: System opens problem files in developer's preferred code editor
- FR11: System validates solution correctness using pre-written test cases
- FR12: Developer receives clear test results showing passed/failed test cases
- FR13: System tracks test execution history for each problem attempt

**Epic 4: Progress Tracking & Status**
- FR14: Developer can view their overall progress
- FR15: System tracks completion status for each problem
- FR16: Developer can view their solution history
- FR17: System records timestamps for problem attempts and completions
- FR44: System maintains data integrity across sessions with zero progress loss

**Epic 5: CLI Configuration & Customization**
- FR22: Developer can configure default code editor via configuration file
- FR23: Developer can configure default code editor via environment variables
- FR24: Developer can override configuration settings via command-line flags
- FR25: Developer can filter problems by difficulty preference in configuration
- FR26: Developer can set output verbosity preferences
- FR27: Developer can enable or disable celebration features
- FR28: Developer can enable or disable programming humor
- FR29: System respects standard environment conventions (NO_COLOR, TTY detection)
- FR45: System supports both global and project-specific configuration

**Epic 6: Advanced Output & Reporting**
- FR30: System provides human-friendly terminal output with colors and formatting
- FR31: Developer can request machine-parseable output (JSON format)
- FR32: System displays progress visualizations (progress bars, ASCII charts)
- FR33: System provides clear error messages when commands fail
- FR34: System outputs structured data to stdout and errors to stderr
- FR46: Commands can be composed and chained using shell operators
- FR47: System provides consistent exit codes for success and different error types
- FR48: Commands work in non-interactive mode (CI/CD, scripts)
- FR49: System detects TTY vs non-TTY contexts and adjusts output
- FR50: Developer can force decisions via flags to skip confirmations

**Epic 7: Motivation & Celebration System [Phase 2]**
- FR18: System tracks daily practice streaks
- FR36: System celebrates milestone achievements with ASCII art
- FR37: System displays streak counter with visual indicators
- FR38: System provides encouraging messages with programming humor
- FR39: System recognizes specific milestones
- FR40: System displays topic completion progress and celebrates mastery

**Epic 8: Shell Completion & Enhanced UX [Phase 2]**
- FR35: Developer can use shell completion for commands and problem names

**Epic 9: Advanced Analytics & Recommendations [Phase 3]**
- FR19: System identifies weak areas based on problem-solving history
- FR20: Developer can view personalized problem recommendations
- FR21: System tracks time spent on each problem for analytics

**Epic 10: Community & Extensibility [Phase 4]**
- FR42: Developer can export their progress data in portable formats
- FR43: Developer can import previously exported progress data
- FR51: Community contributor can submit new problem definitions
- FR52: Community contributor can add support for additional programming languages
- FR53: System provides clear documentation for contribution workflows
- FR54: Developer can create custom problem sets and themed collections

## Epic List

### Epic 1: Project Foundation & Setup
Developers can set up the dsa CLI tool, initialize their practice workspace, and have a working project structure with all infrastructure in place.

**FRs covered:** FR1
**Additional Requirements:** Cobra CLI starter template (5-step initialization), GORM v1.30.1+ with SQLite setup, GitHub Actions CI/CD pipeline, GoReleaser configuration, Standard Go project layout (cmd/, internal/, pkg/), Implementation patterns

**Technical Notes:**
- Initialize using `cobra-cli init` and `cobra-cli add` commands
- Set up GORM AutoMigrate for Problem, Solution, Progress models
- Configure GitHub Actions for testing (matrix: ubuntu/macos/windows)
- Configure GoReleaser for v1.0.0 release
- Establish project structure following architecture patterns

---

### Epic 2: Problem Library & Discovery
Developers can browse and discover DSA problems to practice, filtered by topic and difficulty, with complete problem metadata.

**FRs covered:** FR2, FR3, FR4, FR7, FR41
**NFRs addressed:** NFR6 (browsing <100ms), NFR20-24 (cross-platform support)

**Technical Notes:**
- Implement problem library with 10-15 initial problems (Arrays, Linked Lists, Trees, Sorting)
- Problems stored as embedded definitions in `problems/` directory
- Database queries optimized for <100ms browsing performance
- Problem metadata includes difficulty (easy/medium/hard), topic, tags, description
- Local SQLite storage with GORM models

---

### Epic 3: Test-Driven Practice Workflow
Developers can solve problems using test-driven development with native Go testing, including automatic scaffolding, editor integration, and test execution.

**FRs covered:** FR5, FR6, FR8, FR9, FR10, FR11, FR12, FR13
**NFRs addressed:** NFR1-4 (performance <500ms cold, <100ms warm), NFR7 (zero test overhead), NFR14-15 (Go integration)

**Technical Notes:**
- Scaffold generator creates boilerplate Go files with problem description
- Native `go test` integration with zero wrapper overhead
- Editor detection via $EDITOR environment variable
- Test results parsed and formatted for clear pass/fail display
- Test execution history tracked in database

---

### Epic 4: Progress Tracking & Status
Developers can track their progress and view completion statistics showing overall progress, problem-by-problem status, and solution history.

**FRs covered:** FR14, FR15, FR16, FR17, FR44
**NFRs addressed:** NFR5 (status <300ms), NFR8-13 (100% data integrity, offline operation, ACID transactions)

**Technical Notes:**
- Progress dashboard with problems solved by difficulty/topic
- Completion status tracking (not_started, in_progress, completed)
- Solution history with timestamps
- GORM transactions ensure atomic updates
- Database integrity with zero data loss guarantee

---

### Epic 5: CLI Configuration & Customization
Developers can customize dsa to their preferences and workflow with configuration via files, environment variables, and command-line flags.

**FRs covered:** FR22, FR23, FR24, FR25, FR26, FR27, FR28, FR29, FR45
**NFRs addressed:** NFR16-18 ($EDITOR, NO_COLOR, UNIX conventions), NFR31-36 (usability, defaults work for 80%)

**Technical Notes:**
- Viper configuration with precedence: flags > env > file > defaults
- Config files: ~/.dsa/config.yaml (global), .dsa/config.yaml (project)
- Environment variable prefix: DSA_*
- TTY detection for appropriate output formatting
- Configuration validation with clear error messages

---

### Epic 6: Advanced Output & Reporting
Developers get rich terminal output with colors and visualizations, plus machine-parseable JSON output for scripting and automation.

**FRs covered:** FR30, FR31, FR32, FR33, FR34, FR46, FR47, FR48, FR49, FR50
**NFRs addressed:** NFR17-18 (NO_COLOR, exit codes, UNIX conventions)

**Technical Notes:**
- Human output: ANSI colors, ASCII art, tables, progress bars
- Machine output: --json flag for all commands
- Exit codes: 0 (success), 1 (general), 2 (usage), 3 (database), 4 (test failure)
- Stdout for results, Stderr for logs/errors
- Composable commands for pipes and scripts

---

### Epic 7: Motivation & Celebration System [Phase 2]
Developers feel encouraged and celebrated as they practice, with streak tracking, milestone recognition, and programming humor.

**FRs covered:** FR18, FR36, FR37, FR38, FR39, FR40

**Technical Notes:**
- Streak calculation and visual counter display
- ASCII art celebrations for milestones
- Programming humor in success messages
- Milestone recognition (first problem, first hard, streak milestones)
- Topic mastery celebrations

---

### Epic 8: Shell Completion & Enhanced UX [Phase 2]
Developers get professional shell completion for bash, zsh, and fish shells with dynamic problem name completion.

**FRs covered:** FR35
**NFRs addressed:** NFR19 (bash/zsh/fish completion integration)

**Technical Notes:**
- Cobra built-in completion generation
- Dynamic completion using ValidArgsFunction (queries database for problem names)
- Installation: `dsa completion bash/zsh/fish`

---

### Epic 9: Advanced Analytics & Recommendations [Phase 3]
Developers get personalized insights including weak area identification, problem recommendations, and time tracking analytics.

**FRs covered:** FR19, FR20, FR21

**Technical Notes:**
- Analysis of problem-solving patterns to identify weak topics
- Recommendation engine based on practice history
- Time tracking per problem for performance analytics
- Smart learning features using progress data

---

### Epic 10: Community & Extensibility [Phase 4]
Community can contribute problems, add language support, and developers can share progress data and create custom problem sets.

**FRs covered:** FR42, FR43, FR51, FR52, FR53, FR54
**NFRs addressed:** NFR25-29 (language extensibility, problem format contributions, migrations)

**Technical Notes:**
- Import/export in JSON/CSV formats
- Standardized problem definition format for contributions
- Plugin architecture for additional languages (Rust, Python, TypeScript)
- Custom problem set creation
- Contribution documentation and workflows

---

# Epic Stories

## Epic 1: Project Foundation & Setup

Developers can set up the dsa CLI tool, initialize their practice workspace, and have a working project structure with all infrastructure in place.

### Story 1.1: Initialize Cobra CLI Project

As a **developer**,
I want **to set up the basic Cobra CLI project structure**,
So that **I have a working foundation for building the dsa CLI commands**.

**Acceptance Criteria:**

**Given** I have Go 1.23+ installed on my system
**When** I run the Cobra CLI initialization commands
**Then** The project is created with standard Go layout (cmd/, internal/, pkg/, main.go)
**And** The following commands are scaffolded: root, init, solve, test, status
**And** Running `go run main.go` executes successfully and displays help text
**And** The project follows the architecture's Standard Go Layout pattern
**And** All commands use PascalCase for exported functions and camelCase for unexported

**Given** The Cobra project is initialized
**When** I run `go mod tidy`
**Then** All dependencies are resolved (cobra, viper)
**And** The go.mod file specifies Go 1.23 as minimum version

**Given** The project structure is created
**When** I verify the directory layout
**Then** I see cmd/root.go with global flags defined
**And** I see cmd/init.go, cmd/solve.go, cmd/test.go, cmd/status.go command files
**And** I see internal/ directory ready for application packages
**And** I see testdata/ directory for test fixtures and golden files

### Story 1.2: Configure Database Layer with GORM

As a **developer**,
I want **to set up GORM with SQLite and define the core database models**,
So that **the application can persist problem, solution, and progress data locally**.

**Acceptance Criteria:**

**Given** The Cobra project is initialized
**When** I add GORM dependencies (gorm.io/gorm v1.30.1+, gorm.io/driver/sqlite v1.6.1+)
**Then** The dependencies are added to go.mod
**And** Running `go mod tidy` resolves all GORM packages

**Given** GORM is installed
**When** I create internal/database/models.go
**Then** The file defines Problem, Solution, Progress structs with proper GORM tags
**And** Problem model has: ID (primaryKey), Slug (uniqueIndex), Title, Difficulty, Topic, Description, CreatedAt
**And** Solution model has: ID (primaryKey), ProblemID (index), Code, Language, Passed, CreatedAt
**And** Progress model has: ID (primaryKey), ProblemID (uniqueIndex), Status, Attempts, LastAttempt
**And** All struct tags follow architecture pattern: `gorm:"primaryKey" json:"id"`
**And** All JSON fields use snake_case naming convention

**Given** Models are defined
**When** I create internal/database/connection.go with Initialize() function
**Then** The function opens SQLite database at ~/.dsa/dsa.db
**And** The function runs db.AutoMigrate(&Problem{}, &Solution{}, &Progress{})
**And** The function returns *gorm.DB connection or error
**And** Database initialization handles errors gracefully with wrapped context

**Given** Database initialization exists
**When** I run the application
**Then** The SQLite database file is created at ~/.dsa/dsa.db
**And** All three tables (problems, solutions, progress) are created
**And** Table names are plural snake_case (problems, solutions, progress)
**And** Column names are snake_case (problem_id, created_at, difficulty_level)

### Story 1.3: Implement Workspace Initialization Command

As a **developer**,
I want **to run `dsa init` to set up my practice workspace**,
So that **I can start practicing DSA problems in my local environment**.

**Acceptance Criteria:**

**Given** The database layer is configured
**When** I implement cmd/init.go command
**Then** The command initializes the database connection
**And** The command creates ~/.dsa directory if it doesn't exist
**And** The command creates ~/.dsa/dsa.db SQLite database
**And** The command runs GORM AutoMigrate for all models
**And** The command outputs success message: "Workspace initialized at ~/.dsa"

**Given** I have not initialized dsa yet
**When** I run `dsa init`
**Then** The command completes in <500ms (NFR1: cold start performance)
**And** The ~/.dsa directory is created with proper permissions
**And** The database file is created with tables: problems, solutions, progress
**And** I see output: "✓ Workspace initialized at /Users/[username]/.dsa"
**And** The command exits with code 0 (success)

**Given** I have already initialized dsa
**When** I run `dsa init` again
**Then** The command detects existing workspace
**And** I see message: "Workspace already initialized at ~/.dsa"
**And** The command does not recreate or corrupt existing data (NFR8: data integrity)
**And** The command exits with code 0

**Given** The database initialization fails (e.g., permission denied)
**When** I run `dsa init`
**Then** The command outputs clear error message (NFR31: actionable errors)
**And** The error message suggests resolution (e.g., "Check directory permissions")
**And** The command exits with code 3 (database error per architecture)
**And** No partial/corrupted database is left behind

### Story 1.4: Set Up CI/CD Pipeline

As a **developer**,
I want **automated testing and release workflows configured**,
So that **code quality is enforced and releases are automated**.

**Acceptance Criteria:**

**Given** The project has Go code and tests
**When** I create .github/workflows/ci.yml
**Then** The workflow is triggered on push and pull_request events
**And** The workflow runs tests on matrix: [ubuntu-latest, macos-latest, windows-latest]
**And** The workflow uses Go version 1.23
**And** The workflow runs: `go test -race -coverprofile=coverage.txt -covermode=atomic ./...`
**And** The workflow uploads coverage to codecov
**And** The workflow runs golangci-lint for code quality

**Given** The CI workflow is configured
**When** I create .github/workflows/release.yml
**Then** The workflow triggers on git tags matching `v*` pattern
**And** The workflow depends on CI tests passing first
**And** The workflow uses goreleaser/goreleaser-action@v6
**And** The workflow has permissions to write releases

**Given** Release workflow exists
**When** I create .goreleaser.yml configuration
**Then** The config builds for: linux/darwin/windows with amd64/arm64 architectures
**And** The config sets CGO_ENABLED=0 for static binaries
**And** The config injects version/commit/date via ldflags
**And** The config creates archives (tar.gz for unix, zip for windows)
**And** The config generates checksums
**And** Binary size remains <20MB (NFR23: distribution size)

**Given** CI/CD is fully configured
**When** I push code to GitHub
**Then** CI tests run automatically on all three platforms
**And** Linting checks enforce Go best practices (NFR30)
**And** Test coverage is tracked and reported

**Given** CI/CD is configured
**When** I create and push a git tag (e.g., v1.0.0)
**Then** The release workflow builds binaries for all platforms
**And** GitHub release is created with compiled binaries attached
**And** Release includes: dsa_Linux_x86_64, dsa_Darwin_x86_64, dsa_Windows_x86_64.exe
**And** Checksums file is included for verification

### Story 1.5: Establish Project Documentation

As a **developer**,
I want **comprehensive project documentation**,
So that **I understand how to use, develop, and contribute to dsa**.

**Acceptance Criteria:**

**Given** The project is initialized
**When** I create README.md
**Then** The README includes: project description, core value proposition ("celebrating little victories")
**And** The README includes installation instructions (binary download + `go install`)
**And** The README includes quick start: `dsa init`, `dsa solve [problem]`, `dsa test`
**And** The README includes configuration section (editor, difficulty, verbosity)
**And** The README includes example terminal output showing colored CLI
**And** The README includes link to full documentation
**And** Setup instructions take <5 minutes from clone to first problem (NFR33)

**Given** The project has contribution potential
**When** I create CONTRIBUTING.md
**Then** The document explains how to submit issues and PRs
**And** The document specifies code style requirements (gofmt, golangci-lint)
**And** The document explains the PR review process
**And** The document includes testing requirements (70%+ coverage, table-driven tests)
**And** The document references implementation patterns from architecture.md

**Given** The project needs a license
**When** I create LICENSE file
**Then** The file contains MIT license text
**And** Copyright year is 2025
**And** License allows open source community contributions (supports FR51-54)

**Given** The project needs a changelog
**When** I create CHANGELOG.md
**Then** The file follows Keep a Changelog format
**And** The file has [Unreleased] section for upcoming changes
**And** The file will track versions per semantic versioning (v1.0.0, v1.1.0, etc.)

**Given** All documentation is created
**When** A new developer clones the repository
**Then** They can understand the project purpose from README
**And** They can set up their environment in <5 minutes following docs (NFR33)
**And** They can find contribution guidelines easily
**And** Documentation explains the "celebrating little victories" differentiator

---

## Epic 2: Problem Library & Discovery

### Story 2.1: Seed Initial Problem Library

As a **user**,
I want **an initial library of curated DSA problems available after workspace initialization**,
So that **I can immediately start practicing without manual setup**.

**Acceptance Criteria:**

**Given** I have initialized a workspace with `dsa init`
**When** I run `dsa list`
**Then** I see at least 20 pre-seeded problems in the library
**And** Problems cover core topics: Arrays, Linked Lists, Trees, Graphs, Sorting, Searching (FR2)
**And** Each problem has difficulty level: Easy, Medium, or Hard (FR3)
**And** Each problem includes test cases and boilerplate code (FR10)
**And** All problems are stored in the local SQLite database (FR41)
**And** Problem metadata follows the schema: id, title, description, difficulty, topic, tags, boilerplate_path, test_path

**Given** I inspect the seeded problems
**When** I check the problem files
**Then** Each problem has a corresponding Go file with boilerplate code
**And** Each problem has a corresponding test file with test cases
**And** File names follow snake_case convention (Architecture pattern)
**And** Test files use testify/assert for assertions (Architecture pattern)

---

### Story 2.2: Implement Problem Listing Command

As a **user**,
I want **to list problems with filtering options**,
So that **I can discover problems based on difficulty, topic, or status** (FR2, FR3).

**Acceptance Criteria:**

**Given** I have problems in my workspace
**When** I run `dsa list`
**Then** I see a formatted table with columns: ID, Title, Difficulty, Topic, Status (Unsolved/Solved)
**And** Command executes in <100ms (NFR1: warm start)
**And** Output uses color coding: Green for Solved, Yellow for In Progress, White for Unsolved (Architecture: color-coded output)

**Given** I want to filter by difficulty
**When** I run `dsa list --difficulty easy`
**Then** I see only Easy problems
**When** I run `dsa list --difficulty medium`
**Then** I see only Medium problems
**When** I run `dsa list --difficulty hard`
**Then** I see only Hard problems

**Given** I want to filter by topic
**When** I run `dsa list --topic arrays`
**Then** I see only Array problems
**When** I run `dsa list --topic "linked-lists"`
**Then** I see only Linked List problems

**Given** I want to filter by solved status
**When** I run `dsa list --unsolved`
**Then** I see only unsolved problems
**When** I run `dsa list --solved`
**Then** I see only solved problems

**Given** I want to combine filters
**When** I run `dsa list --difficulty medium --topic trees --unsolved`
**Then** I see only unsolved Medium-level Tree problems

---

### Story 2.3: Implement Problem Details Display

As a **user**,
I want **to view detailed information about a specific problem**,
So that **I can understand the problem requirements before solving** (FR2).

**Acceptance Criteria:**

**Given** I know a problem ID or title
**When** I run `dsa show <problem-id>`
**Then** I see the full problem details:
  - Title
  - Difficulty level
  - Topic/Tags
  - Problem description
  - Example inputs and outputs
  - Constraints
  - File paths (boilerplate and test)
  - Solution status (Unsolved/Solved/In Progress)
**And** Output is formatted with clear sections and readability (FR30)

**Given** The problem has been solved before
**When** I run `dsa show <problem-id>`
**Then** I see additional information:
  - Last solved date
  - Number of attempts
  - Best time/performance (if tracked)
  - Link to my solution file

**Given** I provide an invalid problem ID
**When** I run `dsa show invalid-id`
**Then** I see a helpful error message: "Problem 'invalid-id' not found. Use 'dsa list' to see available problems."
**And** The command exits with non-zero status code (UNIX convention, NFR28)

---

### Story 2.4: Implement Random Problem Selection

As a **user**,
I want **to get a random problem suggestion**,
So that **I can practice without decision fatigue** (FR4).

**Acceptance Criteria:**

**Given** I have unsolved problems in my library
**When** I run `dsa random`
**Then** I see a randomly selected unsolved problem displayed
**And** The display includes: Title, Difficulty, Topic, Description
**And** I see a suggested command: "Run 'dsa solve <problem-id>' to start"

**Given** I want a random problem of specific difficulty
**When** I run `dsa random --difficulty easy`
**Then** I get a random Easy problem
**When** I run `dsa random --difficulty medium`
**Then** I get a random Medium problem
**When** I run `dsa random --difficulty hard`
**Then** I get a random Hard problem

**Given** I want a random problem from a specific topic
**When** I run `dsa random --topic arrays`
**Then** I get a random Array problem
**When** I run `dsa random --topic graphs`
**Then** I get a random Graph problem

**Given** I want to combine filters
**When** I run `dsa random --difficulty hard --topic trees`
**Then** I get a random Hard-level Tree problem that is unsolved

**Given** All problems matching the filter are solved
**When** I run `dsa random --difficulty easy` and all Easy problems are solved
**Then** I see a message: "All Easy problems are solved! Try --difficulty medium or --difficulty hard"
**And** The command suggests next steps

---

### Story 2.5: Add Custom Problem to Library

As a **user**,
I want **to add my own custom problems to the library**,
So that **I can practice company-specific or niche problems** (FR5).

**Acceptance Criteria:**

**Given** I want to add a custom problem
**When** I run `dsa add "Two Sum" --difficulty easy --topic arrays`
**Then** The CLI prompts me for problem description
**And** The CLI creates:
  - A new problem entry in the database
  - A boilerplate Go file at `problems/two_sum.go` (snake_case)
  - A test file at `problems/two_sum_test.go`
**And** Files follow the standard template structure (Architecture: boilerplate generation)
**And** Test file uses testify/assert (Architecture pattern)

**Given** The custom problem is created
**When** I run `dsa list`
**Then** I see my custom problem "Two Sum" in the list
**And** It shows difficulty: Easy, topic: arrays

**Given** I want to add test cases to my custom problem
**When** I edit the generated test file
**Then** I can add table-driven tests following Go conventions (Architecture pattern)
**And** Running `dsa test two-sum` executes my custom test cases

**Given** I want to add optional tags
**When** I run `dsa add "Custom Problem" --difficulty medium --topic graphs --tags "bfs,shortest-path"`
**Then** The problem is created with tags stored in the database
**And** I can later filter by tags using `dsa list --tag bfs`

---

## Epic 3: Test-Driven Practice Workflow

### Story 3.1: Implement Solve Command with Boilerplate Generation

As a **user**,
I want **to start solving a problem with auto-generated boilerplate**,
So that **I can focus on the solution logic immediately** (FR8, FR9).

**Acceptance Criteria:**

**Given** I have selected a problem to solve
**When** I run `dsa solve <problem-id>`
**Then** The CLI generates a solution file at `solutions/<problem-id>.go`
**And** The file contains boilerplate code with function signature from the problem template
**And** The file includes package declaration and necessary imports
**And** The file follows snake_case naming for files (Architecture pattern)
**And** Exported functions use PascalCase (Architecture pattern)

**Given** The solution file is generated
**When** I check the file contents
**Then** I see helpful comments indicating where to write my solution
**And** I see the function signature matching the problem requirements
**And** I see example test case as a comment for reference

**Given** I have already started solving this problem before
**When** I run `dsa solve <problem-id>` again
**Then** The CLI asks: "Solution file already exists. Overwrite? [y/N]"
**And** If I choose 'N', the existing file is preserved
**And** If I choose 'y', a backup is created at `solutions/<problem-id>.go.backup` before overwriting

**Given** The solution is generated
**When** I run `dsa solve <problem-id> --open`
**Then** The solution file opens in my configured editor (from config or $EDITOR)
**And** I can start coding immediately

---

### Story 3.2: Implement Test Execution Command

As a **user**,
I want **to run tests against my solution**,
So that **I can validate my solution works correctly** (FR10, FR11).

**Acceptance Criteria:**

**Given** I have written a solution for a problem
**When** I run `dsa test <problem-id>`
**Then** The CLI executes `go test` on the problem's test file
**And** Test output shows pass/fail status with color coding (green for pass, red for fail)
**And** Failed tests show expected vs actual values clearly
**And** Test execution completes in <2 seconds for typical problems (NFR2)
**And** The CLI uses testify/assert for readable test failures (Architecture pattern)

**Given** All tests pass
**When** I run `dsa test <problem-id>`
**Then** I see: "✓ All tests passed! (5/5)"
**And** The problem status is updated to "Solved" in the database
**And** Progress tracking records the completion (links to Epic 4)
**And** The CLI displays an encouraging message (FR36 Phase 2 preview)

**Given** Some tests fail
**When** I run `dsa test <problem-id>`
**Then** I see: "✗ Tests failed (2/5 passed)"
**And** Failed test cases are displayed with details:
  - Test case name
  - Input values
  - Expected output
  - Actual output
**And** The problem status remains "In Progress"

**Given** I want verbose test output
**When** I run `dsa test <problem-id> --verbose`
**Then** I see detailed test execution logs including all test cases (passed and failed)
**And** Output includes timing information for each test

**Given** I want to run tests with Go's race detector
**When** I run `dsa test <problem-id> --race`
**Then** Tests execute with `go test -race` flag
**And** Any race conditions are reported clearly

---

### Story 3.3: Implement Watch Mode for Continuous Testing

As a **user**,
I want **to automatically re-run tests when I save my solution**,
So that **I get instant feedback without manually running tests** (FR12).

**Acceptance Criteria:**

**Given** I am working on a solution
**When** I run `dsa test <problem-id> --watch`
**Then** The CLI monitors the solution file for changes
**And** Tests automatically re-run whenever I save the file
**And** Test results display immediately in the terminal
**And** The watch process continues until I press Ctrl+C

**Given** Tests are running in watch mode
**When** I save my solution file
**Then** The CLI clears the terminal and shows: "🔄 Re-running tests..."
**And** Test results appear within 1 second (NFR2)
**And** Pass/fail status is clearly indicated with color and symbols

**Given** Watch mode is active
**When** Tests pass after previously failing
**Then** I see a celebratory message: "🎉 Tests now passing!"
**And** The terminal uses green color for success

**Given** Watch mode is active
**When** Tests fail after previously passing
**Then** I see: "⚠️  Tests broken - check your changes"
**And** The terminal uses red color for failure

---

### Story 3.4: Generate Test Cases for Custom Problems

As a **user**,
I want **to add test cases to my custom problems**,
So that **I can validate solutions against expected behavior** (FR13).

**Acceptance Criteria:**

**Given** I have created a custom problem
**When** I run `dsa test-gen <problem-id>`
**Then** The CLI prompts me to input test cases interactively:
  - Test case name/description
  - Input values (with type hints)
  - Expected output
**And** The CLI generates test functions in the test file
**And** Tests follow table-driven test pattern (Architecture: Go conventions)
**And** Tests use testify/assert for assertions (Architecture pattern)

**Given** I provide multiple test cases
**When** I run `dsa test-gen <problem-id>` and add 5 test cases
**Then** The generated test file contains a table-driven test with all 5 cases
**And** Each test case has: name, input, expected output
**And** The test function iterates over the table using `t.Run()` with subtests

**Given** I want to add test cases to an existing test file
**When** I run `dsa test-gen <problem-id> --append`
**Then** New test cases are added to the existing table-driven test
**And** Existing test cases are preserved

**Given** I want to provide test cases via JSON file
**When** I run `dsa test-gen <problem-id> --from-file testcases.json`
**Then** The CLI reads test cases from the JSON file
**And** Generates the corresponding Go test code
**And** JSON schema includes: name, inputs (array), expected (value)

---

### Story 3.5: Implement Solution Submission and History

As a **user**,
I want **to save my solution and maintain a history of attempts**,
So that **I can track my improvement over time** (FR14).

**Acceptance Criteria:**

**Given** I have a passing solution
**When** I run `dsa submit <problem-id>`
**Then** The CLI copies my solution to `solutions/history/<problem-id>/<timestamp>.go`
**And** The database records:
  - Problem ID
  - Submission timestamp
  - Pass/fail status
  - Number of test cases passed
  - Solution file path
**And** I see confirmation: "✓ Solution submitted and saved to history"

**Given** I have multiple submissions for a problem
**When** I run `dsa history <problem-id>`
**Then** I see a list of all my attempts with:
  - Submission date/time
  - Pass/fail status
  - Test results (e.g., "5/5 tests passed")
**And** The list is sorted by most recent first

**Given** I want to view a previous solution
**When** I run `dsa history <problem-id> --show 2`
**Then** The CLI displays the solution code from my 2nd most recent attempt
**And** I can compare it with my current solution

**Given** I want to restore a previous solution
**When** I run `dsa history <problem-id> --restore 2`
**Then** The CLI asks: "Restore solution from <timestamp>? Current solution will be backed up. [y/N]"
**And** If I confirm, the selected solution is copied to `solutions/<problem-id>.go`
**And** The current solution is backed up before restoration

---

### Story 3.6: Implement Benchmark Support for Performance Testing

As a **user**,
I want **to run benchmarks on my solution**,
So that **I can measure and optimize performance** (FR11 extension, NFR1).

**Acceptance Criteria:**

**Given** I have a working solution
**When** I run `dsa bench <problem-id>`
**Then** The CLI executes `go test -bench` on the problem's test file
**And** Benchmark results show:
  - Iterations per second
  - Nanoseconds per operation
  - Allocations per operation
  - Bytes allocated per operation
**And** Results are formatted clearly and color-coded (Architecture: color-coded output)

**Given** I want to compare benchmark results
**When** I run `dsa bench <problem-id> --save`
**Then** Benchmark results are saved with timestamp
**And** Subsequent runs show comparison with previous best:
  - "🚀 15% faster than previous best!"
  - "⚠️  20% slower than previous best"
  - "💾 30% less memory than previous best!"

**Given** I want detailed memory profiling
**When** I run `dsa bench <problem-id> --mem`
**Then** The CLI runs benchmarks with memory profiling enabled
**And** Output shows detailed allocation breakdown
**And** I can identify optimization opportunities

**Given** I want to run benchmarks with CPU profiling
**When** I run `dsa bench <problem-id> --cpuprofile`
**Then** The CLI generates a CPU profile file
**And** I see a message: "Profile saved to <problem-id>.cpu.prof. View with 'go tool pprof'"
**And** I can analyze performance bottlenecks using Go's profiling tools

---

## Epic 4: Progress Tracking & Status

### Story 4.1: Implement Status Dashboard Command

As a **user**,
I want **to see an overview of my progress at a glance**,
So that **I can understand my current standing and what to work on next** (FR14, FR15).

**Acceptance Criteria:**

**Given** I have been solving problems
**When** I run `dsa status`
**Then** I see a formatted dashboard with:
  - Total problems solved (count and percentage)
  - Breakdown by difficulty: Easy (X/Y), Medium (X/Y), Hard (X/Y)
  - Breakdown by topic: Arrays (X/Y), Trees (X/Y), etc.
  - Recent activity: Last 5 problems solved with dates
  - Current streak (if implemented in Phase 2)
**And** The dashboard executes in <100ms (NFR1: warm start)
**And** Output uses color coding for visual clarity (Architecture pattern)

**Given** I have solved problems across multiple difficulties
**When** I run `dsa status`
**Then** I see a progress bar for each difficulty level
**And** Percentages are displayed (e.g., "Easy: 15/30 [50%] ████████░░░░░░░░")
**And** Colors indicate completion level: red (<30%), yellow (30-70%), green (>70%)

**Given** I want to see status for a specific topic
**When** I run `dsa status --topic arrays`
**Then** I see detailed stats only for Array problems
**And** The output shows: total solved, by difficulty, recent solutions

**Given** I want compact status output
**When** I run `dsa status --compact`
**Then** I see a one-line summary: "Progress: 45/100 solved (Easy: 20/30, Med: 15/40, Hard: 10/30) | Streak: 7 days"

---

### Story 4.2: Implement Progress Tracking Database Models

As a **developer**,
I want **database models to track user progress and activity**,
So that **the system can record and query practice history** (FR14, NFR41).

**Acceptance Criteria:**

**Given** I need to track solution attempts
**When** I create the Solution model using GORM
**Then** The model includes fields:
  - ID (primary key)
  - ProblemID (foreign key to Problem)
  - FilePath (string, path to solution file)
  - SubmittedAt (timestamp)
  - Status (enum: Passed, Failed, InProgress)
  - TestsPassed (int, number of tests passed)
  - TestsTotal (int, total number of tests)
**And** The table uses snake_case naming (Architecture: DB conventions)
**And** GORM AutoMigrate creates the table automatically (Architecture pattern)

**Given** I need to track overall progress
**When** I create the Progress model using GORM
**Then** The model includes fields:
  - ID (primary key)
  - ProblemID (foreign key to Problem)
  - FirstSolvedAt (timestamp, nullable)
  - LastAttemptedAt (timestamp)
  - TotalAttempts (int)
  - BestTime (int, milliseconds, nullable)
  - IsSolved (boolean)
**And** The model has indexes on ProblemID and IsSolved for fast queries
**And** The table uses snake_case naming (Architecture: DB conventions)

**Given** The models are defined
**When** The CLI starts
**Then** GORM AutoMigrate runs and creates/updates tables as needed (Architecture pattern)
**And** Database constraints enforce referential integrity
**And** Migrations are non-destructive (existing data preserved)

---

### Story 4.3: Track Problem Completion and Update Status

As a **user**,
I want **my progress to be automatically tracked when I solve problems**,
So that **my stats are always up-to-date without manual input** (FR14).

**Acceptance Criteria:**

**Given** I run tests on a problem for the first time
**When** I run `dsa test <problem-id>`
**Then** A Progress record is created in the database with:
  - ProblemID: <problem-id>
  - FirstSolvedAt: null (not solved yet)
  - LastAttemptedAt: current timestamp
  - TotalAttempts: 1
  - IsSolved: false

**Given** All tests pass for a previously unsolved problem
**When** I run `dsa test <problem-id>` and all tests pass
**Then** The Progress record is updated:
  - FirstSolvedAt: current timestamp
  - LastAttemptedAt: current timestamp
  - TotalAttempts: incremented
  - IsSolved: true
**And** A Solution record is created with status: Passed

**Given** Tests fail
**When** I run `dsa test <problem-id>` and tests fail
**Then** The Progress record is updated:
  - LastAttemptedAt: current timestamp
  - TotalAttempts: incremented
  - IsSolved: remains false
**And** A Solution record is created with status: Failed, TestsPassed: X, TestsTotal: Y

**Given** I re-solve a previously solved problem
**When** I run `dsa test <problem-id>` and all tests pass
**Then** FirstSolvedAt is NOT changed (preserves original solve date)
**And** LastAttemptedAt is updated to current timestamp
**And** A new Solution record is created

---

### Story 4.4: Implement Problem Stats and Analytics

As a **user**,
I want **to view detailed statistics for individual problems**,
So that **I can see my performance history on specific problems** (FR15).

**Acceptance Criteria:**

**Given** I have attempted a problem multiple times
**When** I run `dsa stats <problem-id>`
**Then** I see detailed stats:
  - Problem title, difficulty, topic
  - First solved: <date> or "Not solved yet"
  - Last attempted: <date>
  - Total attempts: <count>
  - Success rate: X% (passed attempts / total attempts)
  - Best time: <duration> (if tracked)
  - All attempts history with timestamps and pass/fail status

**Given** I have never attempted a problem
**When** I run `dsa stats <problem-id>`
**Then** I see: "Problem '<title>' - Not attempted yet"
**And** The output shows problem details (difficulty, topic, description)
**And** Suggestion: "Run 'dsa solve <problem-id>' to start"

**Given** I want to compare my performance across attempts
**When** I run `dsa stats <problem-id> --attempts`
**Then** I see a table showing each attempt:
  - Attempt # | Date | Status | Tests Passed | Notes
  - 1 | 2025-01-15 | Failed | 3/5 | -
  - 2 | 2025-01-16 | Passed | 5/5 | ✓ First solve
  - 3 | 2025-01-20 | Passed | 5/5 | Re-solve

---

### Story 4.5: Export Progress Data

As a **user**,
I want **to export my progress data in various formats**,
So that **I can back up, share, or analyze my practice history** (FR41, FR42, FR43).

**Acceptance Criteria:**

**Given** I want to export my progress
**When** I run `dsa export --format json`
**Then** A JSON file is created at `~/.dsa/exports/progress_<timestamp>.json`
**And** The file contains:
  - All problems with their details
  - All progress records
  - All solution attempts
  - Export metadata (timestamp, version, user)
**And** JSON is properly formatted and valid

**Given** I want to export to CSV
**When** I run `dsa export --format csv`
**Then** A CSV file is created at `~/.dsa/exports/progress_<timestamp>.csv`
**And** The file has columns: ProblemID, Title, Difficulty, Topic, Solved, FirstSolvedAt, TotalAttempts
**And** CSV is properly escaped and can be opened in Excel/Google Sheets

**Given** I want to specify output location
**When** I run `dsa export --format json --output /path/to/backup.json`
**Then** The export is saved to the specified path
**And** I see confirmation: "Progress exported to /path/to/backup.json"

**Given** I want to export only solved problems
**When** I run `dsa export --format json --solved-only`
**Then** The export contains only problems marked as solved
**And** Unsolved problems are excluded from the export

**Given** I want a summary report
**When** I run `dsa export --format summary`
**Then** A human-readable text file is created with:
  - Overall statistics (total solved, by difficulty, by topic)
  - Problem list grouped by difficulty and topic
  - Recent activity summary
  - Formatted for easy reading

---

## Epic 5: CLI Configuration & Customization

### Story 5.1: Implement Configuration File with Viper

As a **user**,
I want **a configuration file to customize CLI behavior**,
So that **I can personalize the tool to my preferences** (FR22, FR23).

**Acceptance Criteria:**

**Given** I initialize a workspace
**When** The workspace is created
**Then** A config file is created at `~/.dsa/config.yaml`
**And** The file uses YAML format for readability
**And** Viper library loads the configuration (Architecture: use Viper)
**And** The config includes default values for all configurable options

**Given** I want to view current configuration
**When** I run `dsa config list`
**Then** I see all configuration options and their current values:
  - editor: "vim"
  - output_format: "table"
  - color_enabled: true
  - default_difficulty: null
  - test_timeout: 30
**And** Output shows which values are defaults vs user-configured

**Given** The config file is missing or corrupted
**When** I run any `dsa` command
**Then** The CLI uses sensible defaults and continues execution
**And** A warning is logged: "Config file not found, using defaults"
**And** The CLI offers to regenerate the config: "Run 'dsa config init' to create a new config file"

---

### Story 5.2: Implement Config Get/Set Commands

As a **user**,
I want **to view and modify configuration values via CLI**,
So that **I can customize settings without manually editing files** (FR22, FR23).

**Acceptance Criteria:**

**Given** I want to view a specific config value
**When** I run `dsa config get editor`
**Then** I see the current value: "editor: vim"
**And** If the value is not set, I see: "editor: <not set> (default: vim)"

**Given** I want to change a config value
**When** I run `dsa config set editor code`
**Then** The config file is updated with the new value
**And** I see confirmation: "✓ Configuration updated: editor = code"
**And** The change takes effect immediately for subsequent commands

**Given** I want to set the default editor
**When** I run `dsa config set editor vim`
**Then** All `dsa solve --open` commands use vim to open files
**When** I run `dsa config set editor code`
**Then** All `dsa solve --open` commands use VS Code to open files

**Given** I want to enable/disable color output
**When** I run `dsa config set color_enabled false`
**Then** All CLI output is monochrome (no ANSI color codes)
**When** I run `dsa config set color_enabled true`
**Then** CLI output uses color coding for better readability

**Given** I provide an invalid config key
**When** I run `dsa config set invalid_key value`
**Then** I see an error: "Invalid configuration key: 'invalid_key'"
**And** I see a list of valid config keys
**And** The command exits with non-zero status code

---

### Story 5.3: Configure Default Output Formats

As a **user**,
I want **to configure default output formats for different commands**,
So that **I don't need to specify format flags repeatedly** (FR30, FR31).

**Acceptance Criteria:**

**Given** I prefer table output for problem lists
**When** I run `dsa config set list_format table`
**Then** All `dsa list` commands use table format by default
**And** I can still override with `dsa list --format json`

**Given** I prefer JSON output for status
**When** I run `dsa config set status_format json`
**Then** All `dsa status` commands output JSON by default
**And** JSON output is properly formatted and valid

**Given** I want compact output by default
**When** I run `dsa config set output_style compact`
**Then** All status and progress commands use compact formatting
**And** I can override with `dsa status --verbose` for detailed output

**Given** I want to configure color preferences
**When** I run `dsa config set color_scheme solarized`
**Then** CLI uses solarized color palette for output
**And** Supported schemes include: default, solarized, monokai, nord

---

### Story 5.4: Configure Editor Integration

As a **user**,
I want **to configure my preferred editor and editor-specific options**,
So that **the CLI integrates seamlessly with my development environment** (FR24, FR25).

**Acceptance Criteria:**

**Given** I want to set my editor
**When** I run `dsa config set editor code`
**Then** The config stores: `editor: code`
**And** `dsa solve --open` uses `code <file>` to open solutions

**Given** I want to set editor-specific arguments
**When** I run `dsa config set editor_args "--goto {file}:{line}"`
**Then** The CLI uses those args when opening files
**And** Placeholders are supported: {file}, {line}, {column}

**Given** I use VS Code
**When** I run `dsa config set editor code --args "--goto {file}:1"`
**Then** `dsa solve two-sum --open` executes: `code --goto solutions/two_sum.go:1`

**Given** I use Vim
**When** I run `dsa config set editor vim --args "+{line}"`
**Then** `dsa solve two-sum --open` executes: `vim +1 solutions/two_sum.go`

**Given** I use Neovim with LSP
**When** I run `dsa config set editor nvim --args "-c 'normal {line}G'"`
**Then** Files open at the correct line number with Neovim

**Given** No editor is configured
**When** I run `dsa solve --open`
**Then** The CLI checks $EDITOR environment variable
**And** If $EDITOR is not set, uses system default (open on macOS, xdg-open on Linux)
**And** A hint is shown: "Set your preferred editor with: dsa config set editor <editor>"

---

### Story 5.5: Implement Configuration Profiles

As a **user**,
I want **to save and switch between configuration profiles**,
So that **I can use different settings for different contexts** (FR26, FR27).

**Acceptance Criteria:**

**Given** I want to create a new profile
**When** I run `dsa config profile create work`
**Then** A new profile "work" is created as a copy of current config
**And** I see confirmation: "✓ Profile 'work' created"
**And** The profile is stored at `~/.dsa/profiles/work.yaml`

**Given** I have multiple profiles
**When** I run `dsa config profile list`
**Then** I see all available profiles:
  - default (active) *
  - work
  - personal
**And** The active profile is marked with an asterisk

**Given** I want to switch profiles
**When** I run `dsa config profile switch work`
**Then** The CLI loads configuration from the "work" profile
**And** I see confirmation: "✓ Switched to profile 'work'"
**And** All subsequent commands use the "work" profile settings

**Given** I want to delete a profile
**When** I run `dsa config profile delete work`
**Then** The CLI asks for confirmation: "Delete profile 'work'? [y/N]"
**And** If confirmed, the profile file is deleted
**And** If it was the active profile, CLI switches to "default"

**Given** I want to export a profile
**When** I run `dsa config profile export work --output work-config.yaml`
**Then** The profile is exported to the specified file
**And** The file can be shared with other users or machines

**Given** I want to import a profile
**When** I run `dsa config profile import --file work-config.yaml --name imported`
**Then** A new profile "imported" is created from the file
**And** I can switch to it with `dsa config profile switch imported`

---

### Story 5.6: Implement Configuration Validation and Reset

As a **user**,
I want **to validate my configuration and reset to defaults if needed**,
So that **I can fix configuration issues easily** (FR28, FR29).

**Acceptance Criteria:**

**Given** I have modified my configuration
**When** I run `dsa config validate`
**Then** The CLI checks all config values for validity:
  - Data types (string, int, bool)
  - Value ranges (e.g., test_timeout > 0)
  - Valid enum values (e.g., color_scheme in [default, solarized, ...])
**And** If valid, I see: "✓ Configuration is valid"
**And** If invalid, I see specific errors with suggestions for fixes

**Given** My configuration is corrupted
**When** I run `dsa config validate`
**Then** I see detailed errors:
  - "Error: 'test_timeout' must be a positive integer, got: 'invalid'"
  - "Error: 'color_scheme' must be one of: default, solarized, monokai, nord. Got: 'unknown'"
**And** I see suggestion: "Run 'dsa config reset' to restore defaults"

**Given** I want to reset all configuration
**When** I run `dsa config reset`
**Then** The CLI asks for confirmation: "Reset all configuration to defaults? [y/N]"
**And** If confirmed, a backup is created at `~/.dsa/config.yaml.backup.<timestamp>`
**And** The config file is replaced with factory defaults
**And** I see: "✓ Configuration reset to defaults. Backup saved to <path>"

**Given** I want to reset a specific key
**When** I run `dsa config reset editor`
**Then** Only the "editor" key is reset to its default value
**And** Other configuration values are preserved
**And** I see: "✓ 'editor' reset to default value: vim"

**Given** I want to view the default configuration
**When** I run `dsa config defaults`
**Then** I see all default values in YAML format
**And** I can use this as a reference for resetting specific values

---

## Epic 6: Advanced Output & Reporting

### Story 6.1: Implement JSON Output Mode

As a **user**,
I want **to output command results in JSON format**,
So that **I can parse and process data programmatically** (FR30, FR31).

**Acceptance Criteria:**

**Given** I want JSON output for problem lists
**When** I run `dsa list --format json`
**Then** Output is valid JSON with structure:
```json
{
  "problems": [
    {
      "id": "two-sum",
      "title": "Two Sum",
      "difficulty": "easy",
      "topic": "arrays",
      "solved": false
    }
  ],
  "total": 20,
  "solved": 5
}
```
**And** JSON is properly formatted with indentation
**And** Output can be piped to jq or other JSON tools

**Given** I want JSON output for status
**When** I run `dsa status --format json`
**Then** Output includes all status data in JSON format:
  - Total problems and solved count
  - Breakdown by difficulty (easy/medium/hard)
  - Breakdown by topic
  - Recent activity array
  - Streak information (if available)
**And** JSON schema is consistent and documented

**Given** I want compact JSON (no formatting)
**When** I run `dsa list --format json --compact`
**Then** Output is minified JSON on a single line
**And** Useful for piping to other tools or APIs

---

### Story 6.2: Implement Table Output with Formatting

As a **user**,
I want **well-formatted table output for CLI commands**,
So that **data is easy to read in the terminal** (FR30).

**Acceptance Criteria:**

**Given** I want to list problems in table format
**When** I run `dsa list --format table` (or just `dsa list`)
**Then** Output is a formatted ASCII table:
```
+------------+-------------------+------------+-----------+--------+
| ID         | Title             | Difficulty | Topic     | Status |
+------------+-------------------+------------+-----------+--------+
| two-sum    | Two Sum           | Easy       | Arrays    | ✓      |
| add-two-no | Add Two Numbers   | Medium     | Linked... | ✗      |
+------------+-------------------+------------+-----------+--------+
```
**And** Columns are aligned and padded properly
**And** Long values are truncated with ellipsis (...)
**And** Table fits within terminal width (respects $COLUMNS or detects automatically)

**Given** I want colored table output
**When** I run `dsa list` with color enabled
**Then** Table uses colors:
  - Headers: Bold
  - Solved problems (✓): Green
  - Unsolved problems (✗): Yellow
  - Difficulty levels: Easy (Green), Medium (Yellow), Hard (Red)
**And** Colors use ANSI escape codes
**And** Colors respect config setting for color_enabled

**Given** Terminal width is small
**When** I run `dsa list` in a narrow terminal (< 80 columns)
**Then** Table adapts by:
  - Hiding less important columns (e.g., ID)
  - Shortening column widths
  - Truncating long text with ellipsis
**And** Table remains readable and doesn't break layout

---

### Story 6.3: Implement CSV Export Format

As a **user**,
I want **to export data in CSV format**,
So that **I can analyze it in spreadsheets or data tools** (FR32, FR43).

**Acceptance Criteria:**

**Given** I want to export problem list as CSV
**When** I run `dsa list --format csv`
**Then** Output is valid CSV with headers:
```csv
ID,Title,Difficulty,Topic,Solved,FirstSolvedAt
two-sum,Two Sum,Easy,Arrays,true,2025-01-15
add-two-numbers,Add Two Numbers,Medium,Linked Lists,false,
```
**And** CSV follows RFC 4180 standard
**And** Values with commas are properly quoted
**And** Empty values are represented correctly

**Given** I want to export progress data as CSV
**When** I run `dsa export --format csv`
**Then** CSV includes all progress fields:
  - ProblemID, Title, Difficulty, Topic, Solved, FirstSolvedAt, TotalAttempts, LastAttemptedAt
**And** CSV can be imported into Excel or Google Sheets
**And** Dates are formatted as ISO 8601 (YYYY-MM-DD)

**Given** I want to save CSV to a file
**When** I run `dsa list --format csv > problems.csv`
**Then** CSV is written to the file without extra formatting
**And** File is a valid CSV that can be imported by spreadsheet tools

---

### Story 6.4: Implement Markdown Output Format

As a **user**,
I want **to generate Markdown tables from CLI output**,
So that **I can easily include results in documentation** (FR33).

**Acceptance Criteria:**

**Given** I want problem list in Markdown
**When** I run `dsa list --format markdown`
**Then** Output is a valid Markdown table:
```markdown
| ID | Title | Difficulty | Topic | Status |
|----|-------|------------|-------|--------|
| two-sum | Two Sum | Easy | Arrays | ✓ |
| add-two-numbers | Add Two Numbers | Medium | Linked Lists | ✗ |
```
**And** Table can be pasted directly into README.md
**And** Table renders correctly in GitHub, GitLab, and other Markdown viewers

**Given** I want status report in Markdown
**When** I run `dsa status --format markdown`
**Then** Output includes:
  - Heading: "# DSA Progress Report"
  - Summary statistics
  - Breakdown tables by difficulty and topic
  - Recent activity list
**And** Output is ready to commit to a progress tracking document

**Given** I want to generate a progress report
**When** I run `dsa report --format markdown --output progress.md`
**Then** A complete Markdown report is generated with:
  - Overall stats
  - Progress charts (ASCII or emoji-based)
  - Problem breakdown tables
  - Recent activity
**And** File can be committed to Git for tracking over time

---

### Story 6.5: Implement Customizable Output Templates

As a **user**,
I want **to define custom output templates**,
So that **I can format output exactly how I need it** (FR34, FR35).

**Acceptance Criteria:**

**Given** I want a custom problem list format
**When** I create a template file at `~/.dsa/templates/list.tmpl`:
```
{{range .Problems}}
- [{{if .Solved}}x{{else}} {{end}}] {{.Title}} ({{.Difficulty}})
{{end}}
```
**And** I run `dsa list --template list`
**Then** Output uses my custom template:
```
- [x] Two Sum (Easy)
- [ ] Add Two Numbers (Medium)
```
**And** Template uses Go text/template syntax (Architecture pattern)

**Given** I want to use template variables
**When** My template includes: `{{.Total}} problems ({{.Solved}} solved)`
**Then** Variables are substituted with actual values
**And** Template has access to all data fields

**Given** I want to share my template
**When** I run `dsa template export list --output mylist.tmpl`
**Then** The template is exported to the specified file
**And** Another user can import it with `dsa template import mylist.tmpl`

**Given** I want to list available templates
**When** I run `dsa template list`
**Then** I see all built-in and custom templates:
  - table (built-in)
  - json (built-in)
  - csv (built-in)
  - markdown (built-in)
  - list (custom)
**And** Description shows what each template does

---

### Story 6.6: Implement Progress Visualization (ASCII Charts)

As a **user**,
I want **to see visual progress charts in the terminal**,
So that **I can quickly understand my progress visually** (FR15).

**Acceptance Criteria:**

**Given** I want to visualize progress
**When** I run `dsa status --visual`
**Then** I see ASCII bar charts showing:
```
Progress by Difficulty:
Easy   [████████████████░░░░] 80% (16/20)
Medium [██████░░░░░░░░░░░░░░] 30% (12/40)
Hard   [███░░░░░░░░░░░░░░░░░] 15% (3/20)

Progress by Topic:
Arrays        [████████████░░░░░░░░] 60% (15/25)
Linked Lists  [██████░░░░░░░░░░░░░░] 30% (6/20)
Trees         [████░░░░░░░░░░░░░░░░] 20% (4/20)
```
**And** Progress bars use Unicode block characters (█ ░)
**And** Bars scale to terminal width automatically

**Given** I want a compact visualization
**When** I run `dsa status --sparkline`
**Then** I see sparkline charts using Unicode:
```
Last 30 days: ▁▂▃▅▇█▇▅▃▂▁▂▃▅▇█
```
**And** Sparklines show activity over time
**And** Useful for at-a-glance progress tracking

**Given** I want emoji-based visualization
**When** I run `dsa status --emoji`
**Then** Progress is shown with emoji indicators:
```
Progress: 🟩🟩🟩🟩🟩🟩🟩🟩⬜⬜ 80%
Easy:     ⭐⭐⭐⭐⭐ (5/5)
Medium:   ⭐⭐⭐⬜⬜ (3/5)
Hard:     ⭐⬜⬜⬜⬜ (1/5)
```
**And** Emoji makes output friendly and engaging

**Given** I want to save visualization to a file
**When** I run `dsa status --visual --output progress.txt`
**Then** The ASCII visualization is saved to the file
**And** File can be shared or committed to Git

---

## Epic 7: Motivation & Celebration System [Phase 2]

### Story 7.1: Implement Daily Streak Tracking

As a **user**,
I want **to track my daily problem-solving streak**,
So that **I stay motivated to practice consistently** (FR36, FR37).

**Acceptance Criteria:**

**Given** I solve a problem today
**When** I run `dsa test <problem-id>` and all tests pass
**Then** My streak counter increments by 1 if I solved problems yesterday
**And** If this is my first solve today, the streak starts at 1
**And** The database records today's activity date

**Given** I check my current streak
**When** I run `dsa status`
**Then** I see my current streak: "🔥 Current Streak: 7 days"
**And** If I haven't solved anything today, I see: "⚠️  Solve a problem today to continue your 7-day streak!"
**And** If my streak is broken, I see: "Streak: 0 days (last active: 2 days ago)"

**Given** I achieve milestone streaks
**When** My streak reaches 7, 14, 30, 60, 90, or 365 days
**Then** I see a special message:
  - 7 days: "🎉 One week streak! Keep going!"
  - 30 days: "🔥 One month streak! You're on fire!"
  - 365 days: "🏆 ONE YEAR STREAK! Incredible dedication!"
**And** Milestones are recorded in the database for history

**Given** I want to see my longest streak
**When** I run `dsa stats --streaks`
**Then** I see:
  - Current streak: X days
  - Longest streak: Y days (from date to date)
  - Total active days: Z
**And** This motivates me to beat my personal record

---

### Story 7.2: Implement Celebration Messages and ASCII Art

As a **user**,
I want **celebratory messages when I solve problems**,
So that **I feel a sense of accomplishment** (FR36, FR38).

**Acceptance Criteria:**

**Given** I solve an Easy problem
**When** All tests pass
**Then** I see a simple celebration: "✅ Great job! Easy problem solved."
**And** Message is encouraging but not overly dramatic

**Given** I solve a Medium problem
**When** All tests pass
**Then** I see: "🎉 Excellent! Medium problem conquered!"
**And** Message acknowledges the increased difficulty

**Given** I solve a Hard problem
**When** All tests pass
**Then** I see ASCII art celebration:
```
╔═══════════════════════════════════════╗
║  🏆  HARD PROBLEM SOLVED!  🏆       ║
║                                       ║
║  Outstanding work! You're crushing    ║
║  the tough challenges!                ║
╔═══════════════════════════════════════╗
```
**And** Celebration is memorable and special
**And** ASCII art is only shown for Hard problems or milestones

**Given** I solve my first problem of a topic
**When** I solve my first Tree problem
**Then** I see: "🌟 First Tree problem solved! New territory unlocked!"
**And** This highlights progress in new areas

**Given** I complete all problems in a topic
**When** I solve the last remaining Array problem
**Then** I see:
```
🎊 TOPIC MASTERED: Arrays 🎊
You've solved all 25 Array problems!
Next challenge: Try Trees or Graphs!
```
**And** Major achievements get special recognition

---

### Story 7.3: Implement Achievement System

As a **user**,
I want **to earn achievements for milestones**,
So that **I have goals to work toward** (FR37, FR39).

**Acceptance Criteria:**

**Given** The achievement system is implemented
**When** I complete specific milestones
**Then** I unlock achievements such as:
  - "First Steps" - Solve your first problem
  - "Getting Started" - Solve 10 problems
  - "Committed" - Solve 25 problems
  - "Dedicated" - Solve 50 problems
  - "Master" - Solve 100 problems
  - "Easy Rider" - Solve all Easy problems
  - "Rising Challenge" - Solve all Medium problems
  - "Hard Core" - Solve all Hard problems
  - "Week Warrior" - 7-day streak
  - "Month Master" - 30-day streak
  - "Completionist" - Solve all problems in library

**Given** I unlock an achievement
**When** I complete the required milestone
**Then** I see a special notification:
```
🏅 Achievement Unlocked: "Getting Started"
You've solved 10 problems! Keep up the great work!
```
**And** Achievement is permanently recorded in database

**Given** I want to view my achievements
**When** I run `dsa achievements`
**Then** I see a list of all achievements:
  - Unlocked achievements with dates (✓ colored green)
  - Locked achievements with requirements (🔒 colored gray)
  - Progress toward next achievement (e.g., "Dedicated: 35/50 problems")
**And** Output is visually appealing and motivating

**Given** I want to see achievement progress
**When** I run `dsa achievements --progress`
**Then** I see progress bars for in-progress achievements:
```
Dedicated (50 problems): [███████░░░░░] 35/50 (70%)
Month Master (30 days):  [████░░░░░░░░] 12/30 (40%)
```

---

### Story 7.4: Implement Encouraging Messages for Failures

As a **user**,
I want **encouraging messages when tests fail**,
So that **I stay motivated even when struggling** (FR38).

**Acceptance Criteria:**

**Given** Some tests fail on my first attempt
**When** I run `dsa test <problem-id>` and 2/5 tests pass
**Then** I see: "Good start! 2/5 tests passing. Keep debugging!"
**And** Message is positive and encouraging, not discouraging

**Given** Tests fail after multiple attempts
**When** I fail tests 3+ times on the same problem
**Then** I see helpful suggestions:
  - "Taking a break can help! Come back with fresh eyes."
  - "Consider reviewing the problem description again."
  - "Try adding print statements to debug your logic."
**And** Messages provide actionable advice

**Given** All tests were passing but now fail (regression)
**When** I break previously passing tests
**Then** I see: "⚠️  Tests were passing before. Check your recent changes."
**And** This helps identify when I introduced a bug

**Given** I'm close to passing
**When** I have 4/5 tests passing
**Then** I see: "Almost there! Just one more test to go. You've got this! 💪"
**And** Message emphasizes how close I am to success

---

### Story 7.5: Implement Progress Milestones and Notifications

As a **user**,
I want **to be notified of progress milestones**,
So that **I can celebrate my improvement** (FR40).

**Acceptance Criteria:**

**Given** I solve my 10th problem
**When** The problem is marked as solved
**Then** I see: "🎊 Milestone! You've solved 10 problems!"
**And** Milestone notifications appear for: 10, 25, 50, 100, 200, 500 problems

**Given** I complete 50% of a difficulty level
**When** I solve 15/30 Easy problems
**Then** I see: "Halfway there! 50% of Easy problems completed (15/30)"

**Given** I complete all problems in a difficulty level
**When** I solve the last Easy problem
**Then** I see:
```
🏆 ALL EASY PROBLEMS SOLVED! 🏆
Time to level up to Medium difficulty!
```

**Given** I improve my solve rate
**When** My problems-per-week increases
**Then** I see: "📈 You're on a roll! 5 problems this week vs 2 last week!"
**And** Positive trends are highlighted

**Given** I achieve a personal best
**When** I solve a problem faster than my previous best
**Then** I see: "⚡ New personal record! 15% faster than your previous best!"

---

### Story 7.6: Implement Motivational Quotes and Tips

As a **user**,
I want **to see motivational quotes and tips**,
So that **I stay inspired and learn best practices** (FR38).

**Acceptance Criteria:**

**Given** I start the CLI
**When** I run any command for the first time in a session
**Then** Occasionally (10% of the time) I see a motivational quote:
  - "Consistency beats intensity. Keep showing up!"
  - "Every expert was once a beginner. You're making progress!"
  - "The only way to get better at solving problems is to solve problems."
**And** Quotes are relevant to DSA practice and growth mindset

**Given** I struggle with a problem
**When** I fail tests 5+ times on the same problem
**Then** I see a helpful tip:
  - "Tip: Start with the brute force solution, then optimize."
  - "Tip: Draw out the problem on paper before coding."
  - "Tip: Check for edge cases: empty input, single element, duplicates."
**And** Tips provide practical problem-solving strategies

**Given** I want to disable motivational messages
**When** I run `dsa config set motivation_enabled false`
**Then** No quotes, tips, or celebration messages are shown
**And** Only factual output is displayed (for users who prefer minimal output)

**Given** I want to see a random tip
**When** I run `dsa tip`
**Then** I see a random DSA tip or best practice
**And** Tips cover: time complexity, space complexity, common patterns, debugging strategies

---

## Epic 8: Shell Completion & Enhanced UX [Phase 2]

### Story 8.1: Implement Bash Shell Completion

As a **user**,
I want **shell completion for bash**,
So that **I can quickly type commands without memorization** (FR46).

**Acceptance Criteria:**

**Given** I install shell completion for bash
**When** I run `dsa completion bash > /usr/local/etc/bash_completion.d/dsa`
**Then** A bash completion script is generated
**And** After sourcing, I can press TAB to autocomplete commands

**Given** Shell completion is installed
**When** I type `dsa s` and press TAB
**Then** Completions suggest: `solve`, `status`, `stats`, `submit`
**And** I can cycle through options with repeated TAB

**Given** I'm completing a problem ID
**When** I type `dsa solve tw` and press TAB
**Then** Problem IDs starting with "tw" are suggested (e.g., "two-sum")
**And** Completions pull from actual problem database

**Given** I'm completing command flags
**When** I type `dsa list --d` and press TAB
**Then** Flags are suggested: `--difficulty`, `--debug`
**And** Flag completions are context-aware per command

---

### Story 8.2: Implement Zsh Shell Completion

As a **user**,
I want **shell completion for zsh**,
So that **I can use autocomplete in my preferred shell** (FR46).

**Acceptance Criteria:**

**Given** I install shell completion for zsh
**When** I run `dsa completion zsh > ~/.zsh/completion/_dsa` and source it
**Then** Zsh completion works with TAB
**And** Descriptions appear alongside completions (zsh feature)

**Given** I use zsh completion
**When** I type `dsa list --difficulty ` and press TAB
**Then** Difficulty levels are suggested with descriptions:
  - easy — Beginner-friendly problems
  - medium — Intermediate challenges
  - hard — Advanced problems
**And** Descriptions help me choose the right option

---

### Story 8.3: Implement Fish Shell Completion

As a **user**,
I want **shell completion for fish**,
So that **I can use the CLI efficiently in fish shell** (FR46).

**Acceptance Criteria:**

**Given** I install shell completion for fish
**When** I run `dsa completion fish > ~/.config/fish/completions/dsa.fish`
**Then** Fish completion works automatically (fish auto-loads completions)

**Given** I use fish shell
**When** I start typing `dsa`
**Then** Commands appear as I type with inline suggestions
**And** Fish's rich completion UI enhances discoverability

---

## Epic 9: Advanced Analytics & Recommendations [Phase 3]

### Story 9.1: Implement Weak Area Identification

As a **user**,
I want **the CLI to identify my weak areas**,
So that **I know what to focus on improving** (FR47, FR48).

**Acceptance Criteria:**

**Given** I have solve history across topics
**When** I run `dsa analyze`
**Then** The CLI identifies weak areas:
  - Topics with low solve rates (< 40%)
  - Difficulty levels where I struggle (high fail rate)
  - Problem types I avoid (least attempted topics)
**And** Output includes specific recommendations

**Given** I have weak areas
**When** I run `dsa analyze --weak`
**Then** I see:
```
🔍 Weak Areas Identified:
1. Trees - 20% solved (3/15) - Recommendation: Start with binary tree basics
2. Graphs - 10% solved (1/10) - Recommendation: Review BFS/DFS fundamentals
3. Hard problems - 5% solved (1/20) - Recommendation: Build up from Medium first
```

**Given** I want practice suggestions based on weaknesses
**When** I run `dsa suggest`
**Then** The CLI recommends problems from my weak areas
**And** Suggestions start with easier problems in those topics to build confidence

---

### Story 9.2: Implement Spaced Repetition Scheduling

As a **user**,
I want **the CLI to suggest problems for review using spaced repetition**,
So that **I retain knowledge long-term** (FR48, FR49).

**Acceptance Criteria:**

**Given** I solved a problem 7 days ago
**When** I run `dsa review`
**Then** The CLI suggests that problem for review
**And** Spaced repetition intervals are: 1 day, 3 days, 7 days, 14 days, 30 days

**Given** Problems are due for review
**When** I run `dsa review --list`
**Then** I see all problems due for review today with last solved date
**And** Problems are sorted by priority (longest time since last review)

**Given** I review a problem
**When** I run `dsa review <problem-id>` and solve it successfully
**Then** The next review date is scheduled based on spaced repetition algorithm
**And** Review history is tracked in the database

---

### Story 9.3: Implement Performance Analytics

As a **user**,
I want **detailed analytics on my problem-solving performance**,
So that **I can track improvement over time** (FR50).

**Acceptance Criteria:**

**Given** I have historical solve data
**When** I run `dsa analytics`
**Then** I see:
  - Solve rate by difficulty over time (graph/chart)
  - Average attempts per problem
  - Most and least practiced topics
  - Streaks and consistency patterns
  - Time spent practicing (if tracked)

**Given** I want to compare performance across time periods
**When** I run `dsa analytics --compare-weeks`
**Then** I see week-over-week comparison:
  - This week: 5 problems solved
  - Last week: 3 problems solved
  - Change: +67% 📈

---

## Epic 10: Community & Extensibility [Phase 4]

### Story 10.1: Import Progress from Backup

As a **user**,
I want **to import progress data from a backup file**,
So that **I can restore my data or migrate between machines** (FR41, FR42).

**Acceptance Criteria:**

**Given** I have an exported progress file
**When** I run `dsa import --file progress.json`
**Then** All problems, progress, and history are imported into the database
**And** Existing data is merged intelligently (no duplicates)
**And** I see summary: "Imported 45 problems, 120 solutions, 60 progress records"

**Given** Import conflicts exist (same problem different data)
**When** I import data
**Then** The CLI asks how to resolve conflicts:
  - Keep existing data
  - Overwrite with imported data
  - Merge (keep newest timestamps)
**And** I can choose a strategy: `--strategy merge|overwrite|keep`

---

### Story 10.2: Implement Plugin System Foundation

As a **user**,
I want **the CLI to support plugins**,
So that **I can extend functionality with custom features** (FR51, FR52).

**Acceptance Criteria:**

**Given** Plugin directory exists at `~/.dsa/plugins/`
**When** I place a plugin binary/script in that directory
**Then** The plugin is discovered and loaded on CLI start

**Given** A plugin is loaded
**When** I run `dsa plugins list`
**Then** I see all installed plugins with: name, version, description, status (enabled/disabled)

**Given** I want to execute a plugin command
**When** I run `dsa plugin run <plugin-name> [args]`
**Then** The plugin executes with access to:
  - Problem database (read-only by default)
  - Configuration data
  - Standard input/output for user interaction

---

### Story 10.3: Implement LeetCode Problem Sync

As a **user**,
I want **to sync problems from LeetCode**,
So that **I can practice LeetCode problems in the CLI** (FR6, FR53).

**Acceptance Criteria:**

**Given** I want to import LeetCode problems
**When** I run `dsa leetcode sync`
**Then** The CLI fetches problem metadata from LeetCode API
**And** Problems are added to the local database
**And** Each problem includes: title, description, difficulty, topic, example test cases

**Given** LeetCode problems are synced
**When** I run `dsa list --source leetcode`
**Then** I see only problems imported from LeetCode
**And** Problems are marked with source: "leetcode"

**Given** I want to filter by LeetCode difficulty
**When** I run `dsa list --source leetcode --difficulty medium`
**Then** Only Medium LeetCode problems are shown

---

### Story 10.4: Implement Problem Sharing

As a **user**,
I want **to share custom problems with others**,
So that **the community can benefit from curated problem sets** (FR54).

**Acceptance Criteria:**

**Given** I created custom problems
**When** I run `dsa share <problem-id> --output problem.json`
**Then** Problem data is exported to a shareable JSON file including:
  - Problem metadata (title, difficulty, topic, description)
  - Test cases
  - Boilerplate code template
  - Solution (optional, if --include-solution flag is used)

**Given** I receive a shared problem
**When** I run `dsa import-problem --file problem.json`
**Then** The problem is added to my library
**And** I can solve it like any other problem

---

### Story 10.5: Implement Community Problem Repository

As a **user**,
I want **to browse and download problems from a community repository**,
So that **I can access curated problem collections** (FR54).

**Acceptance Criteria:**

**Given** A community problem repository exists (e.g., GitHub repo)
**When** I run `dsa community browse`
**Then** I see available problem packs:
  - "Top Interview Questions" (50 problems)
  - "Algorithm Patterns" (30 problems)
  - "Company-Specific: Google" (40 problems)
**And** Each pack shows: name, problem count, difficulty distribution, downloads

**Given** I want to install a problem pack
**When** I run `dsa community install "Top Interview Questions"`
**Then** All problems in the pack are downloaded and added to my library
**And** I see progress: "Downloading... 50/50 problems imported"

---

### Story 10.6: Implement Stats Comparison (Leaderboards)

As a **user**,
I want **to compare my stats with others (optionally)**,
So that **I can see how I rank and stay motivated** (FR54).

**Acceptance Criteria:**

**Given** I opt into leaderboard (privacy-respecting, opt-in only)
**When** I run `dsa leaderboard --opt-in`
**Then** My anonymized stats are uploaded to a community server
**And** I can view leaderboard rankings

**Given** I'm opted into leaderboard
**When** I run `dsa leaderboard`
**Then** I see rankings:
  - Top 10 users by problems solved
  - Top streaks
  - My rank and percentile
**And** All data is anonymized (usernames are hashed or pseudonymous)

**Given** I want to opt out
**When** I run `dsa leaderboard --opt-out`
**Then** My data is removed from the leaderboard
**And** Local tracking continues unaffected

---

**End of Epic Stories**
