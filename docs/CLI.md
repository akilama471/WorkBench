# WorkBench CLI Documentation

WorkBench provides a native, pure-Go command-line interface for managing local development environment services, runtimes, and packages.

## Commands Overview

### 1. Service Management
- **`workbench status`**  
  Displays operational status for all installed services (Apache, MariaDB, MySQL, PHP).
- **`workbench start <service>`**  
  Starts the specified service (`apache`, `mariadb`, or `mysql`).
- **`workbench stop <service>`**  
  Stops the specified service (`apache`, `mariadb`, or `mysql`).
- **`workbench restart <service>`**  
  Restarts the specified service (`apache`, `mariadb`, or `mysql`).

### 2. PHP Version Management
- **`workbench php list`**  
  Lists all installed PHP versions in `bin/php/`.
- **`workbench php current`**  
  Shows the active PHP version.
- **`workbench php use <version>`**  
  Switches the active PHP version and updates Apache configuration.

### 3. Package & ZIP Installation Commands
- **`workbench install apache <path_to_zip>`** / **`workbench -install apache <path_to_zip>`**  
  Extracts the Apache ZIP archive, auto-detects the version, moves files to `bin/apache/[version]/`, and completes initial setup.
- **`workbench install mariadb <path_to_zip>`** / **`workbench -install mariadb <path_to_zip>`**  
  Extracts the MariaDB ZIP archive, auto-detects the version, moves files to `bin/mariadb/[version]/`, and initializes data directories in `data/mariadb/`.
- **`workbench install mysql <path_to_zip>`** / **`workbench -install mysql <path_to_zip>`**  
  Extracts the MySQL ZIP archive, auto-detects the version, moves files to `bin/mysql/[version]/`, and initializes data directories in `data/mysql/`.
- **`workbench install php <path_to_zip>`** / **`workbench -install php <path_to_zip>`**  
  Extracts the PHP ZIP archive, auto-detects the version, moves files to `bin/php/[version]/`, and configures `php.ini`.

