# Contributing to DSA Dojo 🥋

```
╔══════════════════════════════════════════════════════════════╗
║  ████████╗██╗  ██╗ █████╗ ███╗   ██╗██╗  ██╗███████╗        ║
║  ╚══██╔══╝██║  ██║██╔══██╗████╗  ██║██║ ██╔╝██╔════╝        ║
║     ██║   ███████║███████║██╔██╗ ██║█████╔╝ ███████╗        ║
║     ██║   ██╔══██║██╔══██║██║╚██╗██║██╔═██╗ ╚════██║        ║
║     ██║   ██║  ██║██║  ██║██║ ╚████║██║  ██╗███████║        ║
║     ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝        ║
║                                                              ║
║         For joining the dojo and helping us grow! 💜         ║
╚══════════════════════════════════════════════════════════════╝
```

We're excited that you're interested in contributing to DSA Dojo! This document outlines the process and guidelines for contributing.

## 🎯 Ways to Contribute

- **🐛 Report Bugs**: Found a bug? Let us know!
- **✨ Suggest Features**: Have an idea? We'd love to hear it!
- **📚 Add Problems**: Expand our problem library with quality DSA challenges
- **📝 Improve Docs**: Help make our documentation clearer
- **🎨 Enhance UI**: Make the terminal experience even better
- **🧪 Write Tests**: Improve test coverage
- **🔧 Fix Issues**: Tackle open issues

## 🚀 Getting Started

### 1. Fork & Clone

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/YOUR_USERNAME/dsa-dojo.git
cd dsa-dojo
```

### 2. Set Up Development Environment

```bash
# Install dependencies
go mod download

# Build the project
go build -o dsa

# Run tests
go test ./...

# Run integration tests
go test ./cmd -v
```

### 3. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-number-description
```

**Branch Naming Convention**:
- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring
- `test/description` - Test improvements

## 📋 Development Guidelines

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused (single responsibility principle)

### Testing

- **Always** write tests for new features
- Maintain or improve test coverage
- Run `go test ./...` before committing
- Use table-driven tests where appropriate

Example:
```go
func TestFormatDateISO(t *testing.T) {
    tests := []struct {
        name     string
        date     *time.Time
        expected string
    }{
        {"valid date", timePtr(time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)), "2025-01-15"},
        {"nil date", nil, ""},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := formatDateISO(tt.date)
            if result != tt.expected {
                t.Errorf("expected '%s', got '%s'", tt.expected, result)
            }
        })
    }
}
```

### Commit Messages

Write clear, descriptive commit messages:

```
<type>: <subject>

<body (optional)>

<footer (optional)>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code restructuring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples**:
```
feat: add markdown output format for status command

- Implemented MarkdownFormatter for status stats
- Added --format markdown flag
- Updated tests and documentation
```

```
fix: resolve CSV date formatting for nil timestamps

Fixes #42
```

## 🐛 Reporting Bugs

Before creating a bug report:
1. **Search existing issues** to avoid duplicates
2. **Update to latest version** and verify the bug persists
3. **Gather information**: OS, Go version, command that failed

When creating a bug report, include:
- **Clear title**: "Command fails when..." not "It's broken"
- **Steps to reproduce**
- **Expected behavior**
- **Actual behavior**
- **Environment**: OS, Go version, terminal
- **Error messages**: Full stack traces if applicable

**Example**:
```markdown
**Bug**: `dsa list --format csv` panics with nil pointer

**Steps to Reproduce**:
1. Run `dsa list --format csv`
2. Observe panic

**Expected**: CSV output
**Actual**: Panic with error: `panic: runtime error: invalid memory address`

**Environment**:
- OS: macOS 14.2
- Go: 1.21.5
- Terminal: iTerm2

