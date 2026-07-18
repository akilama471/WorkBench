# ROADMAP.md

# DevBox — Phase 1 Implementation Roadmap

This document defines the implementation roadmap for **DevBox Phase 1**.

DevBox is a native, lightweight, cross-platform local development environment manager built with Go.

The Phase 1 objective is to build a stable foundation supporting:

* Apache
* PHP
* MariaDB
* Service management
* PHP version management
* CLI
* Native GUI

Read `AGENTS.md` before implementing any task.

---

# 1. Phase 1 Objective

At the end of Phase 1, a developer should be able to:

1. Install DevBox.
2. Launch DevBox.
3. Initialize the DevBox environment.
4. Install or register Apache, PHP, and MariaDB.
5. View installed PHP versions.
6. Switch between PHP versions.
7. Start Apache.
8. Stop Apache.
9. Restart Apache.
10. Start MariaDB.
11. Stop MariaDB.
12. Restart MariaDB.
13. View service status.
14. View basic logs.
15. Perform the same core operations through:

    * CLI
    * Native GUI

The following must work through the **same Go core engine**:

```text
CLI ───────┐
           ├── Core Engine
GUI ───────┘
```

---

# 2. Phase 1 Rules

## 2.1 Scope Control

Do not implement:

* Nginx
* Redis
* Python
* Mailpit
* HeidiSQL
* ngrok
* Node.js
* Composer manager
* Virtual hosts
* Automatic SSL
* Laravel project creation
* Cloud services
* User accounts
* Telemetry

These belong to future phases.

---

## 2.2 Implementation Order

Complete tasks in order unless a dependency requires a different order.

Do not skip foundational tasks.

---

# 3. Milestone 0 — Project Initialization

## Goal

Create the initial Go project and repository foundation.

---

## Task 0.1 — Initialize Go Module

### Requirements

* Create `go.mod`.
* Use a clear module path.
* Use a currently supported stable Go version.
* Do not add unnecessary dependencies.

### Acceptance Criteria

* `go.mod` exists.
* `go build ./...` succeeds.

---

## Task 0.2 — Create Repository Structure

Create the initial directory structure defined in `AGENTS.md`.

At minimum:

```text
cmd/
internal/
pkg/
ui/
resources/
tests/
```

### Acceptance Criteria

* Repository structure is created.
* No unnecessary top-level directories are added.

---

## Task 0.3 — Create Application Entry Points

Create:

```text
cmd/devbox/main.go
cmd/devbox-cli/main.go
```

### Requirements

* Both applications compile.
* Both use the application/core layer.
* Do not put business logic in `main.go`.

### Acceptance Criteria

```bash
go build ./cmd/devbox
go build ./cmd/devbox-cli
```

both succeed.

---

## Task 0.4 — Add Basic Logger

Create the logging abstraction.

### Requirements

Support at least:

* Info
* Warn
* Error
* Debug

Log categories should be supported conceptually:

```text
application
service
runtime
package
process
configuration
database
```

### Acceptance Criteria

* Logger is reusable.
* Core code does not use random `fmt.Println` for application logging.
* No sensitive information is logged.

---

# 4. Milestone 1 — Filesystem and Environment Foundation

## Goal

Create and manage the DevBox environment directory.

---

## Task 1.1 — Implement Path Manager

Create a path abstraction.

It must resolve:

```text
DevBox Root
bin/
active/
www/
data/
etc/
logs/
backup/
packages/
cache/
```

### Requirements

Do not hardcode personal paths.

Example API:

```go
type Paths interface {
    Root() string
    Bin() string
    Active() string
    WWW() string
    Data() string
    Etc() string
    Logs() string
    Backup() string
    Packages() string
    Cache() string
}
```

### Acceptance Criteria

* Paths are resolved from the DevBox root.
* Tests do not depend on a developer's machine.

---

## Task 1.2 — Implement Environment Initializer

Create the directory initializer.

### Requirements

The initializer must:

* Create missing directories.
* Be idempotent.
* Never delete existing user data.

### Required directories

```text
bin/
active/
www/
data/
etc/
logs/
backup/
packages/
cache/
```

### Acceptance Criteria

* Running initialization twice produces no errors.
* Existing files remain untouched.
* Tests use temporary directories.

---

## Task 1.3 — Add Filesystem Tests

Test:

* Path resolution.
* Directory creation.
* Repeated initialization.
* Existing data preservation.

### Acceptance Criteria

```bash
go test ./...
```

passes.

---

# 5. Milestone 2 — Process Manager

## Goal

