---
stepsCompleted: [1, 2, 3, 4, 7, 8, 9, 10]
inputDocuments:
  - '/Users/noi03_ajaysingh/Documents/LearnGo/dsa/docs/analysis/brainstorming-session-2025-12-10.md'
documentCounts:
  briefs: 0
  research: 0
  brainstorming: 1
  projectDocs: 0
workflowType: 'prd'
lastStep: 11
project_name: 'dsa'
user_name: 'Empire'
date: '2025-12-10'
---

# Product Requirements Document - dsa

**Author:** Empire
**Date:** 2025-12-10

## Executive Summary

**dsa** is a CLI-based DSA practice platform designed to make algorithm practice engaging, effective, and habit-forming for developers who live in the terminal. Built on the proven simplicity of kata-machine's test-driven workflow, it adds the critical missing layer: **motivation through encouragement and celebrating progress**.

**The Problem:** Developers want to build DSA skills, but existing solutions fall short. Kata-machine proves the CLI + local-first model works, but offers no tracking, no motivation, and no sense of progress. LeetCode has gamification but forces you into a browser and owns your data. No platform understands that sustainable practice isn't about grinding—it's about celebrating little victories that compound into mastery.

**The Solution:** A CLI tool that respects developer workflows (local-first, offline, editor-agnostic, multi-language support starting with Go) while adding intelligent tracking and encouragement. When you solve a problem, you're celebrated. When you maintain a streak, you see it. When you have weak areas, you get personalized suggestions—not judgment.

### What Makes This Special

Other platforms track metrics. This platform celebrates victories. The core insight: developers don't need more data about their practice—they need **encouragement that makes showing up rewarding**. Little wins celebrated consistently create habits. Habits create mastery.

The platform combines:
- **kata-machine's simplicity:** Test-driven workflow, scaffolding, local-first
- **Learning science:** Spaced repetition, weak area identification, adaptive scheduling
- **Go tooling excellence:** Native `go test` integration, benchmarking, race detection as first-class features
- **Developer respect:** Offline-capable, import/export your data, works with your tools

**Vision:** Transform DSA practice from a grind into a rewarding habit by celebrating progress and making developers feel accomplished—not just productive.

## Project Classification

**Technical Type:** CLI Tool
**Domain:** General (Developer Productivity/Learning)
**Complexity:** Low-Medium
**Project Context:** Greenfield - new project

**Classification Rationale:**

This is a command-line developer tool (detected signals: CLI commands like `dsa init`, `dsa status`, `dsa solve`; terminal-based workflow; scriptable interface; shell integration requirements). The domain is general developer productivity focused on skill building—outside regulated industries like healthcare or fintech.

The CLI tool classification means the PRD will focus on command structure, output formats, configuration methods, and scripting support rather than visual UI or touch interactions. Multi-language support (Go flagship, expanding to Rust/Python/TypeScript) requires language-specific tooling integration (`go test`, `cargo test`, `pytest`).

**Key Technical Characteristics:**
- Local-first architecture (SQLite embedded database)
- Offline-capable with import/export
- Language-specific test framework integration
- Go as flagship with premium features first

## Success Criteria

### User Success

**Primary Success Metric:** Developers willingly spend 30 minutes daily because the experience is rewarding, not because they "should."

**User Success Indicators:**
- **Engagement Quality:** The 30-minute session feels productive and satisfying, not like grinding through obligations
- **Skill Transfer:** Users notice DSA concepts clicking in their actual coding work—recognizing patterns, applying algorithms naturally
- **Learning Curve Satisfaction:** The challenge feels like leveling up in a game (rewarding struggle) rather than frustrating confusion
- **Habit Formation:** Users return consistently driven by encouragement + visible progress, not guilt or obligation
- **"Aha!" Moments:** Regular breakthrough moments where concepts finally make sense through practice
- **Pattern Recognition:** Users begin seeing connections across problems and applying learned approaches to new challenges

**Success is NOT:** Maximizing problem count or speed. Success is building real, transferable coding skills through an experience developers actively choose to repeat.

### Business Success

**Business Model:** Open source community project

**Primary Business Metric:** Creating a tool valuable enough that the creator (and developers like them) use it for 30 minutes daily.

**Business Success Indicators:**
- **Dogfooding Test:** Creator genuinely uses the platform daily for personal skill development
- **Community Validation:** Other developers adopt, contribute, and recommend it organically
- **GitHub Activity:** Stars, forks, pull requests, and active issues indicate growing community engagement
- **Word-of-Mouth Growth:** Developers discover it through peer recommendations, not marketing
- **Sustained Engagement:** 30-day and 90-day retention metrics show it's not just novelty—it's genuinely useful
- **Contribution Diversity:** Community adds problems, improves features, expands language support
- **Problem Library Growth:** Organic expansion of the problem set through community contributions

**Success Measurement Timeline:**
- **3 months:** Creator using it daily + 10-20 other active users providing feedback
- **6 months:** 100+ GitHub stars, 5+ community contributors, consistent daily active users
- **12 months:** 500+ stars, active problem contributions, multi-language support emerging from community

