# DECISIONS.md

# WorkBench — Architecture Decision Record

This document records important architectural and product decisions made for WorkBench.

The purpose of this document is to prevent future contributors and AI coding agents from repeatedly reconsidering already-decided architectural choices without a strong reason.

Before changing a decision:

1. Understand the original reason.
2. Identify the actual problem with the current decision.
3. Evaluate the impact.
4. Create a new decision record or update the relevant decision.

---

# Decision Statuses

| Status     | Meaning                      |
| ---------- | ---------------------------- |
| Accepted   | Current project decision     |
| Proposed   | Under discussion             |
| Deprecated | No longer recommended        |
| Superseded | Replaced by another decision |

---

# ADR-001 — Use Go as the Primary Programming Language

**Status:** Accepted

## Decision

WorkBench will be implemented primarily using **Go**.

## Context

WorkBench must be:

* Fast.
* Lightweight.
* Cross-platform.
* Easy to distribute.
* Easy to maintain.
* Suitable for system-level process and filesystem management.

The project also requires a native GUI and a shared core engine for both GUI and CLI.

## Reasons

Go provides:

* Simple concurrency.
* Fast compilation.
* Strong standard library.
* Good process and filesystem APIs.
* Easy cross-compilation.
* Simple deployment.
* Low operational complexity.

Go is also a good fit for a developer tool that manages local processes and files.

## Alternatives Considered

### Rust

Rust provides excellent performance and memory safety.

However, for this project:

* Development complexity is higher.
* The project requires rapid iteration.
* Go provides sufficient performance for the WorkBench use case.

### C++

C++ is powerful and suitable for native applications.

However:

* Build complexity is higher.
* Cross-platform dependency management is more complex.
* Memory safety requires greater discipline.

## Consequence

WorkBench will use Go for:

* Core engine.
* CLI.
* Service management.
* Runtime management.
* Package management.
* Application layer.

---

# ADR-002 — Use Fyne for the Native GUI

**Status:** Accepted

## Decision

WorkBench will use **Fyne** for the desktop GUI.

## Context

The GUI must be:

* Native-oriented.
* Lightweight.
* Cross-platform.
* Written in Go.
* Simple to maintain.

## Reasons

Fyne allows the GUI to remain inside the Go ecosystem.

This avoids creating a separate frontend technology stack.

## Alternatives Considered

### Electron

Rejected.

Reasons:

* High memory usage.
* Bundled browser runtime.
* Heavy desktop application footprint.
* Unnecessary web technology.

### React

Rejected for the desktop GUI.

React is a web UI technology and would require a separate desktop runtime or web-based architecture.

### WebView

Rejected.

The GUI should not depend on a browser engine.

### Qt / C++

Not selected because the project is primarily implemented in Go.

## Consequence

The GUI must use Fyne.

The GUI must not introduce:

* Electron.
* React.
* Vue.
* WebView-based architecture.

---

# ADR-003 — The Go Core Engine Is the Product

**Status:** Accepted

## Decision

The WorkBench Core Engine is the primary product.

The GUI and CLI are clients of the core engine.

## Architecture

```text
GUI ───────┐
           ▼
      Application
           ▼
       Core Engine
           ▲
           │
CLI ───────┘
```

## Context

WorkBench must provide both a GUI and CLI.

Duplicating logic between the GUI and CLI would create inconsistent behavior.

## Consequence

The following must be implemented in the core/application layers:

* Service operations.
* PHP version switching.
* Environment management.
* Package management.

The GUI and CLI must not implement separate versions of the same business logic.

---

# ADR-004 — Do Not Use Web Technology for the Desktop GUI

**Status:** Accepted

## Decision

WorkBench will not use web technology as the desktop GUI architecture.

## Rejected Technologies

* Electron.
* React.
* Vue.
* WebView.
* Browser-based desktop UI.

## Context

The primary goal of WorkBench is to provide a fast and lightweight alternative to existing local development environment managers.

A web-based desktop architecture introduces unnecessary runtime overhead.

## Consequence

The GUI must remain native-oriented and use the selected Go GUI technology.

---

# ADR-005 — Use a Modular Service Architecture

**Status:** Accepted

## Decision

Each long-running development service must be implemented as a modular Service.

Examples:

```text
Apache
MariaDB
Nginx
Redis
Mailpit
```

## Context

WorkBench will support more services in the future.

The system must not contain service-specific conditionals throughout the core.

## Consequence

Adding a service should generally follow:

```text
Create Service Module
        ↓
Implement Service Interface
        ↓
Use Existing Infrastructure
        ↓
Register Service
```

The Service Manager must not contain logic such as:

