package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akilama471/WorkBench/internal/app"
)

func main() {
	rootDir := resolveRootDir()

	application, err := app.New(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer application.Close()

	if err := application.Startup(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "status":
		handleStatus(application)
	case "start":
		handleServiceCommand(application, "start", cmdArgs)
	case "stop":
		handleServiceCommand(application, "stop", cmdArgs)
	case "restart":
		handleServiceCommand(application, "restart", cmdArgs)
	case "php":
		handlePHPCommand(application, cmdArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
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

func printUsage() {
	fmt.Println("WorkBench - Local Development Environment Manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  workbench status                     Show status of all services")
	fmt.Println("  workbench start <service>            Start a service (apache, mariadb)")
	fmt.Println("  workbench stop <service>             Stop a service (apache, mariadb)")
	fmt.Println("  workbench restart <service>          Restart a service (apache, mariadb)")
	fmt.Println("  workbench php list                   List installed PHP versions")
	fmt.Println("  workbench php current                Show active PHP version")
	fmt.Println("  workbench php use <version>          Switch active PHP version")
	fmt.Println()
	fmt.Println("Services: apache, mariadb")
}

func handleStatus(application *app.Application) {
	services := []string{"apache", "mariadb"}

	for _, svcID := range services {
		svc, err := application.ServiceManager.GetService(svcID)
		if err != nil {
			fmt.Printf("%-10s  Unknown\n", capitalise(svcID))
			continue
		}
		if !svc.IsInstalled() {
			fmt.Printf("%-10s  Not Installed\n", svc.Name())
		} else {
			fmt.Printf("%-10s  %s\n", svc.Name(), svc.Status())
		}
	}

	phpVersion, err := application.CurrentPHPVersion()
	if err != nil {
		phpVersion = "none"
	}
	if phpVersion == "" {
		phpVersion = "none"
	}
	fmt.Printf("%-10s  %s\n", "PHP", phpVersion)
}

func handleServiceCommand(application *app.Application, action string, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: workbench %s <service>\n", action)
		os.Exit(1)
	}

	serviceID := strings.ToLower(args[0])

	switch action {
	case "start":
		if err := application.StartService(serviceID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s started successfully.\n", capitalise(serviceID))
	case "stop":
		if err := application.StopService(serviceID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s stopped successfully.\n", capitalise(serviceID))
	case "restart":
		if err := application.RestartService(serviceID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s restarted successfully.\n", capitalise(serviceID))
	}
}

func handlePHPCommand(application *app.Application, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: workbench php <list|current|use> [version]\n")
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		versions, err := application.ListPHPVersions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if len(versions) == 0 {
			fmt.Println("No PHP versions installed.")
			return
		}
		active, _ := application.CurrentPHPVersion()
		for _, v := range versions {
			if v == active {
				fmt.Printf("  * %s (active)\n", v)
			} else {
				fmt.Printf("    %s\n", v)
			}
		}
	case "current":
		version, err := application.CurrentPHPVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if version == "" || version == "none" {
			fmt.Println("No active PHP version.")
		} else {
			fmt.Println(version)
		}
	case "use":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: workbench php use <version>\n")
			os.Exit(1)
		}
		if err := application.SwitchPHPVersion(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("PHP version switched to %s.\n", args[1])
	default:
		fmt.Fprintf(os.Stderr, "Unknown PHP command: %s\n", args[0])
		fmt.Fprintf(os.Stderr, "Usage: workbench php <list|current|use> [version]\n")
		os.Exit(1)
	}
}

func capitalise(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
