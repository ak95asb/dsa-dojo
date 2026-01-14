# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Project foundation with Cobra CLI framework
- Database layer with GORM and SQLite support
- Workspace initialization command (`dsa init`)
- CI/CD pipeline with GitHub Actions
- Multi-platform binary distribution (Linux, macOS, Windows) via GoReleaser
- Comprehensive project documentation (README, CONTRIBUTING, CHANGELOG)

### Infrastructure
- GitHub Actions workflows for continuous integration
- Multi-platform testing (ubuntu, macos, windows)
- golangci-lint integration with 7 recommended linters
- Code coverage tracking via codecov
- Automated release workflow with GoReleaser

## [0.1.0] - 2025-12-11

### Added
- Initial project setup and structure
- Core CLI commands scaffolded: init, solve, test, status
- Local SQLite database for progress tracking
- GORM models: Problem, Solution, Progress
- Test-driven workflow foundation
- Cross-platform compatibility (macOS, Linux, Windows)

### Technical
- Go 1.23 support
- Standard Go project layout (cmd/, internal/, pkg/)
- Configuration management with Viper
- Error handling with proper exit codes
- In-memory SQLite testing support

[Unreleased]: https://github.com/empire/dsa/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/empire/dsa/releases/tag/v0.1.0