### Technical Success

**Primary Technical Metric:** The platform respects developer time and feels like a professional tool, not a toy.

**Technical Success Requirements:**
- **Fast Startup:** CLI commands execute instantly (<500ms for most operations)
- **Offline Reliability:** Works completely offline; no network dependency breaks the experience
- **Local-First Integrity:** All data stored locally, no server failures kill user streaks or progress
- **Native Tooling Integration:** `go test` works exactly as developers expect, zero friction
- **Authentic Encouragement:** Celebration and motivation features feel genuine, not gimmicky or patronizing
- **Data Ownership:** Users can export/import their progress—no platform lock-in
- **Extensibility:** Architecture supports community additions (new languages, problems, features)
- **Cross-Platform:** Works on macOS, Linux, Windows without platform-specific quirks

**Technical Quality Bars:**
- Test runner integration feels native (no "wrapper" friction)
- Scaffolding generation produces clean, idiomatic code
- Progress tracking is accurate and reliable
- CLI output is clear, colored appropriately, and helpful
- Database operations never block or slow down the experience

### Measurable Outcomes

**User-Level Metrics:**
- Daily active users with 7+ day streaks (Phase 2)
- Problem completion rates across Easy/Medium/Hard difficulties
- Time to first "completed problem" (onboarding friction indicator)
- Return rate: % of users who come back after first session
- Session duration: targeting 20-40 minute sessions consistently

**Community-Level Metrics:**
- GitHub stars and fork count
- Pull request volume and contributor diversity
- Problem library growth (community-submitted problems)
- Issue engagement (active discussions, feature requests)
- Documentation improvements from community

**Technical-Level Metrics:**
- CLI command execution time (<500ms target)
- Test runner integration success rate
- Data integrity (zero progress loss incidents)
- Cross-platform compatibility (works on all major OSes)
- Setup time from clone to first problem (<5 minutes)

## Product Scope

### MVP - Minimum Viable Product (Phase 1: Core DX Foundation)

**Goal:** Prove the technical foundation works—test-driven workflow, Go tooling integration, and basic CLI commands function flawlessly.

**Must Ship Together:**
1. **Test-Driven Workflow:** Problems come with pre-written tests; developers make them pass
2. **Scaffolding Generator:** Automatic boilerplate creation for new problems
3. **Go Tooling Integration:** Native `go test` execution with proper test framework integration
4. **Core CLI Commands:**
   - `dsa init` - Initialize workspace
   - `dsa solve [problem]` - Start working on a problem
   - `dsa test` - Run tests
   - `dsa status` - Show current progress (basic version)
5. **Local Storage:** SQLite embedded database for progress tracking and problem metadata
6. **Problem Library:** 10-15 curated problems across Easy/Medium/Hard difficulties covering core DSA topics (Arrays, Linked Lists, Trees, Sorting)

**Success Criteria for Phase 1:**
- Creator can use it for 30 minutes to practice a problem without friction
- Test workflow feels native to Go developers
- Setup takes <5 minutes from clone to first problem attempt
- Core mechanics validated: "Does the foundation work?"

**What's NOT in MVP:**
- Advanced tracking/motivation features (streaks, celebrations)
- Spaced repetition scheduling
- Multi-language support
- Import/export functionality

### Growth Features (Phase 2: The Excitement Layer)

**Goal:** Add the motivation and tracking system that transforms practice from "nice to have" into a daily habit.

**Priority Features (Ship immediately after MVP validation):**
1. **Tracking System:** Enhanced `dsa status` dashboard showing:
   - Current streak (consecutive days with at least one problem solved)
   - Problems solved by difficulty (Easy/Medium/Hard counts)
   - Topic coverage (Arrays: 5/15, Trees: 2/10, etc.)
   - Recent activity timeline
2. **Streak Counter:** Daily streak tracking with visual feedback (🔥 7-day streak!)
3. **Celebration Moments:** ASCII art and encouraging messages when:
   - Completing first problem
   - Solving a Hard problem
   - Reaching streak milestones (7, 30, 100 days)
   - Completing all problems in a topic
4. **Easy/Medium/Hard Tracking:** Progress visualization by difficulty level
5. **Goal Setting:** Set personal practice goals and track progress

**Success Criteria for Phase 2:**
- Streaks drive daily return behavior
- Celebrations feel authentic and motivating (not annoying)
- `dsa status` becomes the first command users run each session
- 30/90-day retention improves measurably vs Phase 1

### Growth Features (Phase 3: Smart Learning)

**Goal:** Add learning science principles (spaced repetition, weak area identification) to make practice more effective.

**Features:**
1. **Spaced Repetition Scheduling:** Intelligent problem review suggestions based on:
   - Time since last attempt
   - Success/failure history
   - Topic weakness patterns
2. **Weak Area Identification:** Analytics highlighting topics needing more practice
3. **Personalized Review Suggestions:** "Review these 3 problems from 2 weeks ago"
4. **Performance Analytics:** Track time spent, success rates, improvement trends
5. **Freshness Indicators:** Show when problems were last solved ("Last solved: 14 days ago")

