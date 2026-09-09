# ARCHITECTURE.md

# WorkBench — System Architecture

This document defines the technical architecture of WorkBench.

WorkBench is a native, lightweight, cross-platform local development environment manager written in Go.

This document defines:

* System boundaries
* Core modules
* Domain concepts
* Data flow
* Interfaces
* Dependency rules
* Runtime and service architecture

Read `AGENTS.md` and `ROADMAP.md` before implementing features.

---

# 1. Architectural Philosophy

WorkBench follows this principle:

> **A small native application with a strong core engine and modular development components.**

The system must remain:

* Simple
* Fast
* Modular
* Testable
* Cross-platform
* Native

WorkBench is **not**:

* A web application
* A cloud platform
* A package marketplace
* An enterprise orchestration platform

---

# 2. System Overview

```text
┌─────────────────────────────────────────────────────┐
│                    USER                             │
└──────────────┬───────────────────────┬──────────────┘
               │                       │
               ▼                       ▼
┌────────────────────────┐   ┌────────────────────────┐
│      Native GUI        │   │          CLI           │
│         Fyne           │   │        Go CLI          │
└──────────────┬─────────┘   └──────────────┬─────────┘
               │                            │
               └──────────────┬─────────────┘
                              ▼
┌─────────────────────────────────────────────────────┐
│              APPLICATION LAYER                      │
│                                                     │
│  Application Commands                               │
│  Application State                                  │
│  Application Events                                 │
└─────────────────────────┬───────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────┐
│                   CORE ENGINE                       │
│                                                     │
│  Service Manager                                    │
│  Runtime Manager                                    │
│  Package Manager                                    │
│  Environment Manager                                │
│  Project Manager                                    │
└──────────────┬───────────────┬───────────────┬──────┘
               │               │               │
               ▼               ▼               ▼
        ┌────────────┐ ┌─────────────┐ ┌──────────────┐
        │  SERVICES  │ │  RUNTIMES   │ │   PACKAGES   │
        │            │ │             │ │              │
        │ Apache     │ │ PHP         │ │ Manifests    │
        │ MariaDB    │ │             │ │ Downloads    │
        └──────┬─────┘ └──────┬──────┘ └──────┬───────┘
               │               │               │
               └───────────────┼───────────────┘
                               ▼
┌─────────────────────────────────────────────────────┐
│                INFRASTRUCTURE                       │
│                                                     │
│  Process Manager                                    │
│  Filesystem                                         │
│  Configuration                                      │
│  SQLite                                             │
│  HTTP Downloader                                    │
└─────────────────────────────────────────────────────┘
```

---

# 3. Architectural Layers

WorkBench is organized into five conceptual layers.

```text
Layer 1: Presentation
Layer 2: Application
Layer 3: Core
Layer 4: Domain Modules
Layer 5: Infrastructure
```

---

# 4. Layer 1 — Presentation

## Responsibility

The presentation layer interacts with the user.

Implementations:

```text
GUI
CLI
```

Current GUI:

```text
Fyne
```

The presentation layer may:

* Display information.
* Receive user actions.
* Format errors.
* Trigger application commands.

The presentation layer MUST NOT:

* Start processes directly.
* Read Apache configuration directly.
* Write PHP configuration directly.
* Access SQLite directly.
* Manipulate service binaries directly.

---

# 5. Layer 2 — Application

## Responsibility

The application layer coordinates user actions.

The application layer exposes high-level use cases.

Examples:

```text
StartService
StopService
RestartService
GetServiceStatus
ListPHPVersions
GetActivePHP
UsePHPVersion
GetEnvironmentStatus
```

The application layer is the primary boundary between:

```text
GUI / CLI
```

and:

```text
Core Engine
```

---

# 6. Layer 3 — Core Engine

## Responsibility

The core engine coordinates WorkBench functionality.

Core managers:

```text
Service Manager
Runtime Manager
Package Manager
Environment Manager
Project Manager
```

The core engine must not depend on Fyne.

The CLI must be able to run using the core engine without initializing the GUI.

---

# 7. Layer 4 — Domain Modules

Domain modules represent WorkBench concepts.

Primary domain concepts:

