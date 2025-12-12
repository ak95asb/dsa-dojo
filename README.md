```
╔═══════════════════════════════════════════════════════════════════════════╗
║                                                                           ║
║   ██████╗ ███████╗ █████╗       ██████╗  ██████╗      ██╗ ██████╗       ║
║   ██╔══██╗██╔════╝██╔══██╗      ██╔══██╗██╔═══██╗     ██║██╔═══██╗      ║
║   ██║  ██║███████╗███████║█████╗██║  ██║██║   ██║     ██║██║   ██║      ║
║   ██║  ██║╚════██║██╔══██║╚════╝██║  ██║██║   ██║██   ██║██║   ██║      ║
║   ██████╔╝███████║██║  ██║      ██████╔╝╚██████╔╝╚█████╔╝╚██████╔╝      ║
║   ╚═════╝ ╚══════╝╚═╝  ╚═╝      ╚═════╝  ╚═════╝  ╚════╝  ╚═════╝       ║
║                                                                           ║
║            🥋  Your Personal Data Structures & Algorithms Dojo  🥋        ║
║                      Train. Test. Track. Triumph.                        ║
║                                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-purple?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://github.com/ak95asb/dsa-dojo)
[![PRs Welcome](https://img.shields.io/badge/PRs-Welcome-ff69b4?style=for-the-badge)](CONTRIBUTING.md)

**[Features](#-features) • [Installation](#-installation) • [Quick Start](#-quick-start) • [Commands](#-commands) • [Contributing](#-contributing)**

</div>

---

## 🎮 What is DSA Dojo?

Welcome to your **personal coding dojo** - a retro-inspired, terminal-based CLI tool that transforms grinding LeetCode-style problems into an immersive, gamified experience. Built for developers who love the aesthetic of 80s computing and the discipline of deliberate practice.

**DSA Dojo** helps you:
- 📚 Practice 100+ curated DSA problems right from your terminal
- 🧪 Follow test-driven development with built-in test execution
- 📊 Track your progress with beautiful ASCII dashboards
- 🔥 Build daily solving streaks and earn achievements
- 💾 Export your stats in multiple formats (JSON, CSV, Markdown)

No browser tabs. No distractions. Just you, your terminal, and pure problem-solving flow.

---

## ✨ Features

### 🎯 Core Functionality
- **100+ Curated Problems**: Handpicked DSA challenges across arrays, trees, graphs, dynamic programming, and more
- **Test-Driven Practice**: Generate boilerplate code with test cases, solve with instant feedback
- **Multi-Format Output**: View results as tables, JSON, CSV, or Markdown
- **Smart Filtering**: Filter by difficulty, topic, solved status, or pick random challenges
- **Progress Tracking**: SQLite-backed persistence tracks every solve, attempt, and timestamp

### 🎨 Retro Experience
- **ASCII Art Dashboards**: Beautiful terminal UI with progress bars and stats
- **Keyboard-First**: Blazing fast CLI with shell completions (bash/zsh/fish)
- **Lightweight**: Single binary, no dependencies, runs anywhere Go runs
- **Offline-First**: All data stored locally, practice without internet

### 🧰 Developer-Friendly
- **Cobra CLI Framework**: Industry-standard command structure
- **GORM ORM**: Clean database abstractions
- **Viper Config**: Flexible YAML/JSON configuration
- **Watch Mode**: Auto-run tests on file save during practice
- **Editor Integration**: Opens your favorite $EDITOR for solving

---

## 📦 Installation

### Option 1: Build from Source (Recommended)
```bash
# Clone the repository
git clone https://github.com/ak95asb/dsa-dojo.git
cd dsa-dojo

# Build the binary
go build -o dsa

# Move to PATH (optional)
sudo mv dsa /usr/local/bin/

# Verify installation
dsa --version
```

### Option 2: Go Install
```bash
go install github.com/ak95asb/dsa-dojo@latest
```

### Option 3: Download Pre-built Binary
Download the latest release from [Releases](https://github.com/ak95asb/dsa-dojo/releases) page.

---

## 🚀 Quick Start

```bash
# Initialize your workspace (creates ~/.dsa/ directory with database)
dsa init

# List all problems
dsa list

# Filter problems by difficulty
dsa list --difficulty easy

# Get details about a specific problem
dsa show two-sum

# Pick a random problem to solve
dsa random --difficulty medium

# Start solving (generates boilerplate + tests)
dsa solve two-sum

# Run tests for your solution
dsa test two-sum

# Watch mode (auto-run tests on file save)
dsa watch two-sum

# Submit your solution (marks as solved)
dsa submit two-sum

# Check your progress dashboard
dsa status