**Success Criteria for Phase 3:**
- Users demonstrate improved retention on reviewed problems
- Weak area suggestions feel accurate and helpful
- Practice sessions become more targeted and efficient

### Vision (Phase 4: Polish & Expansion)

**Goal:** Expand platform reach and add features that enhance but aren't essential to core value.

**Future Features:**
1. **Import/Export Functionality:**
   - Export progress as JSON/CSV
   - Share achievements with peers
   - Multi-machine sync via import/export
2. **Multi-Language Support:**
   - Rust (native `cargo test` integration)
   - Python (pytest integration)
   - TypeScript (Jest integration)
   - Each language gets flagship-quality treatment like Go
3. **Advanced Go Challenges:**
   - Concurrency-focused katas (goroutines, channels, race conditions)
   - Memory optimization problems ("solve without allocating")
   - Benchmark-driven challenges (performance as first-class metric)
4. **Community Problem Contributions:** Enable community to submit and curate problems
5. **Custom Problem Sets:** Users create themed problem collections

**Success Criteria for Phase 4:**
- Multi-language support attracts new developer communities
- Community contributions become primary source of new problems
- Platform becomes "the CLI DSA tool" across multiple ecosystems

## User Journeys

### Journey 1: First-Time User - Alex Discovers dsa

**Alex is a mid-level backend engineer** who's been putting off DSA practice for months. They know they should be better at algorithms—interview prep aside, they want to actually *understand* tree traversals and dynamic programming instead of fumbling through them. They've tried Leetcode, but the browser-based experience feels disconnected from their actual coding workflow. They spend their days in the terminal, and switching to a web UI for practice feels... wrong.

Late on a Tuesday, Alex sees a GitHub post about **dsa**—a CLI tool for DSA practice built for Go developers. "Wait, CLI-based? I can use Vim? And it tracks progress?" They clone the repo.

Five minutes later, they run `dsa init` and the tool scaffolds a clean workspace. They type `dsa solve arrays/two-sum` and boom—a Go file opens in their editor with tests already written and a clear TODO. They write their solution, run `dsa test`, and watch the terminal light up green. It feels like normal coding, not "doing exercises."

The breakthrough comes when they run `dsa status` and see "Problems solved: 1/15 | Easy: 1 | Keep going!" It's simple, but it's *their* progress, stored locally, no account needed. Alex spends 30 minutes solving another problem that night. The next morning, they do it again.

### Journey 2: Daily Practicing User - Alex's New Routine (30 Days Later)

**It's been 30 days since Alex started using dsa.** What began as "I should probably practice" has become part of their morning routine—right after coffee, before checking email. The ritual is simple: open the terminal, run `dsa status`, see the streak counter: **🔥 30-day streak!**

That streak number hits different than any GitHub contribution graph. This is *their* commitment, stored locally, no one watching except themselves. The platform celebrates when they solve a Hard problem—ASCII art appears with a message: "You just solved a Hard problem before 9am. You're either a genius or you haven't slept. Either way: 🎉"

Some days are great—they knock out two problems and feel unstoppable. Other days, they only solve one Easy problem to keep the streak alive. **The tool doesn't judge.** It just shows: "Solved today ✓ | Streak continues | A commit a day keeps the impostor syndrome away."

The real magic happens three weeks in when Alex is debugging a production issue and suddenly recognizes it as a graph traversal problem. The solution clicks instantly because they *just practiced this pattern two days ago*. They mutter "I know this one!" to their rubber duck, who remains unimpressed but supportive.

The breakthrough moment comes when `dsa status` shows: "Weak area: Trees (2/12 completed) | Suggested: Review binary-tree problems. (No, not that kind of tree. The data structure.)" Alex realizes the platform is quietly learning their patterns and helping them strategically.

When they finally conquer that brutal linked-list reversal problem, the terminal displays: "LINKED LIST: REVERSED ✓ | Your pointers are finally pointing in the right direction. Unlike your life choices, but we don't judge."

Six months later, Alex hasn't missed a day. Not because they're obsessed—because 30 minutes of practice feels rewarding, not like homework.

### Journey 3: Contributor - Jordan Gives Back

**Jordan has been using dsa for four months** and loves it. They've crushed 87 problems, their friends are jealous of their 120-day streak, and they actually *understand* dynamic programming now (most days, anyway).

One Saturday morning, Jordan realizes: "Wait, this is open source. I could add Rust support!" They've been learning Rust and think "what better way to practice than making a tool for practicing?" (The irony is not lost on them.)

They check the GitHub repo. Clean README, clear contribution guidelines, friendly maintainer responses to issues. Jordan opens `CONTRIBUTING.md` and sees: "Want to add a language? Here's the abstraction pattern. Want to add problems? Here's the format. Want to fix a bug? You're a hero. ❤️"

Jordan starts with something small—adding three new array problems they wish existed. They follow the problem template, write tests, submit a PR. Two days later, the maintainer merges it with a comment: "These problems are 🔥. Thank you!"

**That feels good.** Really good. Better than any corporate code review. Jordan's problems are now helping other developers practice. They're part of the tool that helped them.

