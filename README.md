# WorkBench

**WorkBench** is a free, open-source, native, lightweight, cross-platform local development environment manager built in **pure Go**. It serves as a modern alternative to Laragon, XAMPP, and WAMP.

---

## Key Features

- **Pure Go Engine**: 100% native Go core with **zero CGO dependency** — compiles natively on Windows, macOS, and Linux without requiring GCC or external C compilers.
- **Service Management**: Easily start, stop, restart, and inspect local web services (Apache HTTP Server, MariaDB).
- **PHP Version Switching**: Install and switch seamlessly between multiple PHP versions (`8.1`, `8.2`, `8.3`, etc.). Apache automatically updates its module bindings when switching PHP versions.
- **ZIP Package Installer**: Extract and register Apache, MariaDB, and PHP runtimes directly from local ZIP archives (`bin/<service>/<version>/`).
- **Engine-First Architecture**: Completely decoupled headless core engine with a native CLI interface (`workbench` / `workbench-cli`) and a planned Qt desktop GUI layer.

---

## Directory Structure

```text
WorkBench/
├── active/            # Active runtime symlinks/paths (e.g. active/php)
├── bin/               # Service and runtime binaries
│   ├── apache/        # Installed Apache versions (e.g. bin/apache/2.4.63)
│   ├── mariadb/       # Installed MariaDB versions (e.g. bin/mariadb/11.7.2)
│   └── php/           # Installed PHP versions (e.g. bin/php/8.3.30)
├── build/             # Output executables (workbench.exe, workbench-cli.exe)
├── cache/             # Download & ZIP extraction cache
├── cmd/               # Binary entry points
│   ├── workbench/     # Main executable
│   └── workbench-cli/ # CLI executable alias
├── data/              # Persistent service data (MariaDB data files)
├── docs/              # Specifications, ADR decisions, CLI docs, & architecture
├── etc/               # Config files (etc/apache/httpd.conf, etc/php/php.ini)
├── internal/          # Core engine packages
│   ├── app/           # Application state coordinator
│   ├── cli/           # CLI command handlers
│   ├── config/        # Configuration manager & writer
│   ├── core/          # Service & runtime managers
│   ├── database/      # Embedded SQLite database
│   ├── filesystem/    # Directory path manager & archive extractor
│   ├── package/       # Package installer & ZIP extractor
│   ├── process/       # Cross-platform process manager
│   ├── project/       # Local web project manager
│   └── service/       # Apache and MariaDB service drivers
├── logs/              # Log files (apache/, php/, mariadb/)
└── www/               # Web root directory for your local projects
```

---

## Building from Source

### Prerequisites
- **Go 1.22+** installed on your system.

### Build Executables

```bash
# Build main binary
go build -o build/workbench.exe ./cmd/workbench

# Build CLI binary alias
go build -o build/workbench-cli.exe ./cmd/workbench-cli
```

### Run Tests

```bash
go test ./...
```

---

## CLI Command Reference & Usage

WorkBench provides a fast and intuitive command-line interface. Both `workbench` and `workbench-cli` accept the exact same commands:

### 1. Environment Status

Displays the current running status of Apache, MariaDB, and the active PHP version.

```bash
workbench status
```

**Example Output:**
```text
Apache      Running
MariaDB     Stopped
PHP         8.3.30
```

---

### 2. Service Management

Control background web services:

```bash
# Start a service
workbench start apache
workbench start mariadb

# Stop a service
workbench stop apache
workbench stop mariadb

# Restart a service
workbench restart apache
workbench restart mariadb
```

---

### 3. PHP Version Management

List, inspect, and switch active PHP runtimes:

```bash
# List all installed PHP versions
workbench php list

# Show currently active PHP version
workbench php current

# Switch active PHP version
workbench php use 8.3.30
```

*Note: Switching PHP versions automatically updates the Apache web server configuration and restarts Apache if it is currently running.*

---

### 4. ZIP Package & Runtime Installation

Install new versions of Apache, MariaDB, or PHP directly from ZIP archives:

```bash
# Install Apache from ZIP
workbench install apache C:\Downloads\httpd-2.4.63-win64-VS17.zip

# Install MariaDB from ZIP
workbench install mariadb C:\Downloads\mariadb-11.7.2-winx64.zip

# Install PHP from ZIP
workbench install php C:\Downloads\php-8.3.30-Win32-vs16-x64.zip
```

*(You can also use the `-install` flag format: `workbench -install php <path_to_zip>`)*

#### What Happens During ZIP Installation:
1. **Extraction**: Archives are safely extracted into `cache/extract/`.
2. **Directory Unwrapping**: Top-level single subfolders (e.g. `Apache24/`, `mariadb-11.7.2-winx64/`) are automatically unwrapped.
3. **Version Auto-Detection**: Version strings (e.g. `2.4.63`, `11.7.2`, `8.3.30`) are extracted from the file/folder name.
4. **File Placement**: Runtime files are moved directly into `bin/<service>/<version>/`.
5. **Initial Setup**: Auto-generates default configuration files (`php.ini`, log folders, data directories).

---

### 5. Local Project Management & Framework Scanning

Discover and manage local web projects inside `www/` or external custom directories with auto framework detection (`laravel`, `php`, `node`, `python`, `go`, `generic`):

```bash
# Scan www/ web root for local projects
workbench project scan

# List all tracked projects
workbench project list

# Register a custom project location outside www/
workbench project add C:\Projects\my-laravel-app

# Unregister a project
workbench project remove C:\Projects\my-laravel-app
```


---

## Documentation

For detailed architectural specifications and design decisions, see the [docs/](file:///i:/Project/Hobby/WorkBench/docs) directory:
- [docs/CLI.md](file:///i:/Project/Hobby/WorkBench/docs/CLI.md) — CLI reference
- [docs/ARCHITECTURE.md](file:///i:/Project/Hobby/WorkBench/docs/ARCHITECTURE.md) — Core system architecture
- [docs/SPEC.md](file:///i:/Project/Hobby/WorkBench/docs/SPEC.md) — Product specification
- [docs/DECISIONS.md](file:///i:/Project/Hobby/WorkBench/docs/DECISIONS.md) — Architecture Decision Records (ADRs)

---

## License

WorkBench is licensed under the [GPL-3.0 License](file:///i:/Project/Hobby/WorkBench/LICENSE).