Create a safe platform-aware process management abstraction.

---

## Task 2.1 — Define Process Model

Create a process model.

It should represent:

* PID
* Executable
* Start time where useful
* Running status

Do not expose `os/exec.Cmd` to the entire application.

---

## Task 2.2 — Define Process Manager Interface

Create:

```go
type Manager interface {
    Start(config StartConfig) (*Process, error)
    Stop(pid int) error
    IsRunning(pid int) bool
}
```

`StartConfig` must support:

```text
Executable
Arguments
Working Directory
Environment
```

---

## Task 2.3 — Implement Windows Process Support

Windows is the Phase 1 primary platform.

### Requirements

* Direct process execution.
* No shell command concatenation.
* Arguments passed separately.
* Working directory supported.

### Security Requirement

Do not execute user input through:

```text
cmd.exe /c
```

unless explicitly required and safely controlled.

---

## Task 2.4 — Define Linux and macOS Boundaries

Create platform-specific architecture boundaries.

The application must compile with a design that allows:

```text
process_windows.go
process_linux.go
process_darwin.go
```

Full Linux/macOS process implementation may be completed in a later phase if necessary.

### Acceptance Criteria

* Generic code does not contain Windows-only process logic.
* Platform-specific code is isolated.

---

# 6. Milestone 3 — Service System

## Goal

Create a generic service architecture.

---

## Task 3.1 — Define Service Interface

Implement:

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

---

## Task 3.2 — Define Service Status

Implement:

```text
Unknown
Stopped
Starting
Running
Stopping
Error
```

Use a typed status.

Do not use arbitrary strings throughout the application.

---

## Task 3.3 — Implement Service Manager

Required operations:

```text
Register
Get
All
Start
Stop
Restart
Status
```

### Requirements

* Services registered by ID.
* Duplicate registration must be handled.
* Unknown service IDs must return useful errors.

---

## Task 3.4 — Add Service Tests

Test:

* Service registration.
* Duplicate registration.
* Service lookup.
* Unknown service.
* Service status handling.

---

# 7. Milestone 4 — Apache Service

## Goal

Integrate Apache as the first service.

---

## Task 4.1 — Create Apache Service

Create:

```text
internal/service/apache/
```

The Apache service must:

* Resolve Apache installation.
* Resolve executable.
* Resolve configuration.
* Start Apache.
* Stop Apache.
* Restart Apache.
* Detect status.

---

## Task 4.2 — Apache Path Resolution

Apache should be resolved from:

```text
bin/apache/<version>/
```

Do not hardcode:

```text
C:\Apache
```

or another machine-specific path.

---

## Task 4.3 — Apache Process Startup

Use the Process Manager.

Apache startup must support:

```text
Executable
Arguments
Working Directory
```

The configuration file must be resolved through the configuration/path system.

---

## Task 4.4 — Apache Status Detection

Implement reliable status detection.

Do not rely only on:

```text
"Start() returned nil"
```

A process can start and immediately crash.

Status detection must account for actual process state.

---

## Task 4.5 — Apache Logs

Implement basic log access.

At minimum support:

```text
Error log path
Access log path where available
```

Do not load huge log files entirely into memory unnecessarily.

---

## Task 4.6 — Apache Tests

Test:

* Path resolution.
* Missing installation.
* Missing executable.
* Service state behavior.

Use mocks/fakes where appropriate.

Do not require a real Apache installation for unit tests.

---

# 8. Milestone 5 — MariaDB Service

## Goal

Integrate MariaDB as the second service.

---

## Task 5.1 — Create MariaDB Service

Create:

```text
internal/service/mariadb/
```

Implement:

* Start.
* Stop.
* Restart.
* Status.
* Installation detection.

---

## Task 5.2 — Separate Binary and Data

MariaDB binaries:

```text
bin/mariadb/<version>/
```

MariaDB data:

```text
data/mariadb/
```

Never place user data inside the binary directory.

---

## Task 5.3 — MariaDB Configuration

Resolve configuration through:

```text
etc/mariadb/
```

Do not hardcode configuration paths.

---

## Task 5.4 — MariaDB Status

Implement status detection.

A successful process start is not sufficient.

---

## Task 5.5 — MariaDB Tests

Test:

* Installation detection.
* Missing executable.
* Data directory resolution.
* Service status behavior.

---

# 9. Milestone 6 — Runtime System

## Goal

Create the runtime architecture.

---

## Task 6.1 — Define Runtime Interface

Create a runtime abstraction.

PHP is the first runtime.

The architecture must support future:

```text
Python
Node.js
```