Three months later, Jordan has contributed Rust support (with `cargo test` integration working beautifully), written 15 problems, and helped review other PRs. When new users discover dsa, they're using Jordan's problems. When they run `dsa solve rust/hash-table-basics`, they're running Jordan's code.

The platform's celebratory message when Jordan hit their 100th contribution: "You've contributed 100+ times. At this point, you're basically maintaining this repo. Want commit access? (Seriously, let's talk.)"

Jordan doesn't just use dsa anymore. Jordan *builds* dsa. And that's how an open source community grows—one problem, one language, one contributor at a time.

### Journey Requirements Summary

These three journeys reveal the following capability requirements:

**From Journey 1 (First-Time User - Onboarding):**
- **Fast Setup:** `dsa init` command creates workspace in <5 minutes
- **Core CLI Commands:** init, solve, test, status with intuitive syntax
- **Editor Integration:** Opens problem files in user's preferred editor (Vim, VSCode, etc.)
- **Test-Driven Workflow:** Pre-written tests, scaffolded solution files
- **Scaffolding Generator:** Automatic boilerplate creation for problems
- **Basic Progress Tracking:** Problem completion counts, difficulty breakdowns
- **Local Storage:** SQLite database for progress, no cloud dependency
- **Problem Library:** Curated problems across difficulties and topics

**From Journey 2 (Daily Practice - Habit Loop):**
- **Streak Tracking:** Daily consecutive practice tracking with visual indicators
- **Enhanced Status Dashboard:** Comprehensive view of progress, streaks, weak areas
- **Celebration System:** ASCII art, encouraging messages with programming humor
- **Weak Area Identification:** Analytics showing which topics need more practice
- **Personalized Suggestions:** Intelligent problem recommendations based on history
- **Difficulty Tracking:** Easy/Medium/Hard completion counts and visualization
- **Motivational Messaging:** Authentic, humorous encouragement that feels genuine
- **Progress Persistence:** Reliable local data storage, zero progress loss

**From Journey 3 (Contributor - Community Growth):**
- **Contribution Documentation:** Clear CONTRIBUTING.md with guidelines and templates
- **Problem Templates:** Standardized format for community problem submissions
- **Architecture Documentation:** Clear extensibility patterns for languages and features
- **Language Plugin System:** Abstraction layer supporting multi-language expansion
- **Community Recognition:** Contributor credits, acknowledgment in problems
- **Local Dev Setup:** Easy contributor onboarding, test infrastructure
- **Testing Infrastructure:** Validation for community-submitted problems
- **Maintainer Tooling:** PR review process, quality gates, merge workflows

**Cross-Journey Requirements:**
- **Offline Capability:** All three journeys work without internet connection
- **Data Ownership:** Users control their data (local storage, export/import)
- **Cross-Platform:** Works on macOS, Linux, Windows consistently
- **Performance:** Fast command execution (<500ms for most operations)
- **Go-First Excellence:** Native `go test` integration, idiomatic patterns

## CLI Tool Specific Requirements

### Project-Type Overview

**dsa** is a command-line interface tool designed for terminal-native developers. The CLI architecture prioritizes:
- **Scriptability:** Commands work in scripts, pipes, and CI environments
- **Flexibility:** Hybrid configuration with sensible defaults
- **Professional UX:** Shell completion, colored output, machine-parseable formats
- **UNIX Philosophy:** Do one thing well, compose with other tools

### Command Structure

**Core Commands (MVP):**
```bash
dsa init                    # Initialize workspace
dsa solve [problem]         # Start working on a problem
dsa test                    # Run tests for current problem
dsa status                  # Show progress dashboard
```

**Command Design Principles:**
- **Verb-noun structure:** `dsa [verb] [noun]` (e.g., `dsa solve arrays/two-sum`)
- **Scriptable by default:** All commands work in non-interactive mode
- **Composable:** Commands chain with `&&`, pipe with `|`, redirect with `>`
- **Fast execution:** <500ms for most operations
- **Clear exit codes:** 0 for success, non-zero for errors (enables `&&` chaining)

**Future Commands (Post-MVP):**
```bash
dsa review                  # Show spaced repetition suggestions
dsa stats                   # Detailed analytics
dsa export                  # Export progress data
dsa import                  # Import progress data
```

### Output Formats

**Dual Output Strategy:**

**1. Human-Friendly (Default):**
- Colored terminal output (using ANSI codes)
- ASCII art for celebrations
- Formatted tables for status
- Progress bars for long operations
- Emoji indicators (🔥 for streaks, ✓ for completion)

**Example:**
```bash
$ dsa status
🔥 30-day streak!
─────────────────────────────────────
Problems Solved: 47/15
  Easy:   20 ████████████░░░░░░░░
  Medium: 18 █████████░░░░░░░░░░░
  Hard:    9 ████░░░░░░░░░░░░░░░░

Weak Areas: Trees (2/12)
Suggested: Review binary-tree problems
```

**2. Machine-Parseable (--json flag):**
- JSON output for scripting
- Structured, consistent schema
- No colors, no formatting
- Enables programmatic use

