# SPEC.md

# DevBox — Product Specification

## Version

Phase 1

## Status

Implementation Specification

---

# 1. Product Definition

DevBox is a free and open-source, native, lightweight, cross-platform local development environment manager.

DevBox allows developers to install, manage, and use local development software from one simple application.

Phase 1 focuses on:

* Apache
* PHP
* MariaDB

DevBox must provide both:

* Native GUI
* CLI

Both interfaces must use the same Go core engine.

---

# 2. Product Goals

DevBox must be:

* Fast.
* Lightweight.
* Native.
* Simple.
* Cross-platform.
* Modular.
* Developer-focused.

The application must start quickly.

The user should not need to understand complex server configuration to use the basic features.

---

# 3. Phase 1 Scope

## Included

### Services

* Apache
* MariaDB

### Runtime

* PHP

### Environment

* DevBox directory structure.
* Service management.
* PHP version management.
* Basic package metadata.
* SQLite metadata storage.

### Interfaces

* Native Fyne GUI.
* CLI.

---

## Excluded

The following are explicitly excluded from Phase 1:

* Nginx.
* Redis.
* Python.
* Node.js.
* Mailpit.
* HeidiSQL.
* ngrok.
* Composer management.
* Virtual hosts.
* Automatic SSL.
* Laravel project creation.
* Cloud synchronization.
* User accounts.
* Telemetry.

---

# 4. First Launch Behavior

When DevBox starts for the first time:

```text id="c71f7k"
Start DevBox
      │
      ▼
Resolve DevBox Root
      │
      ▼
Initialize Environment
      │
      ▼
Initialize Database
      │
      ▼
Discover Installed Software
      │
      ▼
Build Application State
      │
      ▼
Show Dashboard
```

---

## 4.1 First Launch Requirements

DevBox must:

1. Resolve the DevBox root directory.
2. Create missing directories.
3. Preserve existing directories.
4. Initialize SQLite.
5. Run database migrations.
6. Discover installed PHP versions.
7. Detect Apache installation.
8. Detect MariaDB installation.
9. Display the dashboard.

If initialization fails, DevBox must show a clear error.

---

## 4.2 First Launch Must Not

DevBox must not:

* Delete user files.
* Delete `www/`.
* Delete `data/`.
* Delete `backup/`.
* Download software automatically without user action.
* Start Apache automatically.
* Start MariaDB automatically.

---

# 5. DevBox Environment

DevBox uses a root directory.

The root must contain:

```text id="a5i1bc"
DevBox/
│
├── bin/
├── active/
├── www/
├── data/
├── etc/
├── logs/
├── backup/
├── packages/
└── cache/
```

---

# 6. Directory Specification

## 6.1 `bin/`

Contains installed software.

```text id="6h9j3m"
bin/
├── apache/
├── php/
└── mariadb/
```

PHP example:

```text id="v0qvqx"
bin/php/
├── 8.1.29/
├── 8.2.20/
└── 8.3.30/
```

Each software version must have its own directory.

---

## 6.2 `active/`

Represents the currently active runtime.

PHP:

```text id="gq87r4"
active/php
```

The active PHP version must be uniquely identifiable.

---

## 6.3 `www/`

Default web root.

The user may place web projects here.

DevBox must not automatically delete contents of this directory.

---

## 6.4 `data/`

Contains persistent service data.

MariaDB:

```text id="5wq5a7"
data/mariadb/
```

Service data must be separated from software binaries.

---

## 6.5 `etc/`

Contains configuration.

Example:

```text id="8x21q8"
etc/
├── apache/
├── php/
└── mariadb/
```

---

## 6.6 `logs/`

Contains application and service logs.

Example:

```text id="0c7o8k"
logs/
├── devbox.log
├── apache/
└── mariadb/
```

---

## 6.7 `backup/`

Contains user or system backups.

DevBox must not automatically delete backup data.

---

## 6.8 `packages/`

Contains package-related data.

This may include:

* Package metadata.
* Installed package records.
* Package manifests.

---

## 6.9 `cache/`

Contains temporary cached data.

Cache data may be safely recreated.

---

# 7. Dashboard Specification

The main GUI screen is the **Dashboard**.

The dashboard must show:

```text id="9y1y8e"
Apache
MariaDB
PHP
```

Example:

```text id="3g4vkg"
┌──────────────────────────────────────────────┐
│ DevBox                                       │
├──────────────────────────────────────────────┤
│ Apache       ● Running        [Stop]          │
│ MariaDB      ● Stopped        [Start]         │
│ PHP          8.3.30                            │
│                                              │
│ PHP Version: [ 8.3.30 ▼ ]                     │
└──────────────────────────────────────────────┘
```

The exact visual design may evolve.