without rewriting the core.

---

## Task 6.2 — Create PHP Runtime

Create:

```text
internal/runtime/php/
```

Implement:

* PHP installation discovery.
* Installed version listing.
* Active version detection.
* PHP version switching.

---

## Task 6.3 — PHP Version Discovery

Discover versions from:

```text
bin/php/
```

Example:

```text
8.1.29
8.2.20
8.3.30
```

The discovery process must validate that the directory represents a usable PHP installation.

Do not treat every directory as a valid PHP version.

---

## Task 6.4 — Active PHP State

Implement active PHP tracking.

Use:

```text
active/php
```

or a platform-safe equivalent.

The application must have a single authoritative active PHP version.

---

# 10. Milestone 7 — PHP Version Switching

## Goal

Implement safe PHP switching.

---

## Task 7.1 — Validate Target PHP Version

Before switching:

* Confirm version exists.
* Confirm installation is valid.
* Confirm required PHP files exist.

---

## Task 7.2 — Implement Atomic Switching

The switch must follow:

```text
Validate
    ↓
Prepare
    ↓
Update
    ↓
Verify
    ↓
Commit
```

If anything fails:

```text
Rollback
```

---

## Task 7.3 — Apache PHP Integration

Ensure Apache uses the active PHP version.

Do not duplicate PHP installations.

Do not copy all PHP files into an `active` directory.

---

## Task 7.4 — PHP Switching Tests

Test:

* Successful switching.
* Missing version.
* Invalid installation.
* Failed switch rollback.
* Active version remains unchanged after failure.

---

# 11. Milestone 8 — Package Metadata System

## Goal

Create the foundation for future package installation.

Phase 1 does not require a complete public package registry.

---

## Task 8.1 — Define Package Model

Package model must support:

```text
ID
Name
Version
Platform
Architecture
Checksum
Download URL
Install Path
Package Type
```

---

## Task 8.2 — Implement Manifest Parser

Create manifest parsing.

Use JSON initially.

Example:

```json
{
  "id": "php",
  "version": "8.3.30",
  "platform": "windows",
  "architecture": "x64",
  "download": {
    "url": "...",
    "sha256": "..."
  }
}
```

---

## Task 8.3 — Implement Package Validation

Validate:

* Required fields.
* Supported platform.
* Supported architecture.
* Valid checksum format.

---

## Task 8.4 — Implement Checksum Verification

Support SHA-256.

Never mark a package as valid before checksum verification.

---

# 12. Milestone 9 — SQLite Database

## Goal

Persist DevBox metadata.

---

## Task 9.1 — Add SQLite

Add SQLite using a suitable lightweight Go implementation.

Before adding a dependency:

* Review license.
* Review maintenance.
* Review cross-platform support.

---

## Task 9.2 — Create Database Layer

Create:

```text
internal/database/
```

The database layer must:

* Open database.
* Initialize schema.
* Run migrations.
* Close cleanly.

---

## Task 9.3 — Create Initial Schema

Initial metadata should support:

```text
settings
installed_packages
projects
```

Do not over-design the schema.

---

## Task 9.4 — Implement Migrations

Migrations must be versioned.

Never silently destroy data.

---

# 13. Milestone 10 — Application Layer

## Goal

Create a unified application layer for CLI and GUI.

---

## Task 10.1 — Create Application Container

The application container should initialize:

```text
Logger
Paths
Database
Process Manager
Service Manager
Runtime Manager
Package Manager
```

---

## Task 10.2 — Create Application API

The CLI and GUI should call application-level operations.

Examples:

```text
StartService
StopService
RestartService
GetServiceStatus
ListPHPVersions
GetActivePHP
UsePHPVersion
```

The UI must not access internal infrastructure directly.

---

# 14. Milestone 11 — CLI

## Goal

Build the first usable DevBox interface.

---

## Task 11.1 — Implement `status`

Command:

```bash
devbox status
```

Display:

```text
Apache     Running
MariaDB    Stopped
PHP        8.3.30
```

---

## Task 11.2 — Implement Service Commands

Commands:

```bash
devbox start apache
devbox stop apache
devbox restart apache
```

Also support MariaDB.

---

## Task 11.3 — Implement PHP Commands

Commands:

```bash
devbox php list
devbox php current
devbox php use 8.3.30
```

---

## Task 11.4 — CLI Error Handling

The CLI must:

* Return non-zero exit codes on failure.
* Display useful errors.
* Avoid stack traces for normal failures.

---

# 15. Milestone 12 — Native GUI

## Goal

