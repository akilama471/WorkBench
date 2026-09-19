package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/akilama471/WorkBench/internal/app"
	"github.com/akilama471/WorkBench/internal/gui"
)

func main() {
	rootDir := resolveRootDir()

	application, err := app.New(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing WorkBench backend: %v\n", err)
		os.Exit(1)
	}
	defer application.Close()

	if err := application.Startup(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during WorkBench startup: %v\n", err)
		os.Exit(1)
	}

	ui := gui.NewUI(application)
	ui.Run()
}

func resolveRootDir() string {
	if envRoot := os.Getenv("WORKBENCH_ROOT"); envRoot != "" {
		return envRoot
	}

	exe, err := os.Executable()
	if err == nil {
		return filepath.Dir(exe)
	}

	return "."
}