The behavior must remain consistent.

---

# 8. Service Status Specification

The following service statuses are supported:

```text id="x3qu21"
Unknown
Stopped
Starting
Running
Stopping
Error
```

---

## 8.1 Status Meaning

### Unknown

DevBox cannot determine the current service state.

### Stopped

The service is not running.

### Starting

DevBox is attempting to start the service.

### Running

The service process is running and has passed basic verification.

### Stopping

DevBox is attempting to stop the service.

### Error

The service failed to start, stop, or verify correctly.

---

# 9. Apache Specification

## 9.1 Apache Installation

Apache is installed under:

```text id="23b8a2"
bin/apache/<version>/
```

DevBox must detect whether Apache is installed.

If Apache is not installed, the UI must display:

```text id="c3y0w0"
Not Installed
```

---

## 9.2 Start Apache

When the user selects:

```text id="yr1jcg"
Start Apache
```

DevBox must:

1. Confirm Apache is installed.
2. Resolve the Apache executable.
3. Resolve the Apache configuration.
4. Start the Apache process.
5. Verify that the process remains active.
6. Update the service status.

Success:

```text id="k2y5ba"
Apache → Running
```

Failure:

```text id="e5h1a6"
Apache → Error
```

The user must receive a meaningful error.

---

## 9.3 Stop Apache

When the user selects:

```text id="0r7x3k"
Stop Apache
```

DevBox must:

1. Detect the Apache process.
2. Request process termination.
3. Wait for process termination.
4. Verify the process is no longer running.
5. Update status.

Success:

```text id="d7r3av"
Apache → Stopped
```

---

## 9.4 Restart Apache

Restart must behave as:

```text id="9z36p6"
Stop
  ↓
Verify Stopped
  ↓
Start
  ↓
Verify Running
```

If the start operation fails:

```text id="9v3lqw"
Apache → Error
```

---

# 10. MariaDB Specification

## 10.1 MariaDB Installation

MariaDB binaries are stored under:

```text id="q7l8w5"
bin/mariadb/<version>/
```

MariaDB data is stored under:

```text id="5jd7g8"
data/mariadb/
```

---

## 10.2 Start MariaDB

When the user selects:

```text id="3i5q09"
Start MariaDB
```

DevBox must:

1. Confirm MariaDB is installed.
2. Resolve the executable.
3. Resolve the data directory.
4. Resolve configuration.
5. Start MariaDB.
6. Verify the process remains active.
7. Update status.

---

## 10.3 Stop MariaDB

DevBox must:

1. Detect the MariaDB process.
2. Request termination.
3. Wait for termination.
4. Verify the process has stopped.

DevBox must not delete MariaDB data.

---

## 10.4 Restart MariaDB

Restart must follow:

```text id="z6r5k5"
Stop
  ↓
Verify Stopped
  ↓
Start
  ↓
Verify Running
```

---

# 11. PHP Specification

PHP is managed as a runtime.

PHP versions are installed independently.

Example:

```text id="1n9s9s"
8.1.29
8.2.20
8.3.30
```

---

# 12. PHP Version Discovery

DevBox must scan:

```text id="1z9n3q"
bin/php/
```

A directory is a valid PHP installation only if it passes PHP installation validation.

The application must not assume:

```text id="b7a1yf"
Every directory = PHP version
```

---

# 13. PHP Active Version

Only one PHP version may be active at a time.

Example:

```text id="v98q65"
Installed:
8.1.29
8.2.20
8.3.30

Active:
8.3.30
```

The active version must be persisted.

---

# 14. PHP Version Switching

When the user selects:

```text id="3vhrv6"
PHP 8.2.20
```

DevBox must:

1. Validate that PHP 8.2.20 exists.
2. Validate the PHP installation.
3. Prepare the switch.
4. Update active PHP state.
5. Verify the new active version.
6. Commit the change.

---

## 14.1 PHP Switch Success

The application state becomes:

```text id="8gd8pj"
Active PHP: 8.2.20
```

---

## 14.2 PHP Switch Failure

If switching fails:

```text id="x1n1y3"
Previous PHP remains active
```

The user must receive an error.

The application must not leave the active PHP state partially changed.

---

# 15. PHP and Apache Integration

Apache must use the active PHP runtime.

When the active PHP version changes:

```text id="zhv6ir"
PHP 8.3.30
     ↓
Switch
     ↓
PHP 8.2.20
```

The application must ensure Apache's PHP integration uses the active PHP version.

The exact integration implementation is platform-specific.

The user must not manually copy PHP files between versions.

---

# 16. CLI Specification

The CLI executable is:

```text id="v2z7wr"
devbox-cli
```

The CLI should provide a user-friendly command structure.