**Example:**
```bash
$ dsa status --json
{
  "streak": 30,
  "problems_solved": 47,
  "problems_total": 15,
  "by_difficulty": {
    "easy": 20,
    "medium": 18,
    "hard": 9
  },
  "weak_areas": ["trees"],
  "suggested_problems": ["binary-tree-traversal"]
}
```

**Output Requirements:**
- Respect `NO_COLOR` environment variable
- Support `--quiet` flag for minimal output
- Error messages go to stderr, data to stdout
- Colored output auto-disabled in non-TTY contexts

### Configuration Schema

**Hybrid Configuration Approach:**

**Priority Order (highest to lowest):**
1. **CLI flags:** `dsa solve --editor=vim`
2. **Environment variables:** `DSA_EDITOR=vim`
3. **Config file:** `~/.dsa/config.yaml`
4. **Defaults:** Built-in sensible defaults

**Config File Location:**
```
~/.dsa/config.yaml          # Global config
.dsa/config.yaml            # Project-specific (overrides global)
```

**Config File Schema:**
```yaml
# ~/.dsa/config.yaml
editor: "vim"               # Default: $EDITOR or vim
difficulty: "all"           # Filter: all, easy, medium, hard
output_format: "human"      # human or json
verbosity: "normal"         # quiet, normal, verbose
language: "go"              # go (rust/python/typescript future)

# Celebration preferences
celebrations: true          # Enable/disable ASCII art
humor: true                 # Enable/disable programming jokes
```

**Environment Variables:**
```bash
DSA_EDITOR=code             # Override editor
DSA_DIFFICULTY=medium       # Filter difficulty
DSA_OUTPUT=json             # Override output format
DSA_VERBOSITY=quiet         # Override verbosity
```

**CLI Flags:**
```bash
--editor=vim                # Override editor
--difficulty=hard           # Filter difficulty
--json                      # Output as JSON
--quiet                     # Minimal output
--verbose                   # Detailed output
```

### Scripting Support

**Scriptability Requirements:**

**1. Exit Codes:**
```bash
0   # Success
1   # General error
2   # Command usage error
3   # Test failure
4   # Problem not found
```

**2. Non-Interactive Mode:**
- Detect TTY vs non-TTY automatically
- No prompts in non-interactive mode
- Use flags for decisions (e.g., `--force` to skip confirmations)

**3. Composability Examples:**
```bash
# Chain commands
dsa solve arrays/two-sum && dsa test && dsa status

# Pipe output
dsa status --json | jq '.streak'

# Conditional execution
dsa test || echo "Tests failed, streak broken"

# CI/CD usage
#!/bin/bash
dsa solve "$PROBLEM" --quiet
dsa test --json > test-results.json
```

**4. Performance:**
- Cold start: <500ms
- Warm commands: <100ms
- Database operations: Non-blocking
- No network dependencies

### Shell Completion (MVP)

**Completion Support:**
- **Bash:** `dsa completion bash`
- **Zsh:** `dsa completion zsh`
- **Fish:** `dsa completion fish`

**Completion Scope (MVP):**
```bash
# Command completion
dsa <TAB>
  init    solve    test    status    completion

# Problem completion
dsa solve <TAB>
  arrays/    linked-lists/    trees/    sorting/

# Problem name completion
dsa solve arrays/<TAB>
  two-sum    three-sum    merge-intervals
```

**Installation:**
```bash
# Bash
dsa completion bash > /etc/bash_completion.d/dsa

# Zsh
dsa completion zsh > "${fpath[1]}/_dsa"

# Fish
dsa completion fish > ~/.config/fish/completions/dsa.fish
```

### Implementation Considerations

**CLI Framework:**
- **Recommendation:** Cobra (Go standard for CLI tools)
- **Alternatives:** urfave/cli, stdlib flag package
- **Rationale:** Cobra provides command structure, flag parsing, shell completion, and help generation out-of-box

**Terminal Output:**
- **Colored output:** `fatih/color` or `gookit/color`
- **Progress indicators:** `schollz/progressbar` or `cheggaaa/pb`
- **Tables:** `olekukonko/tablewriter`

**Configuration:**
- **Config parsing:** `spf13/viper` (pairs with Cobra)
- **Env var support:** Built into Viper
- **Config validation:** Schema validation on load

**Shell Completion:**
- Built into Cobra (`cobra-cli` generates completion commands)
- Auto-generates for bash/zsh/fish
- Dynamic completion for problem names (reads from database)

**Cross-Platform:**
- ANSI color support: Windows 10+, macOS, Linux
- Path handling: `filepath` package for cross-platform paths
- Editor detection: Check $EDITOR, fallback to platform defaults

**Quality Bars:**
- Help text: Every command has `--help`
- Error messages: Clear, actionable (not cryptic)
- Command naming: Consistent verb-noun structure
- Flag naming: Follow GNU/POSIX conventions (e.g., `-v` verbose, `-q` quiet)

## Project Scoping & Phased Development

### MVP Strategy & Philosophy

**MVP Approach:** Experience MVP

