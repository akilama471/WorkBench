# AGENTS.md

# DevBox — AI Coding Agent Instructions

You are an autonomous senior software architect and Go engineer working on **DevBox**.

DevBox is a **free and open-source, native, lightweight, cross-platform local development environment manager** designed as an alternative to Laragon, XAMPP, and WAMP.

Your job is to build this project incrementally, safely, and with production-quality engineering practices.

---

# 1. Project Vision

DevBox is not merely an Apache/PHP/MariaDB bundle.

DevBox is a **local development environment orchestration platform**.

The long-term goal is to provide developers with a simple, fast, lightweight native application for managing local development runtimes, services, and tools.

The product should eventually support:

* Apache
* PHP
* MariaDB
* Nginx
* Python
* Redis
* Mailpit
* HeidiSQL
* ngrok
* Node.js
* Composer
* Other development tools

The initial implementation must remain focused.

---

# 2. Non-Negotiable Requirements

These requirements are mandatory.

## 2.1 Native GUI

The GUI MUST NOT use:

* Electron
* React
* Vue
* Angular
* Svelte
* HTML/CSS-based UI
* WebView-based UI
* Browser-based UI
* Any web technology for the desktop GUI

Do not introduce a web frontend.

The application is a **native desktop application**.

The preferred GUI technology is:

> **Fyne**

Use Go-native GUI architecture.

If the current repository already contains a GUI technology decision, preserve the existing decision unless the user explicitly changes it.

---

## 2.2 Programming Language

The primary programming language is:

> **Go**

Do not rewrite the project in Rust, C++, Java, C#, or another language.

Use Go idiomatic conventions.

---

## 2.3 Open Source

The project is free and open source.

Do not introduce:

* Proprietary cloud dependencies
* Paid SaaS dependencies
* Vendor lock-in
* Required online accounts
* Mandatory telemetry
* Mandatory cloud infrastructure

The application must work locally and offline after required packages are installed.

---

## 2.4 Cross-Platform

The project is Windows-first but must be designed for cross-platform support.

Target platforms:

1. Windows
2. Linux
3. macOS

Do not implement Windows-specific logic directly inside generic packages.

Use platform-specific files where necessary:

```text
*_windows.go
*_linux.go
*_darwin.go
```

Use Go build tags when appropriate.

---

## 2.5 Lightweight and Fast

Performance is a product requirement.

Avoid:

* Unnecessary background polling
* Busy loops
* Excessive goroutines
* Heavy dependency trees
* Repeated filesystem scans
* Repeated process spawning
* Unnecessary database queries
* Memory-heavy abstractions

The application should start quickly.

The GUI must remain responsive.

---

# 3. Product Terminology

Use the following terminology consistently.

## Service

A long-running development server or daemon.

Examples:

* Apache
* MariaDB
* Nginx
* Redis
* Mailpit

## Runtime

A programming language runtime or execution environment.

Examples:

* PHP
* Python
* Node.js

## Tool

A developer utility.

Examples:

* HeidiSQL
* ngrok
* Composer
* Git

Do not incorrectly model every component as a "service".

---

# 4. Initial Product Scope

The first milestone is the **Phase 1 MVP**.

Phase 1 MUST support:

* Application installation/root detection
* DevBox directory creation
* Apache integration
* PHP integration
* MariaDB integration
* Start Apache
* Stop Apache
* Restart Apache
* Start MariaDB
* Stop MariaDB
* Restart MariaDB
* Service status detection
* PHP version listing
* PHP version switching
* Active PHP runtime management
* Basic log access
* CLI
* Native GUI
* SQLite persistence
* Basic package metadata handling

Do NOT implement the following in Phase 1 unless explicitly requested:

* Nginx
* Redis
* Python
* Mailpit
* HeidiSQL
* ngrok
* Node.js
* Composer management
* Virtual hosts
* Automatic SSL
* Laravel project creation
* Cloud synchronization
* User accounts
* Telemetry
* Package marketplace

Do not expand the scope without an explicit requirement.

---