# Export stats as CSV
dsa status --format csv > my-progress.csv
```

---

## 📖 Commands

### Problem Discovery
| Command | Description |
|---------|-------------|
| `dsa list` | List all problems with filters |
| `dsa show <slug>` | Display problem details with examples |
| `dsa random` | Pick a random problem |

### Practice Workflow
| Command | Description |
|---------|-------------|
| `dsa solve <slug>` | Generate boilerplate code and tests |
| `dsa test <slug>` | Run tests for your solution |
| `dsa watch <slug>` | Auto-run tests on file changes |
| `dsa submit <slug>` | Mark problem as solved |

### Progress & Stats
| Command | Description |
|---------|-------------|
| `dsa status` | View progress dashboard |
| `dsa export` | Export progress data (JSON/CSV) |

### Configuration
| Command | Description |
|---------|-------------|
| `dsa config get <key>` | Get config value |
| `dsa config set <key> <value>` | Set config value |
| `dsa init` | Initialize workspace |

### Output Formats
All commands support `--format` flag:
- `--format table` (default) - Pretty ASCII tables
- `--format json` - Machine-readable JSON
- `--format csv` - Spreadsheet-friendly CSV
- `--format markdown` - GitHub-friendly Markdown

---

## 🎨 Screenshots

### List Problems
```
┌──────────────────────────────────────────────────────────────────────┐
│ ID            │ Title                │ Difficulty │ Topic     │ ✓    │
├──────────────────────────────────────────────────────────────────────┤
│ two-sum       │ Two Sum              │ Easy       │ Arrays    │ ✓    │
│ add-two-num   │ Add Two Numbers      │ Medium     │ Linked    │      │
│ longest-sub   │ Longest Substring    │ Medium     │ Strings   │ ✓    │
└──────────────────────────────────────────────────────────────────────┘
Total: 100 problems (23 solved)
```

### Status Dashboard
```
DSA Progress Dashboard

Overall Progress: 23/100 problems solved (23%)

By Difficulty:
┌────────────────────────────────────────────────┐
│ Easy      │ 15/38   │ █████████░░░░░░  39%    │
│ Medium    │ 6/42    │ ███░░░░░░░░░░░░  14%    │
│ Hard      │ 2/20    │ ██░░░░░░░░░░░░░  10%    │
└────────────────────────────────────────────────┘

Recent Activity:
  ✓ Two Sum (Easy) - 2025-01-14
  ✓ Valid Parentheses (Easy) - 2025-01-13
  ✓ Binary Search (Easy) - 2025-01-12
```

---

## 🏗️ Architecture

```
dsa-dojo/
├── cmd/              # Cobra command definitions
│   ├── root.go       # Root command and global flags
│   ├── list.go       # List problems command
│   ├── solve.go      # Solve problem command
│   ├── status.go     # Status dashboard command
│   └── ...
├── internal/
│   ├── database/     # GORM models and migrations
│   ├── problem/      # Problem service layer
│   ├── progress/     # Progress tracking service
│   ├── output/       # Output formatters (JSON, CSV, Table)
│   └── generator/    # Code and test generators
├── data/
│   └── problems/     # Problem library (JSON)
└── main.go           # Application entry point
```

---

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Config Management**: [Viper](https://github.com/spf13/viper)
- **Database**: SQLite with [GORM](https://gorm.io/)
- **Testing**: Go's built-in `testing` package + integration tests
- **Output**: Custom ASCII table rendering, `encoding/json`, `encoding/csv`

---

## 🤝 Contributing

We welcome contributions! Whether it's:
- 🐛 Bug fixes
- ✨ New features
- 📝 Documentation improvements
- 🎨 UI/UX enhancements
- 📚 Adding new problems to the library

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit with descriptive messages
6. Push to your fork
7. Open a Pull Request

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Inspired by the retro computing aesthetic of the 1980s
- Problem set curated from classic DSA interview questions
- Built with love for terminal enthusiasts and keyboard warriors

---

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/ak95asb/dsa-dojo/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/ak95asb/dsa-dojo/discussions)
- 📧 **Contact**: Open an issue or discussion

---

<div align="center">

**Made with 💜 and ☕ by [ak95asb](https://github.com/ak95asb)**

*Train hard. Code harder. Level up your DSA skills in the dojo.*

```
█▀▀ █▀█ █▀▄ █▀▀   █ █▄░█   ▀█▀ █░█ █▀▀   █▀▄ █▀█ ░░█ █▀█
█▄▄ █▄█ █▄▀ ██▄   █ █░▀█   ░█░ █▀█ ██▄   █▄▀ █▄█ █▄█ █▄█
```

</div>