The MVP strategy prioritizes delivering the core test-driven workflow experience that makes developers want to practice DSA in the terminal. Phase 1 validates that the foundation (CLI + Go tooling + test workflow) creates genuine value. Phase 2 immediately layers in the motivation system (streaks, celebrations) to achieve the "30 mins daily" success metric.

**Strategic Rationale:**
- **Phase 1 proves:** "Does the CLI DSA practice workflow feel native to Go developers?"
- **Phase 2 delivers:** "Does the motivation system make it habit-forming?"
- **Phase 3+ enhances:** Learning science and expansion after core value validated

This staged approach de-risks the most critical assumption: that developers will adopt a CLI-based DSA tool if the experience is excellent, before investing in advanced features.

**Resource Requirements:**
- **MVP (Phase 1):** Solo developer or 2-person team (1 backend + 1 problem curator)
- **Skills needed:** Go expertise, CLI tool design, SQLite database, test framework integration
- **Timeline estimate:** 4-6 weeks for MVP assuming full-time focus
- **Phase 2:** Same team, additional 2-3 weeks to add tracking/celebration layer

### MVP Feature Set (Phase 1) - Detailed

**Core User Journey Supported:** Journey 1 (First-time user onboarding) and basic Journey 2 (daily practice without advanced tracking)

**Must-Have Capabilities (from Product Scope):**

1. **Test-Driven Workflow**
   - Pre-written tests for every problem
   - Developers write solutions to make tests pass
   - Native `go test` integration (zero wrapper friction)

2. **Scaffolding Generator**
   - `dsa solve [problem]` creates boilerplate automatically
   - Opens file in user's preferred editor
   - Includes problem description, test cases, solution template

3. **Go Tooling Integration**
   - Native `go test` execution
   - Go module structure
   - Idiomatic Go patterns in generated code

4. **Core CLI Commands**
   - `dsa init` - Initialize workspace
   - `dsa solve [problem]` - Start problem
   - `dsa test` - Run tests
   - `dsa status` - Basic progress (problems solved, difficulty breakdown)

5. **Local Storage**
   - SQLite database for progress tracking
   - Problem metadata (difficulty, topic, tags)
   - Solution history and timestamps
   - Completely offline, no network dependency

6. **Problem Library (10-15 curated problems)**
   - Arrays (3-4 problems): two-sum, three-sum, merge-intervals
   - Linked Lists (2-3 problems): reverse-list, detect-cycle
   - Trees (3-4 problems): traversals, BST operations
   - Sorting (2-3 problems): quicksort, mergesort
   - Mix of Easy/Medium/Hard difficulties

**What's NOT in MVP (Explicitly Out of Scope for Phase 1):**
- Advanced tracking (streaks, weak area identification)
- Celebration system (ASCII art, encouraging messages)
- Spaced repetition scheduling
- Multi-language support (Go only for MVP)
- Import/export functionality
- Shell completion (moved to Phase 2 based on CLI discussion)

**MVP Success Criteria:**
- Creator can practice 30 minutes without friction
- Setup takes <5 minutes from clone to first problem
- Test workflow feels native to Go developers
- All 6 core capabilities function reliably

### Post-MVP Features (Phased Roadmap)

**Phase 2: Excitement Layer (Ship immediately after MVP validation)**

**Goal:** Transform practice from "nice to have" into daily habit through motivation and tracking.

**5 Priority Features:**
1. Enhanced `dsa status` dashboard (streak counter, topic coverage, recent activity)
2. Streak tracking with visual feedback (🔥 7-day streak!)
3. Celebration moments (ASCII art, programming humor, milestone recognition)
4. Easy/Medium/Hard tracking with progress visualization
5. Shell completion (bash/zsh/fish) for professional CLI UX

**Phase 2 Success Criteria:**
- Streaks drive daily return behavior
- Celebrations feel authentic and motivating
- 30/90-day retention improves vs Phase 1

**Phase 3: Smart Learning (After habit formation validated)**

**Goal:** Apply learning science principles for more effective practice.

**5 Features:**
1. Spaced repetition scheduling (intelligent review suggestions)
2. Weak area identification (analytics highlighting topics needing practice)
3. Personalized review suggestions ("Review these 3 from 2 weeks ago")
4. Performance analytics (time spent, success rates, improvement trends)
5. Freshness indicators ("Last solved: 14 days ago")

**Phase 4: Vision (Polish & Expansion)**

**Goal:** Expand reach and add features that enhance but aren't essential.

**5 Future Features:**
1. Import/export functionality (JSON/CSV, multi-machine sync)
2. Multi-language support (Rust, Python, TypeScript with native tooling)
3. Advanced Go challenges (concurrency katas, memory optimization, benchmarking)
4. Community problem contributions (enable submit/curate workflow)
5. Custom problem sets (themed collections)

### Risk Mitigation Strategy

**Technical Risks:**

**Risk 1: Native `go test` integration complexity**
- **Mitigation:** Prototype test integration in week 1, validate it feels native
- **Fallback:** If native integration too complex, ship with `dsa test` wrapper initially