---

## 16.1 Status

Command:

```bash id="e8v8a8"
devbox status
```

Example output:

```text id="ax4x6q"
Apache     Running
MariaDB    Stopped
PHP        8.3.30
```

---

## 16.2 Start

```bash id="xx6lvr"
devbox start apache
devbox start mariadb
```

---

## 16.3 Stop

```bash id="n8q0q5"
devbox stop apache
devbox stop mariadb
```

---

## 16.4 Restart

```bash id="x5lj0c"
devbox restart apache
devbox restart mariadb
```

---

## 16.5 PHP List

```bash id="q4jz0u"
devbox php list
```

Example:

```text id="pr5dkt"
8.1.29
8.2.20
8.3.30 *
```

`*` represents the active PHP version.

---

## 16.6 PHP Current

```bash id="hrwqvr"
devbox php current
```

Example:

```text id="n5f8m0"
8.3.30
```

---

## 16.7 PHP Use

```bash id="5k4y57"
devbox php use 8.2.20
```

Success:

```text id="w5s7gt"
PHP 8.2.20 is now active.
```

Failure:

```text id="4n9z7h"
Failed to activate PHP 8.2.20.
```

---

# 17. CLI Exit Codes

The CLI must use meaningful exit codes.

```text id="j9u8b1"
0  Success
1  General error
2  Invalid command or argument
3  Service error
4  Runtime error
5  Package error
```

The exact mapping may evolve.

The CLI must never return exit code `0` for a failed operation.

---

# 18. GUI Interaction Specification

## 18.1 Start Button

When the user clicks `Start`:

```text id="z8qf1b"
Button Click
    ↓
Application Command
    ↓
Service Start
    ↓
Status Update
    ↓
UI Refresh
```

The UI must not freeze unnecessarily during long operations.

---

## 18.2 Stop Button

The UI must:

* Start the stop operation.
* Show a stopping state.
* Disable conflicting actions if necessary.
* Display the final result.

---

## 18.3 PHP Selector

The PHP selector must:

* Display installed versions.
* Mark the active version.
* Allow selection.
* Show switching progress or state.
* Display failure messages.

---

# 19. Error Handling Specification

Errors must be understandable to the developer.

Bad:

```text id="s3y4u8"
Error 123
```

Good:

```text id="j7n3af"
Failed to start Apache: Apache executable was not found.
```

---

## 19.1 Normal Errors

Normal operational errors must be shown as user-friendly messages.

Do not display stack traces.

---

## 19.2 Technical Logs

Detailed technical information may be written to:

```text id="r9g4v0"
logs/devbox.log
```

---

# 20. Data Safety Specification

DevBox must prioritize user data safety.

The following directories must never be automatically deleted:

```text id="m5q5f4"
www/
data/
backup/
```

DevBox must not:

* Delete a project during service operations.
* Delete MariaDB data during upgrades.
* Delete backups automatically.
* Overwrite a valid runtime during package installation.

---

# 21. Package Installation Specification

When package installation is implemented:

```text id="3nj5q7"
Download
   ↓
Verify Checksum
   ↓
Extract Temporary
   ↓
Validate
   ↓
Move Final
   ↓
Register
```

If any step fails:

```text id="j7v9x2"
Final Installation Remains Unchanged
```

---

# 22. Startup Performance

DevBox should start quickly.

The application must avoid:

* Unnecessary network requests.
* Full filesystem scans on every UI refresh.
* Loading huge log files.
* Blocking the UI thread with long operations.

The exact startup performance target may be defined after the first working prototype.

---

# 23. Cross-Platform Behavior

DevBox must preserve the same product behavior on:

```text id="j4o4x3"
Windows
Linux
macOS
```

Platform-specific implementation details may differ.

The user-facing concepts must remain consistent.

Example:

```text id="5y1l5e"
Start Apache
```

must mean the same thing on all supported platforms.

---

# 24. Phase 1 Definition of Done

Phase 1 is complete when a developer can:

1. Start DevBox.
2. See the environment dashboard.
3. See Apache status.
4. See MariaDB status.
5. See the active PHP version.
6. Start Apache.
7. Stop Apache.
8. Restart Apache.
9. Start MariaDB.
10. Stop MariaDB.
11. Restart MariaDB.
12. See installed PHP versions.
13. Switch PHP versions.
14. Use the CLI for the same core operations.
15. Use the GUI for the same core operations.
16. Run the project tests successfully.

---

# 25. Product Principle

The most important DevBox product principle is:

> **A developer should be able to manage their local development environment without fighting the environment manager.**

DevBox should make local development infrastructure:

```text id="o5c4x3"
Simple
Fast
Predictable
Native
```

The implementation must always prioritize this experience.