```go
if service == "apache" {
    ...
}
```

---

# ADR-006 — Runtime and Service Are Different Concepts

**Status:** Accepted

## Decision

A Runtime and a Service are separate domain concepts.

## Examples

```text
PHP      → Runtime
Apache   → Service
MariaDB  → Service
```

## Context

PHP is a language runtime.

Apache and MariaDB are long-running processes.

They have different lifecycle and version-management requirements.

## Consequence

PHP must not be modeled as a Service.

Future runtimes such as:

```text
Python
Node.js
```

must use the Runtime architecture.

---

# ADR-007 — PHP Versions Are Installed Independently

**Status:** Accepted

## Decision

Each PHP version must have its own installation directory.

Example:

```text
bin/php/
├── 8.1.29/
├── 8.2.20/
└── 8.3.30/
```

## Context

Developers frequently need to switch between PHP versions.

The active PHP version must be changed without modifying or corrupting installed versions.

## Consequence

PHP version switching must not:

* Copy all PHP files.
* Overwrite another PHP installation.
* Modify the original installation.

---

# ADR-008 — Use an Active PHP State Instead of Copying PHP

**Status:** Accepted

## Decision

WorkBench will track the active PHP version separately.

Example:

```text
active/php
```

## Context

Copying PHP installations creates:

* Slow switching.
* Duplicate files.
* Increased disk usage.
* Risk of partial state.

## Decision

The active PHP version will be represented by an active state mechanism.

The exact platform-specific implementation may differ.

## Consequence

Switching PHP should be lightweight and transactional.

---

# ADR-009 — PHP Version Switching Must Be Transactional

**Status:** Accepted

## Decision

PHP version switching must follow:

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

Failure must trigger rollback.

## Context

A failed PHP switch must not leave WorkBench in an inconsistent state.

## Consequence

If switching from:

```text
PHP 8.3
```

to:

```text
PHP 8.2
```

fails, PHP 8.3 must remain active.

---

# ADR-010 — Keep Software Binaries Separate from User Data

**Status:** Accepted

## Decision

Software binaries and persistent data must use separate directories.

Example:

```text
bin/mariadb/
data/mariadb/
```

## Context

Software upgrades and package changes must not risk user data.

## Consequence

Installing or updating MariaDB must not automatically delete or overwrite the MariaDB data directory.

---

# ADR-011 — `www/`, `data/`, and `backup/` Are User-Sensitive Directories

**Status:** Accepted

## Decision

WorkBench must treat the following directories as user-sensitive:

```text
www/
data/
backup/
```

## Consequence

WorkBench must never automatically delete these directories during:

* Service operations.
* Runtime switching.
* Package installation.
* Application startup.

Any future destructive operation must require explicit user intent and appropriate safeguards.

---

# ADR-012 — Use SQLite for WorkBench Metadata

**Status:** Accepted

## Decision

WorkBench will use SQLite for local metadata.

## Context

WorkBench needs to store:

* Installed package metadata.
* Settings.
* Project metadata.
* Migration state.

A full database server is inappropriate for a local environment manager.

## Reasons

SQLite is:

* Embedded.
* Lightweight.
* Cross-platform.
* Reliable.
* Easy to distribute.

## Consequence

SQLite stores metadata only.

SQLite is not the source of truth for software binaries.

---

# ADR-013 — Filesystem Reality Takes Precedence over Metadata

**Status:** Accepted

## Decision

The filesystem is the source of truth for installed software.

## Example

If SQLite says:

```text
PHP 8.3.30 is installed
```

but the installation directory does not exist, WorkBench must treat PHP 8.3.30 as unavailable.

## Context

Users may manually modify or remove files.

Metadata can become stale.

## Consequence

WorkBench must be able to detect and reconcile stale metadata.

---

# ADR-014 — Package Installation Must Use Temporary Extraction

**Status:** Accepted

## Decision

Packages must be extracted into a temporary location before being moved to the final installation directory.

## Flow

```text
Download
   ↓
Verify Checksum
   ↓
Extract Temporary
   ↓
Validate
   ↓
Move Final
```

## Context

Direct extraction into the final directory can create partial installations.

## Consequence

A failed package installation must not corrupt an existing valid installation.

---

# ADR-015 — Verify Package Checksums

**Status:** Accepted

## Decision

WorkBench must verify package checksums before installation.

SHA-256 is the initial supported checksum algorithm.

## Context

WorkBench will eventually download software packages.

Downloaded files must not be blindly trusted.

## Consequence

A package that fails checksum validation must not be installed.

---

# ADR-016 — Avoid Shell Command Concatenation