**Risk 2: SQLite ORM performance with local data**
- **Mitigation:** Design schema to avoid blocking operations (<100ms queries)
- **Fallback:** Choose appropriate Go ORM (GORM, sqlx) or use stdlib database/sql for maximum performance

**Risk 3: Cross-platform CLI compatibility (Windows/macOS/Linux)**
- **Mitigation:** Use Go's cross-compilation, test on all platforms early
- **Fallback:** Ship macOS/Linux first, Windows as Phase 1.5

**Market Risks:**

**Risk 1: Developers prefer browser-based tools (Leetcode habit)**
- **Validation approach:** Launch with 10-20 beta testers who are terminal-native developers
- **Learning goal:** Do terminal developers adopt CLI DSA practice if UX is excellent?
- **Mitigation:** MVP focuses on terminal-first developers, validates niche before expanding

**Risk 2: Kata-machine satisfies the market (no need for tracking/motivation)**
- **Validation approach:** Phase 1 → Phase 2 comparison shows retention improvement
- **Learning goal:** Does motivation layer drive sustained engagement vs pure CLI practice?
- **Mitigation:** If tracking doesn't help, focus on problem quality and Go-specific features

**Risk 3: Open source competition (someone clones concept)**
- **Mitigation:** Execution quality matters more than idea. Ship fast, iterate based on feedback.
- **Strategy:** Community becomes moat—contributors, problem library, multi-language support

**Resource Risks:**

**Risk 1: Solo developer timeline pressure**
- **Contingency:** Cut problem library to 8-10 instead of 10-15 for MVP
- **Contingency:** Ship without shell completion initially (add in Phase 2)
- **Contingency:** Start with Easy/Medium only, add Hard problems post-launch

**Risk 2: Problem curation takes longer than expected**
- **Mitigation:** Start with well-known problems (two-sum, reverse-list) that have clear test cases
- **Fallback:** Ship 8 problems instead of 15, add more weekly post-launch

**Risk 3: Community adoption slower than expected**
- **Strategy:** Focus on personal dogfooding test first (creator uses 30 mins daily)
- **Strategy:** Share in Go community channels (Reddit r/golang, Gopher Slack) after MVP polish
- **Acceptance:** Success is creator + 10-20 active users, not thousands initially

### Scope Validation Checklist

Before proceeding to implementation, validate:

- ✅ **MVP delivers on core value:** Test-driven DSA practice feels native to terminal
- ✅ **Success metric achievable:** Phase 1 + Phase 2 enables "30 mins daily" goal
- ✅ **Resource-realistic:** Solo dev or 2-person team can ship MVP in 4-6 weeks
- ✅ **Risk-aware:** Technical, market, and resource risks identified with mitigation
- ✅ **Phases sequenced correctly:** Foundation → Motivation → Learning Science → Expansion

## Functional Requirements

### Problem Management

- **FR1:** Developer can initialize a DSA practice workspace in their local environment
- **FR2:** Developer can browse available problems by topic (arrays, linked lists, trees, sorting)
- **FR3:** Developer can browse available problems by difficulty (easy, medium, hard)
- **FR4:** Developer can start working on a specific problem by name or identifier
- **FR5:** System provides problem description, constraints, and examples when developer starts a problem
- **FR6:** System provides pre-written test cases for every problem
- **FR7:** Developer can view problem metadata (difficulty, topic, tags, description)

### Test-Driven Workflow

- **FR8:** System generates scaffolded solution file with boilerplate code when developer starts a problem
- **FR9:** Developer can run tests against their solution using native Go testing tools
- **FR10:** System opens problem files in developer's preferred code editor automatically
- **FR11:** System validates solution correctness using pre-written test cases
- **FR12:** Developer receives clear test results showing passed/failed test cases
- **FR13:** System tracks test execution history for each problem attempt

### Progress Tracking & Analytics

- **FR14:** Developer can view their overall progress (total problems solved, by difficulty, by topic)
- **FR15:** System tracks completion status for each problem (not started, in progress, completed)
- **FR16:** Developer can view their solution history for previously attempted problems
- **FR17:** System records timestamps for problem attempts and completions
- **FR18:** System tracks daily practice streaks (consecutive days with at least one problem solved) [Phase 2]
- **FR19:** System identifies weak areas based on problem-solving history and patterns [Phase 3]
- **FR20:** Developer can view personalized problem recommendations based on practice history [Phase 3]
- **FR21:** System tracks time spent on each problem for performance analytics [Phase 3]

### CLI Configuration & Customization

- **FR22:** Developer can configure default code editor via configuration file
- **FR23:** Developer can configure default code editor via environment variables
- **FR24:** Developer can override configuration settings via command-line flags
- **FR25:** Developer can filter problems by difficulty preference in configuration
- **FR26:** Developer can set output verbosity preferences (quiet, normal, verbose)
- **FR27:** Developer can enable or disable celebration features in configuration
- **FR28:** Developer can enable or disable programming humor in output messages
- **FR29:** System respects standard environment conventions (NO_COLOR, TTY detection)

### Output & Reporting

