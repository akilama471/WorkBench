# CONTRIBUTING.md

# Contributing to WorkBench

Thank you for your interest in contributing to **WorkBench**.

WorkBench is a free and open-source, native, lightweight, cross-platform local development environment manager built with Go.

The goal of WorkBench is to provide developers with a fast, simple, and native alternative to tools such as Laragon, XAMPP, and WAMP.

Every contribution is welcome, including:

* Bug fixes
* Documentation improvements
* Tests
* Performance improvements
* Cross-platform improvements
* Service integrations
* Runtime integrations
* UI improvements

Before contributing, please read:

* `AGENTS.md`
* `ROADMAP.md`
* `ARCHITECTURE.md`

These files define the project's engineering rules and architecture.

---

# 1. Project Philosophy

WorkBench follows these principles:

> **Native. Fast. Simple. Modular. Developer-focused.**

Contributions should preserve these principles.

WorkBench is intentionally not a web application.

Do not introduce:

* Electron
* React
* Vue
* WebView-based UI
* Browser-based desktop UI

The desktop application must remain native.

---

# 2. Before You Start

Before creating a pull request:

1. Read `AGENTS.md`.
2. Read `ROADMAP.md`.
3. Read `ARCHITECTURE.md`.
4. Search existing issues and pull requests.
5. Avoid duplicating existing work.
6. Confirm that your change fits the project scope.

If you are unsure whether an idea fits WorkBench, open a discussion or issue before implementing a large change.

---

# 3. Development Environment

## Required Tools

You should have:

* Go
* Git

The project is primarily developed on Windows during Phase 1.

Future development will include:

* Linux
* macOS

---

# 4. Repository Structure

The project follows this structure:

```text id="w55wy5"
workbench/
│
├── cmd/
├── internal/
├── pkg/
├── ui/
├── resources/
├── tests/
│
├── AGENTS.md
├── ROADMAP.md
├── ARCHITECTURE.md
└── CONTRIBUTING.md
```

Before adding a new package, determine which architectural layer it belongs to.

Do not create random top-level directories.

---

# 5. Branching Strategy

Use the following branch naming convention.

## Feature

```text id="1hl1xz"
feature/<short-description>
```

Example:

```text id="1hl1xz"
feature/php-version-manager
```

## Bug Fix

```text id="10ny7c"
fix/<short-description>
```

Example:

```text id="4d14v2"
fix/apache-status-detection
```

## Performance

```text id="zjv8jh"
perf/<short-description>
```

Example:

```text id="30ah5f"
perf/reduce-startup-time
```

## Documentation

```text id="b2g9jw"
docs/<short-description>
```

Example:

```text id="k1cy41"
docs/package-manifest
```

## Refactoring

```text id="4c6cx4"
refactor/<short-description>
```

Example:

```text id="h0r2q1"
refactor/process-manager
```

---

# 6. Branch Rules

Do not commit directly to the main branch.

Use a feature or fix branch.

Before creating a pull request:

```bash id="9u7h2c"
git pull
```

or rebase your branch according to project maintainer instructions.

Keep branches focused.

Do not combine:

```text id="6o0e7c"
PHP feature
+
UI redesign
+
README rewrite
```

in a single unrelated pull request.

---

# 7. Commit Messages

Use clear and concise commit messages.

Recommended format:

```text id="r7b1hx"
type: short description
```

Types:

```text id="r20w8j"
feat
fix
refactor
perf
test
docs
build
chore
```

Examples:

```text id="u11e0c"
feat: add PHP version discovery
```

```text id="efxjvq"
fix: handle Apache startup failure
```

```text id="lkw0q4"
test: add PHP switching rollback tests
```

```text id="70y6jv"
docs: document package manifests
```

---

# 8. Commit Guidelines

A commit should represent one logical change.

Good:

```text id="35japq"
feat: add PHP version discovery
```

Bad:

```text id="9z8u4y"
update everything
```

Avoid huge commits whenever possible.

Prefer:

```text id="49m3k0"
feat: add service abstraction
feat: add Apache service
test: add Apache service tests
```