**Status:** Accepted

## Decision

WorkBench must not build shell commands by concatenating user input.

Bad:

```go
command := "start " + userInput
```

## Decision

Process arguments must be passed as structured arguments.

## Context

WorkBench executes external software.

Unsafe shell command construction creates unnecessary security risks.

## Consequence

The Process Manager owns process execution.

Services must not implement unsafe shell execution.

---

# ADR-017 — Abstract Process Management

**Status:** Accepted

## Decision

WorkBench will use a Process Manager abstraction.

## Context

Apache and MariaDB require process management.

The application must also support Windows, Linux, and macOS.

## Consequence

Services must not directly depend on operating system process APIs.

Platform-specific process code must be isolated.

---

# ADR-018 — Isolate Platform-Specific Code

**Status:** Accepted

## Decision

Platform-specific behavior must be isolated using platform-specific Go files and implementations.

Examples:

```text
process_windows.go
process_linux.go
process_darwin.go
```

## Context

WorkBench is cross-platform.

Operating system behavior differs.

## Consequence

Generic application and domain code must not contain platform-specific implementation details.

---

# ADR-019 — Use Constructor Dependency Injection

**Status:** Accepted

## Decision

Important dependencies should be passed through constructors.

Example:

```go
func NewApacheService(
    paths filesystem.Paths,
    processes process.Manager,
) *ApacheService
```

## Context

Dependency injection improves:

* Testability.
* Explicit dependencies.
* Architecture clarity.

## Consequence

Avoid global mutable state.

Avoid hidden dependencies.

---

# ADR-020 — Prefer the Go Standard Library

**Status:** Accepted

## Decision

The Go standard library should be preferred where it provides suitable functionality.

## Context

WorkBench must remain lightweight.

Unnecessary dependencies increase:

* Binary complexity.
* Maintenance burden.
* Security review requirements.

## Consequence

A third-party dependency must have a clear justification.

---

# ADR-021 — No Automatic Network Activity on Startup

**Status:** Accepted

## Decision

WorkBench must not make unnecessary network requests during startup.

## Context

WorkBench is a local development environment manager.

Startup should be fast and predictable.

## Consequence

Package registry checks and updates must be explicit user actions or controlled background operations introduced in a future specification.

---

# ADR-022 — No Telemetry in the Initial Product

**Status:** Accepted

## Decision

WorkBench will not include telemetry in the initial product.

## Context

WorkBench is a free and open-source developer tool.

The initial priority is trust, simplicity, and privacy.

## Consequence

The application must not collect usage analytics by default.

Any future telemetry proposal requires a separate architectural decision.

---

# ADR-023 — Favour Explicit State over Aggressive Polling

**Status:** Accepted

## Decision

WorkBench should use application state and internal events where appropriate instead of continuously polling the filesystem or processes.

## Context

Aggressive polling can:

* Waste resources.
* Increase CPU usage.
* Make the application less efficient.

## Consequence

The application should use events such as:

```text
ServiceStarted
ServiceStopped
ServiceError
PHPVersionChanged
```

where practical.

Polling may still be used when required for process verification.

---

# ADR-024 — Phase 1 Must Remain Small

**Status:** Accepted

## Decision

Phase 1 will focus only on:

```text
Apache
PHP
MariaDB
```

## Context

The project has a large future roadmap.

Adding too many features early increases architectural risk.

## Consequence

The following must not be added to Phase 1 without an explicit roadmap change:

```text
Nginx
Redis
Python
Node.js
Mailpit
HeidiSQL
ngrok
```

---

# ADR-025 — The User Experience Is More Important Than Feature Count

**Status:** Accepted

## Decision

WorkBench prioritizes:

```text
Speed
Simplicity
Predictability
```

over the number of supported tools.

## Context

The purpose of WorkBench is to simplify local development environments.

A large number of poorly integrated features would contradict the product goal.

## Consequence

A feature should not be added only because another environment manager has it.

The feature must solve a real developer problem and fit the WorkBench architecture.

---

# 26. Changing an Architectural Decision

A decision may be changed when:

* The original assumptions are no longer valid.
* The current architecture causes a significant problem.
* A new requirement cannot be reasonably supported.
* A better solution has been demonstrated.

When changing a decision:

1. Identify the existing ADR.
2. Explain the problem.
3. Explain the alternatives.
4. Explain the new decision.
5. Document consequences.

Do not silently reverse an accepted decision.

---

# 27. Final Architectural Principle

WorkBench must remain:

> **A small native core with modular services and runtimes.**

The most important architectural rule is:

```text
Do not make the system more complicated than the developer problem requires.
```
