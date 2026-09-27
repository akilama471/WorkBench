---
description: Build commands and Go workflow rules for WorkBench project
globs: *
---

# WorkBench Rules

When working on this project, always adhere to the following rules:

1. **Build Location**: Always build executables into the `build/` directory.
2. **GUI Build Command**: Build `cmd/workbench-gui` as `build/workbench.exe` using the following command:
   ```bash
   go build -o build/workbench.exe ./cmd/workbench-gui
   ```
3. **CLI Build Command**: Build `cmd/workbench-cli` as `build/workbench-cli.exe` using the following command:
   ```bash
   go build -o build/workbench-cli.exe ./cmd/workbench-cli
   ```
4. **Golang Skills**: Always follow and utilize the Go language skills and conventions located in the `.agents/skills` directory.
5. **Git Commits**: After completing every successful step or feature, make a Git commit with a meaningful commit message.