over one massive commit.

---

# 9. Go Coding Standards

Use idiomatic Go.

Before submitting code, run:

```bash id="qyk2p5"
gofmt
```

Run:

```bash id="qyjkl7"
go test ./...
```

Run:

```bash id="k1u5bp"
go vet ./...
```

Use clear names.

Good:

```go id="0ap9mb"
ProcessManager
```

Bad:

```go id="6df0m4"
ProcMgr
```

Avoid unnecessary abbreviations.

---

# 10. Architecture Rules

The architecture defined in `ARCHITECTURE.md` is mandatory.

The dependency direction is:

```text id="9nyxy6"
GUI / CLI
    ↓
Application
    ↓
Core
    ↓
Domain
    ↓
Infrastructure
```

Do not bypass layers.

---

## 10.1 GUI Rules

The GUI must not:

* Start processes directly.
* Stop processes directly.
* Read configuration files directly.
* Write configuration files directly.
* Access SQLite directly.

Correct:

```text id="sl46ty"
GUI
 ↓
Application
 ↓
Core
```

---

## 10.2 Service Rules

Services must use the Service abstraction.

Do not implement Apache-specific logic inside the generic Service Manager.

Bad:

```go id="r2z1yz"
if service == "apache" {
    ...
}
```

Good:

```text id="n2xwtr"
Service Manager
        ↓
ApacheService
```

---

## 10.3 Runtime Rules

PHP is a Runtime.

Do not model PHP as a Service.

Future runtimes must be implemented through the Runtime architecture.

---

# 11. Adding a New Service

When adding a new service such as Redis, follow this pattern.

## Step 1 — Create the Service Package

Example:

```text id="7w0t5m"
internal/service/redis/
```

---

## Step 2 — Implement the Service Interface

Implement:

```text id="5rj9u7"
ID
Name
Start
Stop
Restart
Status
IsInstalled
```

---

## Step 3 — Use Existing Infrastructure

Use:

```text id="r7f2h5"
Process Manager
Filesystem
Configuration Manager
Logger
```

Do not create a second process manager.

---

## Step 4 — Register the Service

Register Redis through the application/service composition layer.

---

## Step 5 — Add Tests

Test:

* Installation detection.
* Start behavior.
* Stop behavior.
* Restart behavior.
* Status behavior.

---

## Step 6 — Update Documentation

Update the relevant documentation.

Do not add Redis to Phase 1 unless the roadmap is intentionally updated.

---

# 12. Adding a New Runtime

When adding a runtime such as Python:

```text id="c3s6m0"
internal/runtime/python/
```

The runtime should follow the Runtime abstraction.

The runtime must support:

* Installed version discovery.
* Active version detection where applicable.
* Version selection.

Do not duplicate PHP-specific code blindly.

Extract shared logic only when the abstraction is genuinely common.

---

# 13. Adding a New Tool

Tools such as:

* HeidiSQL
* ngrok
* Composer

should not be implemented as Services unless they actually run as long-running services.

Use the Tool domain when it is introduced.

Do not force every executable into the Service architecture.

---

# 14. Adding Platform-Specific Code

Platform-specific code must be isolated.

Use:

```text id="e6hx72"
*_windows.go
*_linux.go
*_darwin.go
```

Examples:

```text id="r0q6xq"
process_windows.go
process_linux.go
process_darwin.go
```

Do not write:

```go id="l16y1e"
if runtime.GOOS == "windows" {
    ...
}
```

throughout the entire codebase.

Use platform boundaries.

---

# 15. Testing Requirements

New core behavior should include tests.

Test important behavior such as:

* Path resolution.
* Directory initialization.
* Version discovery.
* Version switching.
* Rollback.
* Service status.
* Package validation.

Use:

```go id="f4ih9u"
t.TempDir()
```

for filesystem tests.

Do not use personal machine paths.

Bad:

```text id="oqpdyv"
C:\Users\John\Desktop\WorkBench
```

---

# 16. External Dependencies

Before adding a dependency, ask:

1. Can the Go standard library solve this?
2. Is the dependency actively maintained?
3. Is the license compatible?
4. Does it support Windows, Linux, and macOS?
5. Does it add unnecessary complexity?

Do not add a dependency for trivial functionality.

---

# 17. Security Requirements

Security issues should be treated seriously.

Do not:

* Execute user input through a shell.
* Concatenate shell commands.
* Trust package filenames.
* Extract archives without path validation.
* Ignore package checksums.
* Execute arbitrary downloaded binaries without validation.

If you discover a security vulnerability, do not immediately publish detailed exploit instructions in a public issue.

Contact the project maintainers through the project's security reporting process.

---

# 18. Pull Request Requirements

A pull request should include:

1. A clear title.
2. A short explanation of the problem.
3. A summary of the solution.
4. Testing performed.
5. Any platform limitations.

Example:

```markdown
## Problem

Apache status was reported as running after the process exited immediately.

## Solution

Added process state verification after startup.

## Testing

- go test ./...
- go vet ./...

## Platform

Windows tested.
Linux/macOS implementation unchanged.
```

---

# 19. Pull Request Checklist

Before opening a pull request:

* [ ] I read `AGENTS.md`.
* [ ] I read `ROADMAP.md`.
* [ ] I read `ARCHITECTURE.md`.
* [ ] My branch has a clear name.
* [ ] My changes are focused.
* [ ] I did not introduce Electron.
* [ ] I did not introduce React.
* [ ] I did not introduce WebView-based UI.
* [ ] I followed the architecture.
* [ ] I added tests where appropriate.
* [ ] `go test ./...` passes.
* [ ] `go vet ./...` passes.
* [ ] Go code is formatted.
* [ ] I did not hardcode personal paths.
* [ ] I did not delete user data.
* [ ] I updated documentation where necessary.

---

# 20. Documentation Contributions

Documentation is important.

Documentation should be:

* Clear.
* Technically accurate.
* Easy for new contributors to understand.

When changing architecture, update:

```text id="yd7s5x"
ARCHITECTURE.md
```

When changing the implementation plan, update:

```text id="q0k2bp"
ROADMAP.md
```

When changing contributor rules, update:

```text id="j8e2p4"
CONTRIBUTING.md
```

---

# 21. Issue Guidelines

Before opening an issue, search existing issues.

A good bug report should include:

* WorkBench version or commit.
* Operating system.
* Architecture.
* Steps to reproduce.
* Expected behavior.
* Actual behavior.
* Relevant logs.

Example:

```markdown
## Environment

OS: Windows 11
Architecture: amd64

## Steps

1. Install PHP 8.3.
2. Install PHP 8.2.
3. Switch to PHP 8.2.

## Expected

PHP 8.2 becomes active.

## Actual

PHP 8.3 remains active.

## Logs

...
```

Do not include sensitive personal information.

---

# 22. Feature Requests

Feature requests should explain:

1. The problem.
2. The proposed solution.
3. Why it fits WorkBench.
4. Possible architectural impact.

Do not only write:

```text
"Add Redis"
```

Explain the developer problem the feature solves.

---

# 23. Code Review Philosophy

Code review focuses on:

* Correctness.
* Architecture.
* Maintainability.
* Security.
* Cross-platform compatibility.
* Performance.

Code review is not about personal coding style preferences.

Contributors are expected to respond constructively to feedback.

---

# 24. Maintainer Principles

Maintainers should prioritize:

1. Stability.
2. Simplicity.
3. Security.
4. Cross-platform compatibility.
5. Developer experience.

Do not merge features only because they are technically interesting.

WorkBench must remain focused.

---

# 25. Project Direction

WorkBench is built around this principle:

> **The core engine is the product. The GUI and CLI are interfaces to the core engine.**

All contributions should strengthen this architecture.

The long-term goal is to make WorkBench:

```text id="m5ssg0"
Fast
Native
Open Source
Cross-Platform
Modular
Developer-Friendly
```

Thank you for contributing to WorkBench.