**Error**:
```
panic: runtime error: invalid memory address or nil pointer dereference
    at cmd/output_csv.go:42
```
```

## ✨ Suggesting Features

We love feature ideas! When suggesting:
- **Check existing issues** first
- **Explain the use case**: Why is this needed?
- **Describe the solution**: How should it work?
- **Consider alternatives**: Are there other approaches?
- **Show examples**: Mock up the command/output

## 📚 Adding Problems to the Library

Want to contribute a new problem?

1. **Check for duplicates** in `data/problems/`
2. **Create problem JSON** following this structure:

```json
{
  "id": "unique-slug",
  "title": "Problem Title",
  "difficulty": "easy|medium|hard",
  "topic": "arrays|strings|trees|etc",
  "description": "Clear problem description",
  "examples": [
    {
      "input": "...",
      "output": "...",
      "explanation": "..."
    }
  ],
  "constraints": ["List of constraints"],
  "hints": ["Optional hints"],
  "solution_template": "Go code template",
  "test_cases": [
    {
      "input": "...",
      "expected": "...",
      "description": "..."
    }
  ]
}
```

3. **Submit a PR** with your problem

**Quality Guidelines**:
- Problem should be well-explained and unambiguous
- Include 3-5 test cases (edge cases + normal cases)
- Provide a working solution template
- Difficulty should match complexity accurately

## 🔄 Pull Request Process

### Before Submitting

- [ ] Code builds without errors (`go build`)
- [ ] All tests pass (`go test ./...`)
- [ ] Code follows Go conventions (`gofmt`, `golint`)
- [ ] New code has tests
- [ ] Documentation updated (if needed)
- [ ] Commit messages are clear
- [ ] Branch is up to date with `main`

### Submitting

1. **Push your branch** to your fork
2. **Open a PR** against `main` branch
3. **Fill out the PR template** completely
4. **Link related issues** (e.g., "Fixes #42")
5. **Request review** (optional, but encouraged)

### After Submitting

- **Respond to feedback** promptly
- **Make requested changes** in new commits
- **Keep the PR focused**: One feature/fix per PR
- **Be patient**: Reviews may take a few days

### PR Title Format

```
<type>: <description>

Examples:
feat: add benchmark command for performance testing
fix: resolve panic in CSV export with empty data
docs: update installation instructions for Windows
```

## 🧪 Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestFormatDateISO ./cmd

# Run integration tests only
go test ./cmd -run Integration

# Skip integration tests (for quick iteration)
go test -short ./...
```

## 📁 Project Structure

```
dsa-dojo/
├── cmd/              # Command implementations
│   ├── root.go       # Root command
│   ├── list.go       # List command
│   ├── *_test.go     # Test files
│   └── ...
├── internal/
│   ├── database/     # Database models & migrations
│   ├── problem/      # Problem service layer
│   ├── progress/     # Progress tracking
│   ├── output/       # Output formatters
│   └── generator/    # Code generators
├── data/
│   └── problems/     # Problem library (JSON)
├── docs/
│   └── sprint-artifacts/  # BMAD workflow artifacts
└── main.go
```

## 💬 Code of Conduct

This project follows a Code of Conduct to ensure a welcoming environment. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## ❓ Questions?

- **General questions**: [GitHub Discussions](https://github.com/ak95asb/dsa-dojo/discussions)
- **Bug reports**: [GitHub Issues](https://github.com/ak95asb/dsa-dojo/issues)
- **Feature requests**: [GitHub Discussions](https://github.com/ak95asb/dsa-dojo/discussions)

---

## 🙏 Recognition

All contributors will be recognized in our README. Thank you for helping make DSA Dojo awesome!

---

<div align="center">

```
█░█ ▄▀█ █▀█ █▀█ █▄█   █▀▀ █▀█ █▀▄ █ █▄░█ █▀▀ ▄█
█▀█ █▀█ █▀▀ █▀▀ ░█░   █▄▄ █▄█ █▄▀ █ █░▀█ █▄█ ░▄
```

**Train hard. Code harder. Welcome to the dojo.** 🥋

</div>