- **FR30:** System provides human-friendly terminal output with colors and formatting by default
- **FR31:** Developer can request machine-parseable output (JSON format) for scripting
- **FR32:** System displays progress visualizations (progress bars, ASCII charts) for status commands
- **FR33:** System provides clear error messages when commands fail
- **FR34:** System outputs structured data to stdout and errors to stderr following UNIX conventions
- **FR35:** Developer can use shell completion for commands and problem names [Phase 2]

### Motivation & Celebration System [Phase 2]

- **FR36:** System celebrates milestone achievements with ASCII art and encouraging messages [Phase 2]
- **FR37:** System displays streak counter with visual indicators in status output [Phase 2]
- **FR38:** System provides encouraging messages with programming humor when problems are solved [Phase 2]
- **FR39:** System recognizes specific milestones (first problem, first hard problem, streak milestones) [Phase 2]
- **FR40:** System displays topic completion progress and celebrates topic mastery [Phase 2]

### Data Management

- **FR41:** System stores all progress data locally using embedded database (no network dependency)
- **FR42:** Developer can export their progress data in portable formats (JSON/CSV) [Phase 4]
- **FR43:** Developer can import previously exported progress data [Phase 4]
- **FR44:** System maintains data integrity across sessions with zero progress loss
- **FR45:** System supports both global and project-specific configuration

### Scripting & Automation

- **FR46:** Commands can be composed and chained using shell operators (pipes, redirects, logical operators)
- **FR47:** System provides consistent exit codes for success and different error types
- **FR48:** Commands work in non-interactive mode (CI/CD, scripts) without requiring user input
- **FR49:** System detects TTY vs non-TTY contexts and adjusts output accordingly
- **FR50:** Developer can force decisions via flags to skip confirmations in automated contexts

### Community & Extensibility [Phase 4]

- **FR51:** Community contributor can submit new problem definitions following standardized format [Phase 4]
- **FR52:** Community contributor can add support for additional programming languages [Phase 4]
- **FR53:** System provides clear documentation for contribution workflows [Phase 4]
- **FR54:** Developer can create custom problem sets and themed collections [Phase 4]

## Non-Functional Requirements

### Performance

**NFR1:** CLI commands execute with cold start time <500ms
**NFR2:** Warm command execution (after first run) completes in <100ms
**NFR3:** Database query operations complete without blocking user interaction (<100ms)
**NFR4:** File scaffolding and generation operations complete in <200ms
**NFR5:** Status dashboard rendering completes in <300ms regardless of solution history size
**NFR6:** Problem library browsing operations return results in <100ms
**NFR7:** Test execution performance matches native `go test` performance (zero overhead)

### Reliability & Data Integrity

**NFR8:** System maintains 100% data integrity across sessions with zero progress loss
**NFR9:** System operates completely offline with no network dependencies for core functionality
**NFR10:** Database operations use transactions to ensure atomic updates
**NFR11:** System gracefully handles interrupted operations (crash, kill) without data corruption
**NFR12:** Configuration file parsing fails safely with clear error messages on invalid syntax
**NFR13:** System provides automatic recovery mechanisms for corrupted local database

### Integration & Compatibility

**NFR14:** System integrates with native Go testing framework without wrapper overhead
**NFR15:** Generated code follows idiomatic Go patterns and conventions
**NFR16:** System respects user's `$EDITOR` environment variable and standard editor detection conventions
**NFR17:** CLI output respects `NO_COLOR` environment variable and TTY detection for appropriate formatting
**NFR18:** System follows UNIX conventions for exit codes, stdin/stdout/stderr, and signal handling
**NFR19:** Shell completion integrates seamlessly with bash, zsh, and fish completion systems

### Portability & Cross-Platform Support

**NFR20:** System runs without modification on macOS, Linux, and Windows operating systems
**NFR21:** System handles platform-specific path conventions correctly across all platforms
**NFR22:** ANSI color output renders correctly on Windows 10+, macOS, and Linux terminals
**NFR23:** Binary distribution size remains <20MB for single-binary CLI distribution
**NFR24:** System requires no external dependencies beyond Go runtime for core functionality

### Maintainability & Extensibility

**NFR25:** Codebase architecture supports addition of new programming languages without core refactoring
**NFR26:** Problem definition format allows community contributions without code changes
**NFR27:** Configuration schema supports backward compatibility when adding new settings
**NFR28:** CLI command structure supports addition of new commands without breaking existing workflows
**NFR29:** Database schema supports migrations for future feature additions
**NFR30:** Code follows Go best practices and passes standard linters (golint, go vet, staticcheck)

### Usability & Developer Experience

**NFR31:** Error messages provide actionable guidance rather than cryptic technical errors
**NFR32:** Help text (`--help`) for each command is comprehensive and includes usage examples
**NFR33:** Setup process from clone to first problem takes <5 minutes
**NFR34:** Configuration defaults work for 80% of users without requiring customization
**NFR35:** Command naming follows intuitive verb-noun conventions recognizable to CLI users
**NFR36:** Output formatting provides clear visual hierarchy with appropriate use of color and spacing