# 5. Core Architecture

Use the following high-level architecture:

```text
┌──────────────────────────────────────────────┐
│                  DevBox GUI                  │
│                   Fyne                       │
└──────────────────────┬───────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────┐
│              DevBox Core Engine              │
│                    Go                        │
│                                              │
│  Service Manager                             │
│  Runtime Manager                             │
│  Package Manager                             │
│  Configuration Manager                       │
│  Process Manager                              │
│  Environment Manager                         │
│  Project Manager                              │
└──────┬─────────┬─────────┬───────────────────┘
       │         │         │
       ▼         ▼         ▼
    Apache      PHP      MariaDB
```

## Critical Architecture Rule

The GUI MUST NOT directly:

* Start processes
* Stop processes
* Modify Apache configuration
* Modify PHP configuration
* Modify MariaDB configuration
* Scan installed packages
* Directly write application configuration

The GUI must call the application/core layer.

Correct:

```text
GUI
  ↓
Application Layer
  ↓
Core Manager
  ↓
Service/Runtime Manager
  ↓
Process/File System
```

Incorrect:

```text
GUI
  ↓
os/exec
```

---

# 6. Recommended Repository Structure

Use this structure unless the repository already contains a well-designed equivalent.

```text
devbox/
│
├── cmd/
│   ├── devbox/
│   │   └── main.go
│   │
│   └── devbox-cli/
│       └── main.go
│
├── internal/
│   │
│   ├── app/
│   │   ├── application.go
│   │   └── container.go
│   │
│   ├── core/
│   │   ├── service_manager.go
│   │   ├── runtime_manager.go
│   │   └── environment_manager.go
│   │
│   ├── service/
│   │   ├── service.go
│   │   ├── manager.go
│   │   │
│   │   ├── apache/
│   │   │   └── apache.go
│   │   │
│   │   └── mariadb/
│   │       └── mariadb.go
│   │
│   ├── runtime/
│   │   ├── runtime.go
│   │   ├── manager.go
│   │   │
│   │   └── php/
│   │       ├── php.go
│   │       ├── manager.go
│   │       └── version.go
│   │
│   ├── package/
│   │   ├── package.go
│   │   ├── manager.go
│   │   ├── manifest.go
│   │   └── repository.go
│   │
│   ├── process/
│   │   ├── process.go
│   │   ├── manager.go
│   │   ├── process_windows.go
│   │   ├── process_linux.go
│   │   └── process_darwin.go
│   │
│   ├── config/
│   │   ├── config.go
│   │   ├── loader.go
│   │   └── writer.go
│   │
│   ├── filesystem/
│   │   ├── paths.go
│   │   ├── directories.go
│   │   └── archive.go
│   │
│   ├── project/
│   │   ├── project.go
│   │   ├── manager.go
│   │   └── detector.go
│   │
│   ├── database/
│   │   ├── database.go
│   │   └── migrations.go
│   │
│   └── logger/
│       └── logger.go
│
├── pkg/
│   ├── models/
│   │   ├── service.go
│   │   ├── runtime.go
│   │   ├── package.go
│   │   └── project.go
│   │
│   └── events/
│       └── events.go
│
├── ui/
│   ├── app.go
│   ├── window.go
│   │
│   ├── pages/
│   │   ├── dashboard.go
│   │   ├── services.go
│   │   ├── runtimes.go
│   │   ├── projects.go
│   │   └── settings.go
│   │
│   ├── components/
│   │   ├── service_card.go
│   │   ├── runtime_selector.go
│   │   └── status_badge.go
│   │
│   └── state/
│       └── app_state.go
│
├── resources/
│   ├── icons/
│   └── manifests/
│
├── scripts/
│
├── tests/
│
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── AGENTS.md
```

Do not create random top-level directories.

---

# 7. Application Root Directory

The default DevBox environment root should contain:

```text
DevBox/
│
├── bin/
│   ├── apache/
│   │   └── <version>/
│   │
│   ├── php/
│   │   ├── <version>/
│   │   └── ...
│   │
│   └── mariadb/
│       └── <version>/
│
├── active/
│   └── php
│
├── www/
│
├── data/
│   └── mariadb/
│
├── etc/
│   ├── apache/
│   ├── php/
│   └── mariadb/
│
├── logs/
│   ├── apache/
│   ├── php/
│   └── mariadb/
│
├── backup/
│
├── packages/
│
├── cache/
│
└── devbox.db
```

## Directory Rules

The directory manager must be idempotent.

Running initialization multiple times must not destroy existing data.

Never delete:

```text
www/
data/
backup/
```

without an explicit destructive user action.

---

# 8. Service Architecture

Create a generic service abstraction.

Recommended interface:

```go
type Service interface {
    ID() string
    Name() string

    Start() error
    Stop() error
    Restart() error

    Status() Status
    IsInstalled() bool
}
```

Status values:

```go
type Status int

const (
    StatusUnknown Status = iota
    StatusStopped
    StatusStarting
    StatusRunning
    StatusStopping
    StatusError
)
```

Use typed errors where useful.

Example:

```go
var (
    ErrServiceNotInstalled = errors.New("service is not installed")
    ErrServiceAlreadyRunning = errors.New("service is already running")
)
```

---

# 9. Service Manager

The service manager owns registered services.

Conceptually:

```go
type Manager struct {
    services map[string]Service
}
```

Required operations:

```go
Start(id string) error
Stop(id string) error
Restart(id string) error
Status(id string) Status
All() []Service
```

The service manager must not know Apache implementation details.

Correct:

```text
Service Manager
    ↓
Service Interface
    ↓
ApacheService
```

Incorrect:

```go
if id == "apache" {
    exec.Command("httpd.exe")
}
```

---

# 10. Apache Service

Apache must be implemented as a service module.

Apache responsibilities:

* Resolve installed Apache version
* Resolve Apache executable
* Resolve Apache configuration
* Start Apache
* Stop Apache
* Restart Apache
* Detect running status
* Read Apache logs

Apache must not own generic process management.

Use the process manager.

Conceptual startup:

```text
ApacheService
    ↓
Resolve Apache installation
    ↓
Resolve httpd executable
    ↓
Resolve httpd.conf
    ↓
ProcessManager.Start()
```

Do not hardcode absolute paths.

---

# 11. MariaDB Service

MariaDB must be implemented as a service module.

MariaDB responsibilities:

* Resolve MariaDB installation
* Resolve server executable
* Resolve data directory
* Resolve configuration
* Start MariaDB
* Stop MariaDB
* Restart MariaDB
* Detect running status

The data directory must be separate from the MariaDB binary directory.

Correct:

```text
bin/mariadb/11.x/
data/mariadb/
```

Never store user database data inside:

```text
bin/mariadb/
```

---

# 12. Runtime Architecture

Runtime interface:

```go
type Runtime interface {
    ID() string
    Name() string

    InstalledVersions() []Version
    ActiveVersion() *Version
    Use(version string) error
}
```

PHP is the first runtime.

Do not model PHP as a service.

PHP is a runtime.

---

# 13. PHP Version Management

PHP versions must be stored independently.

Example:

```text
bin/php/
├── 8.1.29/
├── 8.2.20/
└── 8.3.30/
```

The active PHP version must be tracked separately.

Use:

```text
active/php
```

or an equivalent platform-safe abstraction.

Do not copy the entire PHP installation when switching versions.

Do not duplicate PHP binaries.

The PHP manager must:

1. Discover installed PHP versions
2. Validate a version
3. Get active PHP version
4. Switch active PHP version
5. Notify dependent components
6. Ensure Apache uses the selected PHP runtime

PHP switching must be atomic from the user's perspective.

If switching fails, the previous active version must remain active.

---

# 14. PHP Switching Transaction

Implement PHP switching safely.

Required conceptual flow:

```text
1. Validate target version
2. Validate PHP installation
3. Validate required files
4. Prepare new active state
5. Update runtime pointer/configuration
6. Update Apache integration
7. Verify configuration
8. Commit active version
9. Restart dependent service if required
```