```text
Service
Runtime
Package
Project
Environment
```

These concepts should not be tightly coupled to operating system APIs.

---

# 8. Layer 5 — Infrastructure

Infrastructure communicates with the operating system and external systems.

Infrastructure includes:

```text
Process Manager
Filesystem
Configuration Manager
SQLite
HTTP Downloader
Archive Extractor
```

Infrastructure implementations may be platform-specific.

---

# 9. Dependency Direction

Dependencies must flow inward.

Correct:

```text
GUI
 ↓
Application
 ↓
Core
 ↓
Domain
 ↓
Infrastructure
```

Infrastructure must not control the application layer.

Incorrect:

```text
Filesystem → GUI
```

Incorrect:

```text
Fyne → ApacheService internal fields
```

---

# 10. Repository Architecture

The repository is organized as follows:

```text
workbench/
│
├── cmd/
│
├── internal/
│   ├── app/
│   ├── core/
│   ├── service/
│   ├── runtime/
│   ├── package/
│   ├── process/
│   ├── config/
│   ├── filesystem/
│   ├── project/
│   ├── database/
│   └── logger/
│
├── pkg/
│   ├── models/
│   └── events/
│
├── ui/
│
└── resources/
```

---

# 11. Application Container

The application container initializes system dependencies.

Conceptually:

```text
Application
│
├── Logger
├── Paths
├── Database
├── Process Manager
├── Configuration Manager
├── Service Manager
├── Runtime Manager
└── Package Manager
```

The application container is responsible for dependency composition.

It should be possible to create:

```text
Production Application
Test Application
```

with different infrastructure implementations.

---

# 12. Service Domain

A service is a long-running development process.

Examples:

```text
Apache
MariaDB
Nginx
Redis
Mailpit
```

---

## 12.1 Service Interface

Conceptually:

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

## 12.2 Service Lifecycle

```text
Unknown
   │
   ▼
Stopped
   │
   ▼
Starting
   │
   ▼
Running
   │
   ▼
Stopping
   │
   ▼
Stopped
```

Error path:

```text
Starting
   │
   ▼
Error
```

---

## 12.3 Service Manager

The Service Manager owns registered services.

```text
Service Manager
│
├── apache
└── mariadb
```

Future:

```text
├── nginx
├── redis
└── mailpit
```

The Service Manager does not know the implementation details of each service.

---

# 13. Apache Architecture

Apache is a Service.

```text
ApacheService
     │
     ├── Installation Resolver
     ├── Configuration Resolver
     ├── Process Manager
     └── Log Resolver
```

Startup flow:

```text
Start Apache
      │
      ▼
Resolve Installation
      │
      ▼
Resolve Executable
      │
      ▼
Resolve Configuration
      │
      ▼
Build Process Configuration
      │
      ▼
Process Manager.Start()
      │
      ▼
Verify Process
      │
      ▼
Running
```

Apache must not directly manage operating system process APIs.

---

# 14. MariaDB Architecture

MariaDB is a Service.

```text
MariaDBService
      │
      ├── Installation Resolver
      ├── Data Directory Resolver
      ├── Configuration Resolver
      └── Process Manager
```

Important separation:

```text
MariaDB Binary
        │
        ▼
bin/mariadb/<version>/

MariaDB Data
        │
        ▼
data/mariadb/
```

The MariaDB binary lifecycle must never delete the data directory.

---

# 15. Runtime Domain

A Runtime is a language execution environment.

Examples:

```text
PHP
Python
Node.js
```

PHP is the first implementation.

---

## 15.1 Runtime Interface

Conceptually:

```go
type Runtime interface {
    ID() string
    Name() string

    InstalledVersions() []Version
    ActiveVersion() *Version
    Use(version string) error
}
```

---

# 16. PHP Architecture

PHP is a Runtime.

```text
PHP Runtime
│
├── Version Discovery
├── Version Validation
├── Active Version
├── Version Switching
└── Integration Coordination
```

PHP installations:

```text
bin/php/
├── 8.1.29/
├── 8.2.20/
└── 8.3.30/
```

Active runtime:

```text
active/php
```

---

# 17. PHP Version Switching

PHP switching is a transaction.

