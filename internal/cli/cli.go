package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akilama471/WorkBench/internal/app"
)

// Run executes the WorkBench command line interface.
func Run() {
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
	case "project":
		handleProjectCommand(application, cmdArgs)
	case "install", "-install":
		handleInstallCommand(application, cmdArgs)
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
		dir := filepath.Dir(exe)
		if strings.ToLower(filepath.Base(dir)) == "build" {
			return filepath.Dir(dir)
		}
		return dir
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
	fmt.Println("  workbench project list               List all tracked local projects")
	fmt.Println("  workbench project scan               Scan www/ directory for local projects")
	fmt.Println("  workbench project add <path>         Register a custom project path")
	fmt.Println("  workbench project remove <path>      Unregister a tracked project")
	fmt.Println("  workbench install <service> <zip>    Extract and install runtime package from zip")
	fmt.Println("  workbench -install <service> <zip>   Extract and install runtime package from zip")
	fmt.Println()
	fmt.Println("Services: apache, mariadb, php")
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

func handleProjectCommand(application *app.Application, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: workbench project <list|scan|add|remove> [path]\n")
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		projects, err := application.ListProjects()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if len(projects) == 0 {
			fmt.Println("No projects tracked. Run 'workbench project scan' to discover projects in www/")
			return
		}
		fmt.Printf("%-20s  %-12s  %s\n", "NAME", "TYPE", "PATH")
		fmt.Println(strings.Repeat("-", 60))
		for _, p := range projects {
			fmt.Printf("%-20s  %-12s  %s\n", p.Name, p.Type, p.Path)
		}

	case "scan":
		fmt.Println("Scanning www/ for local web projects...")
		projects, err := application.ScanProjects()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if len(projects) == 0 {
			fmt.Println("No projects found in www/ directory.")
			return
		}
		fmt.Printf("Scan complete. Discovered %d project(s):\n\n", len(projects))
		fmt.Printf("%-20s  %-12s  %s\n", "NAME", "TYPE", "PATH")
		fmt.Println(strings.Repeat("-", 60))
		for _, p := range projects {
			fmt.Printf("%-20s  %-12s  %s\n", p.Name, p.Type, p.Path)
		}

	case "add":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: workbench project add <path>\n")
			os.Exit(1)
		}
		p, err := application.AddProject(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Project '%s' (%s) registered successfully at %s\n", p.Name, p.Type, p.Path)

	case "remove":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: workbench project remove <path>\n")
			os.Exit(1)
		}
		if err := application.RemoveProject(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Project at '%s' removed from tracked list.\n", args[1])

	default:
		fmt.Fprintf(os.Stderr, "Unknown project command: %s\n", args[0])
		fmt.Fprintf(os.Stderr, "Usage: workbench project <list|scan|add|remove> [path]\n")
		os.Exit(1)
	}
}

func handleInstallCommand(application *app.Application, args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: workbench install <apache|mariadb|php> <path_to_zip>\n")
		fmt.Fprintf(os.Stderr, "   or: workbench -install <apache|mariadb|php> <path_to_zip>\n")
		os.Exit(1)
	}

	serviceType := args[0]
	zipPath := args[1]

	fmt.Printf("Extracting and installing %s from %s...\n", capitalise(serviceType), zipPath)
	version, err := application.InstallPackage(serviceType, zipPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s version %s installed successfully to bin/%s/%s.\n", capitalise(serviceType), version, strings.ToLower(serviceType), version)
}

func capitalise(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}