If any step fails:

```text
Rollback
```

Do not leave the application in a partially switched state.

---

# 15. Process Manager

Create a platform-independent process abstraction.

Example:

```go
type StartConfig struct {
    Executable  string
    Args        []string
    Directory   string
    Environment []string
}
```

Interface:

```go
type Manager interface {
    Start(config StartConfig) (*Process, error)
    Stop(pid int) error
    IsRunning(pid int) bool
}
```

Platform-specific implementation is allowed.

Generic packages must not directly use Windows-only APIs.

Do not use shell strings such as:

```text
"start apache && ..."
```

Prefer direct process execution.

Avoid shell injection.

Arguments must be passed as argument arrays.

---

# 16. Configuration Manager

Configuration files must be managed through a dedicated configuration layer.

The configuration manager should support:

```go
Load(path string) ([]byte, error)
Save(path string, data []byte) error
Backup(path string) error
```

Before changing an important configuration file:

```text
Create backup
Apply change
Validate
```

Configuration updates must be safe.

Do not perform uncontrolled string replacement.

When possible, use structured configuration manipulation.

If a format is not structurally parseable, isolate text manipulation in a dedicated configuration adapter.

---

# 17. Package Manager

The package manager is a foundational feature.

It must be designed for:

* PHP
* Apache
* MariaDB
* Future runtimes
* Future services
* Future tools

Do not hardcode package URLs inside service implementations.

Use manifests.

Example:

```json
{
  "id": "php",
  "version": "8.3.30",
  "platform": "windows",
  "architecture": "x64",
  "download": {
    "url": "https://example.com/package.zip",
    "sha256": "..."
  },
  "archive": "zip",
  "install": {
    "directory": "bin/php/8.3.30"
  }
}
```

The package manager must:

1. Resolve package
2. Download package
3. Verify checksum
4. Extract package
5. Validate installation
6. Register package
7. Clean temporary files

Never install a package without checksum verification when a checksum is available.

Never extract directly into the final directory.

Use a temporary directory.

Recommended:

```text
cache/downloads/
cache/extract/
```

Then atomically move the validated package into place.

---

# 18. SQLite

Use SQLite for application metadata.

SQLite stores:

* Installed package metadata
* Active runtime selection
* Project metadata
* User settings
* Application state

SQLite is NOT the source of truth for the actual binaries.

The filesystem is the source of truth for installed software.

If SQLite says PHP exists but the directory is missing:

```text
Filesystem reality wins.
```

The application must be able to reconcile database metadata with the filesystem.

---

# 19. Database Migrations

Use explicit migrations.

Example:

```text
internal/database/
├── database.go
└── migrations.go
```

Migrations must be:

* Ordered
* Idempotent
* Versioned

Never silently drop user data during migrations.

---

# 20. Event System

Implement a lightweight internal event system.

Events may include:

```text
ServiceStarted
ServiceStopped
ServiceStatusChanged
ServiceError
PHPVersionChanged
PackageInstalled
PackageRemoved
```

The GUI should be able to react to events.

Avoid excessive polling.

If polling is necessary, use a reasonable interval and avoid aggressive loops.

---

# 21. CLI

The CLI is a first-class interface.

The GUI is not the core product.

The CLI must use the same core engine as the GUI.

Required initial commands:

```bash
devbox status
devbox start apache
devbox stop apache
devbox restart apache
devbox start mariadb
devbox stop mariadb
devbox php list
devbox php current
devbox php use 8.3.30
```

Expected:

```bash
devbox status
```

Example:

```text
Apache     Running
MariaDB    Running
PHP        8.3.30
```

The CLI must return non-zero exit codes on failure.

Do not print stack traces to normal users.

Provide useful human-readable errors.

---

# 22. GUI

Use Fyne.

The first GUI should be intentionally simple.

The dashboard should display:

```text
Apache       Running
MariaDB      Running
PHP          8.3.30
```

Actions:

```text
Start
Stop
Restart
Change PHP
```

The GUI must not become a complex dashboard.

Avoid:

* Excessive animations
* Decorative UI
* Large graphical widgets
* Web-style layouts
* Unnecessary navigation layers

The application is a developer utility.

Prioritize clarity and speed.

---

# 23. GUI State

The GUI must maintain application state separately from core implementation details.

The GUI should receive:

```text
Service Status
Runtime Status
Events
```

The GUI should call application methods.

Do not expose internal service implementations directly to Fyne components.

Bad:

```go
ui.ApacheService.Process.Cmd
```

Good:

```go
ui.AppState.Services["apache"].Status
```

---

# 24. Error Handling

Errors must be explicit.

Bad:

```go
return fmt.Errorf("failed")
```

Good:

```go
return fmt.Errorf("failed to start Apache: %w", err)
```

Add context at abstraction boundaries.

Do not swallow errors.

Never do:

```go
_ = os.Remove(path)
```

unless the error is intentionally handled and documented.

---

# 25. Logging

Use structured application logging.

Log categories:

```text
application
service
runtime
package
process
configuration
database
```

Logs should help diagnose:

* Apache startup failures
* MariaDB startup failures
* PHP switching failures
* Package installation failures

Do not log secrets.

Do not log full environment variables.

Do not log sensitive user data unnecessarily.

---

# 26. Testing Requirements

Every core module must be testable.

Prioritize tests for:

* Path resolution
* Directory initialization
* PHP version discovery
* PHP version switching
* Package manifest parsing
* Package checksum validation
* Service registration
* Service status transitions
* Process manager behavior
* Configuration backup behavior

Do not write tests that depend on a developer's personal machine path.

Avoid:

```text
C:\Users\Akila\...
```

Use temporary directories.

---

# 27. Test Strategy

Use:

```go
t.TempDir()
```

for filesystem tests.

Use interfaces for:

* Process execution
* File downloads
* Clock/time where necessary

Do not make tests depend on the real internet.

Package download tests must use mocked or local test sources.

---

# 28. Security Rules

The application executes external binaries.

Security is important.

Never:

* Execute user input through a shell
* Concatenate shell commands
* Trust package filenames
* Extract archives without path validation
* Allow archive path traversal
* Execute arbitrary downloaded files without validation

Prevent archive traversal.

Reject paths such as:

```text
../../evil.exe
```

Validate extracted paths before writing.

---

# 29. Dependency Rules

Before adding a dependency:

1. Check whether the Go standard library solves the problem.
2. Check dependency maintenance status.
3. Check license compatibility.
4. Check transitive dependency size.
5. Check cross-platform support.

Do not add dependencies for trivial functionality.

Prefer the standard library.

---

# 30. Coding Style

Follow idiomatic Go.

Use:

```text
gofmt
go vet
```

Prefer small interfaces.

Avoid premature abstractions.

Do not create interfaces only because "clean architecture" says so.

Interfaces should exist where they provide:

* Testability
* Platform abstraction
* Dependency inversion
* Stable boundaries

Use clear names.

Avoid abbreviations unless conventional.

Good:

```go
ProcessManager
ConfigurationManager
PackageRepository
```

Bad:

```go
ProcMgr
CfgMgr
PkgRepo
```

---

# 31. Architecture Rules

The following dependencies are preferred:

```text
UI
 ↓
Application
 ↓
Core
 ↓
Domain/Managers
 ↓
Infrastructure
```

Infrastructure includes:

* OS process APIs
* Filesystem
* SQLite
* HTTP downloads

Core logic must not depend directly on Fyne.

The core must be usable from the CLI without the GUI.

---

# 32. Implementation Workflow

When starting work on a task, follow this workflow.

## Step 1 — Inspect

Inspect:

* Repository structure
* Existing code
* `go.mod`
* README
* Existing tests
* Existing architecture

Do not assume the repository is empty.

## Step 2 — Plan

Before making significant changes:

1. Identify affected modules.
2. Identify architectural impact.
3. Identify potential platform issues.
4. Identify test requirements.