```text
┌───────────────┐
│ Target PHP    │
└──────┬────────┘
       ▼
┌───────────────┐
│ Validate      │
└──────┬────────┘
       ▼
┌───────────────┐
│ Prepare       │
└──────┬────────┘
       ▼
┌───────────────┐
│ Update        │
└──────┬────────┘
       ▼
┌───────────────┐
│ Verify        │
└──────┬────────┘
       ▼
┌───────────────┐
│ Commit        │
└───────────────┘
```

Failure:

```text
Failure
   │
   ▼
Rollback
   │
   ▼
Previous PHP remains active
```

---

# 18. Service and Runtime Relationship

Services and runtimes are different concepts.

```text
PHP
 │
 └── Runtime

Apache
 │
 └── Service
```

Apache may depend on PHP.

Conceptually:

```text
Apache
  │
  ▼
PHP Runtime
```

The dependency is:

```text
Apache → Active PHP
```

The PHP Runtime must not directly control Apache.

The application/core layer coordinates changes.

---

# 19. Package Domain

A Package represents installable software metadata.

Examples:

```text
PHP 8.3.30
Apache 2.4.x
MariaDB 11.x
```

Package metadata is separate from the installed binary.

---

# 20. Package Installation Flow

```text
Resolve Manifest
       │
       ▼
Download
       │
       ▼
Verify Checksum
       │
       ▼
Extract to Temporary Directory
       │
       ▼
Validate Installation
       │
       ▼
Move to Final Directory
       │
       ▼
Register Metadata
```

Never:

```text
Download
   ↓
Extract directly into final directory
```

---

# 21. Package Manager

The Package Manager coordinates:

```text
Package Repository
Package Downloader
Checksum Validator
Archive Extractor
Package Registry
```

Conceptually:

```text
Package Manager
│
├── Repository
├── Downloader
├── Validator
├── Extractor
└── Registry
```

---

# 22. Package Repository

The repository provides package metadata.

Initial implementation may use:

```text
Local manifests
```

Future implementation may use:

```text
Remote package registry
```

The Package Manager must not assume that packages are always remote.

---

# 23. Filesystem Architecture

The filesystem layer owns path and file operations.

It must provide abstractions for:

```text
Path Resolution
Directory Creation
File Existence
File Copy
File Move
File Removal
Archive Extraction
```

The filesystem layer must not contain Apache or PHP business logic.

---

# 24. WorkBench Environment

The environment is represented by a root directory.

```text
WorkBench/
```

The root contains:

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

---

# 25. Active Runtime Model

The active runtime is a separate concept from installed runtimes.

Example:

```text
Installed:
8.1.29
8.2.20
8.3.30

Active:
8.3.30
```

Do not modify or delete installed runtimes during switching.

---

# 26. Configuration Architecture

Configuration files are infrastructure resources.

Examples:

```text
etc/apache/
etc/php/
etc/mariadb/
```

Configuration updates must be coordinated through a configuration manager.

---

## 26.1 Configuration Update Flow

```text
Load
  │
  ▼
Backup
  │
  ▼
Modify
  │
  ▼
Validate
  │
  ▼
Save
```

If validation fails:

```text
Restore Backup
```

---

# 27. Process Architecture

The Process Manager abstracts operating system process operations.

```text
Service
   │
   ▼
Process Manager
   │
   ▼
Operating System
```

The Service layer must not directly use:

```go
exec.Command()
```

unless the process abstraction itself is being implemented.

---

# 28. Platform Architecture

Generic code:

```text
internal/process/
```

Platform-specific implementations:

```text
process_windows.go
process_linux.go
process_darwin.go
```

The same conceptual interface must be used on every platform.

---

# 29. Database Architecture

SQLite stores metadata.

The database is not the source of truth for binaries.

```text
Filesystem
    +
SQLite Metadata
```

If they disagree:

```text
Filesystem Reality Wins
```

The application should reconcile stale metadata.

---

# 30. Event Architecture

WorkBench uses lightweight internal events.

Example:

```text
PHPVersionChanged
ServiceStarted
ServiceStopped
ServiceError
PackageInstalled
```

Event flow:

```text
Core
  │
  ▼
Event Bus
  │
  ├── GUI
  ├── Logger
  └── Application State
```