Build a minimal Fyne GUI.

---

## Task 12.1 — Create Fyne Application

Create:

```text
ui/
```

The GUI must use the application layer.

---

## Task 12.2 — Build Dashboard

Display:

```text
Apache
MariaDB
PHP
```

Each service must show:

```text
Name
Status
Action
```

---

## Task 12.3 — Service Actions

Add:

```text
Start
Stop
Restart
```

The GUI must call application methods.

---

## Task 12.4 — PHP Selector

Add a PHP version selector.

The selector must:

* Display installed PHP versions.
* Display active PHP.
* Allow switching.
* Display errors.

---

## Task 12.5 — GUI State Updates

The GUI must update when:

```text
Service status changes
PHP version changes
```

Avoid unnecessary aggressive polling.

---

# 16. Milestone 13 — Logs

## Goal

Provide basic log visibility.

---

## Task 13.1 — Application Logs

Implement application log storage.

---

## Task 13.2 — Service Logs

Allow basic access to:

```text
Apache logs
MariaDB logs
```

Do not load massive files fully into memory.

---

## Task 13.3 — Log UI

Add a simple log view if practical.

Do not build a full log analysis system.

---

# 17. Milestone 14 — Testing and Stabilization

## Goal

Make Phase 1 stable.

---

## Task 14.1 — Run Full Test Suite

Run:

```bash
go test ./...
```

Fix all failures.

---

## Task 14.2 — Run Static Checks

Run:

```bash
go vet ./...
```

Fix issues caused by the project.

---

## Task 14.3 — Format Code

Run Go formatting.

All Go files must be formatted.

---

## Task 14.4 — Review Error Handling

Review all:

* File operations.
* Process operations.
* Database operations.
* Configuration operations.

No important errors may be silently ignored.

---

## Task 14.5 — Review Data Safety

Confirm:

* `www/` is never deleted automatically.
* `data/` is never deleted automatically.
* `backup/` is never deleted automatically.
* PHP switching can roll back.
* Package installation does not partially overwrite a valid installation.

---

# 18. Phase 1 Final Acceptance Criteria

Phase 1 is complete only when all of the following are true.

## Architecture

* [ ] Go core exists.
* [ ] GUI uses Fyne.
* [ ] No Electron.
* [ ] No React.
* [ ] No WebView.
* [ ] CLI and GUI share the core.
* [ ] Platform-specific process code is isolated.

## Environment

* [ ] DevBox root is resolved.
* [ ] Required directories are created.
* [ ] Initialization is idempotent.
* [ ] User data is preserved.

## Apache

* [ ] Apache can be detected.
* [ ] Apache can start.
* [ ] Apache can stop.
* [ ] Apache can restart.
* [ ] Apache status can be detected.

## MariaDB

* [ ] MariaDB can be detected.
* [ ] MariaDB can start.
* [ ] MariaDB can stop.
* [ ] MariaDB can restart.
* [ ] MariaDB status can be detected.
* [ ] Data directory is separate from binaries.

## PHP

* [ ] PHP versions can be discovered.
* [ ] Active PHP can be identified.
* [ ] PHP versions can be switched.
* [ ] Failed switching can roll back.

## Package System

* [ ] Package manifests can be parsed.
* [ ] Package metadata can be validated.
* [ ] SHA-256 checksums can be verified.

## Database

* [ ] SQLite is integrated.
* [ ] Migrations work.
* [ ] Metadata can be persisted.

## CLI

* [ ] `status` works.
* [ ] Service commands work.
* [ ] PHP commands work.
* [ ] Exit codes are correct.

## GUI

* [ ] Fyne application starts.
* [ ] Dashboard displays services.
* [ ] Service actions work.
* [ ] PHP selector works.
* [ ] GUI uses the core layer.

## Quality

* [ ] `go test ./...` passes.
* [ ] `go vet ./...` passes.
* [ ] Code is formatted.
* [ ] No personal paths are hardcoded.
* [ ] No unrelated features were added.

---

# 19. After Phase 1

Do not immediately implement every planned feature.

The next step should be a **Phase 1 review**.

Evaluate:

* Startup time.
* Memory usage.
* Process reliability.
* PHP switching reliability.
* Apache compatibility.
* MariaDB stability.
* GUI usability.
* Cross-platform architecture.

Only after this review should Phase 2 begin.

Potential Phase 2 features:

```text
Nginx
Redis
Mailpit
Python
Node.js
Composer
Virtual Hosts
SSL
Project Detection
```

The Phase 1 foundation must remain stable before adding additional runtimes and services.