## Step 3 — Implement

Implement the smallest complete change.

Do not mix unrelated features.

## Step 4 — Test

Run relevant tests.

At minimum:

```bash
go test ./...
```

## Step 5 — Format

Run:

```bash
gofmt -w .
```

or format only affected Go files.

## Step 6 — Static Checks

Run:

```bash
go vet ./...
```

## Step 7 — Review

Check:

* Cross-platform compatibility
* Error handling
* Data safety
* Scope creep
* Architectural violations

---

# 33. Agent Behavior Rules

You are an autonomous coding agent.

You MUST:

* Read existing code before editing.
* Preserve working functionality.
* Prefer incremental changes.
* Keep changes focused.
* Explain architectural tradeoffs internally through code structure and comments where useful.
* Write tests for important core behavior.
* Run tests after meaningful changes.
* Fix failures caused by your changes.
* Avoid unnecessary rewrites.

You MUST NOT:

* Rewrite the entire project without a clear reason.
* Introduce Electron.
* Introduce React.
* Introduce WebView.
* Introduce web GUI technology.
* Add unrelated features.
* Add cloud services without explicit instruction.
* Add telemetry.
* Add user accounts.
* Delete user data.
* Remove `www/`, `data/`, or `backup/` automatically.
* Hardcode personal filesystem paths.
* Hardcode Windows-only behavior into generic code.

---

# 34. Phase 1 Implementation Order

Implement Phase 1 in this exact order unless a dependency requires otherwise.

## Phase 1.1 — Project Foundation

* Initialize Go module.
* Create repository structure.
* Create application entry points.
* Create logger.
* Create filesystem path abstraction.
* Create directory initializer.
* Add tests.

## Phase 1.2 — Process Manager

* Create process abstraction.
* Implement Windows process support.
* Design Linux/macOS boundaries.
* Add process tests where practical.

## Phase 1.3 — Service Abstraction

* Create Service interface.
* Create Service Manager.
* Create status model.
* Add service registration.

## Phase 1.4 — Apache

* Implement Apache service.
* Resolve Apache executable.
* Implement start/stop/restart.
* Implement status detection.
* Implement log access.

## Phase 1.5 — MariaDB

* Implement MariaDB service.
* Separate binary and data directories.
* Implement start/stop/restart.
* Implement status detection.

## Phase 1.6 — PHP Runtime

* Implement PHP runtime abstraction.
* Discover PHP versions.
* Track active PHP.
* Implement version switching.
* Add rollback behavior.

## Phase 1.7 — SQLite

* Add SQLite.
* Add migrations.
* Store installed package metadata.
* Store active runtime state.

## Phase 1.8 — CLI

* Add `status`.
* Add service commands.
* Add PHP commands.

## Phase 1.9 — Fyne GUI

* Add dashboard.
* Add service cards.
* Add status display.
* Add start/stop actions.
* Add PHP version selector.

## Phase 1.10 — Stabilization

* Improve error messages.
* Add tests.
* Run `go test ./...`.
* Run `go vet ./...`.
* Review Windows behavior.
* Update README.

---

# 35. Definition of Done

A feature is not complete until:

* The code compiles.
* Relevant tests pass.
* `go test ./...` passes where applicable.
* `go vet ./...` passes where applicable.
* The implementation follows the architecture.
* Errors are handled.
* No unrelated features were added.
* No user data is at risk.
* The feature works through the core layer.
* The GUI does not bypass the core.
* The CLI and GUI use the same underlying logic where applicable.

---

# 36. Final Product Philosophy

DevBox must be:

> **Fast. Native. Simple. Modular. Developer-focused.**

Do not turn DevBox into an enterprise management platform.

Do not turn DevBox into a cloud platform.

Do not turn DevBox into a web application.

Build a powerful local development tool with a small, understandable, maintainable codebase.

The core architectural principle is:

```text
Simple UI
    ↓
Strong Core
    ↓
Modular Services and Runtimes
```

Build the foundation correctly before adding features.