Events should be used to reduce unnecessary polling.

---

# 31. Application State

Application state represents the current WorkBench state.

Example:

```text
Apache: Running
MariaDB: Stopped
PHP: 8.3.30
```

The GUI consumes application state.

The GUI does not reconstruct application state by independently scanning the filesystem.

---

# 32. CLI Architecture

The CLI is a presentation layer.

```text
CLI Command
    │
    ▼
Application Command
    │
    ▼
Core Manager
```

Example:

```text
workbench php use 8.3.30
        │
        ▼
UsePHPVersion()
        │
        ▼
PHP Manager
```

---

# 33. GUI Architecture

The GUI is also a presentation layer.

```text
Fyne UI
   │
   ▼
Application Layer
   │
   ▼
Core Engine
```

The GUI must not contain domain logic.

Bad:

```text
if apacheRunning {
    exec.Command(...)
}
```

Good:

```text
app.StartService("apache")
```

---

# 34. Dependency Injection

Use constructor injection for important dependencies.

Example:

```go
func NewApacheService(
    paths filesystem.Paths,
    processes process.Manager,
    config config.Manager,
) *ApacheService
```

Avoid global mutable state.

Avoid global service registries.

The application container should compose dependencies.

---

# 35. Test Architecture

Core logic must be testable without:

* Apache installed.
* PHP installed.
* MariaDB installed.
* Internet access.

Use interfaces for external boundaries.

Test doubles may be used for:

```text
Process Manager
Package Downloader
Filesystem
```

Use real temporary directories for filesystem integration tests.

---

# 36. Error Architecture

Errors should preserve context.

Example:

```go
return fmt.Errorf("failed to start Apache: %w", err)
```

Errors should identify:

```text
What failed
Which component failed
Why it failed
```

Do not expose internal stack traces to normal users.

---

# 37. Security Architecture

WorkBench executes external software.

Security boundaries include:

```text
Package Download
Archive Extraction
Process Execution
Configuration Changes
```

Required protections:

* SHA-256 package verification.
* Archive traversal protection.
* No unsafe shell concatenation.
* No arbitrary command execution from user input.
* Safe temporary directories.
* Controlled executable paths.

---

# 38. Future Extension Architecture

Future services must be added as modules.

Example:

```text
internal/service/nginx/
internal/service/redis/
internal/service/mailpit/
```

Future runtimes:

```text
internal/runtime/python/
internal/runtime/node/
```

Future tools:

```text
internal/tool/heidisql/
internal/tool/ngrok/
```

Do not modify the entire core engine to add a new service.

The goal is:

```text
Add Module
    ↓
Register Module
    ↓
Use Existing Core Infrastructure
```

---

# 39. Architectural Anti-Patterns

Do not introduce these patterns.

## GUI-Owned Business Logic

```text
Fyne → Apache Process
```

## Service-Specific Filesystem Logic

```text
ApacheService → random os.WriteFile()
```

## Hardcoded Paths

```text
C:\WorkBench
```

## Shell Command Concatenation

```text
"start " + userInput
```

## Binary/Data Coupling

```text
bin/mariadb/data/
```

## Runtime Copying

```text
PHP 8.3 → active/php
```

## Database as Binary Truth

```text
SQLite says PHP exists
```

when the filesystem does not contain it.

---

# 40. Final Architecture Principle

The WorkBench architecture must preserve this dependency flow:

```text
┌──────────────┐
│     GUI      │
└──────┬───────┘
       │
┌──────▼───────┐
│     CLI      │
└──────┬───────┘
       │
┌──────▼────────────────────┐
│    Application Layer      │
└──────┬────────────────────┘
       │
┌──────▼────────────────────┐
│       Core Engine         │
└──────┬────────────────────┘
       │
┌──────▼────────────────────┐
│      Domain Modules       │
└──────┬────────────────────┘
       │
┌──────▼────────────────────┐
│      Infrastructure       │
└───────────────────────────┘
```

The core design rule is:

> **The GUI and CLI are clients. The Go Core Engine is the product.**

Build the architecture so that WorkBench can eventually support many runtimes, services, and tools without rewriting the foundation.
