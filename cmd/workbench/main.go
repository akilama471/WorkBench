//go:build cgo

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	apppkg "github.com/akilama471/WorkBench/internal/app"
)

func main() {
	rootDir := resolveRootDir()

	application, err := apppkg.New(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer application.Close()

	if err := application.Startup(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fyneApp := app.New()
	window := fyneApp.NewWindow("DevBox")
	window.Resize(fyne.NewSize(600, 400))

	content := buildDashboard(application)
	window.SetContent(content)

	window.ShowAndRun()
}

func buildDashboard(application *apppkg.Application) fyne.CanvasObject {
	title := widget.NewLabel("DevBox Dashboard")
	title.TextStyle = fyne.TextStyle{Bold: true}

	apacheStatus := widget.NewLabel("Apache   -- Unknown")
	mariadbStatus := widget.NewLabel("MariaDB  -- Unknown")
	phpVersion := widget.NewLabel("PHP      -- Unknown")

	refreshBtn := widget.NewButton("Refresh", func() {
		apacheSvc, err := application.ServiceManager.GetService("apache")
		if err == nil {
			apacheStatus.SetText(fmt.Sprintf("Apache   -- %s", apacheSvc.Status()))
		}
		mariadbSvc, err := application.ServiceManager.GetService("mariadb")
		if err == nil {
			mariadbStatus.SetText(fmt.Sprintf("MariaDB  -- %s", mariadbSvc.Status()))
		}
		pv, err := application.CurrentPHPVersion()
		if err == nil && pv != "" {
			phpVersion.SetText(fmt.Sprintf("PHP      -- %s", pv))
		} else {
			phpVersion.SetText("PHP      -- none")
		}
	})

	startApacheBtn := widget.NewButton("Start Apache", func() {
		if err := application.StartService("apache"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		refreshBtn.OnTapped()
	})
	stopApacheBtn := widget.NewButton("Stop Apache", func() {
		if err := application.StopService("apache"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		refreshBtn.OnTapped()
	})
	restartApacheBtn := widget.NewButton("Restart Apache", func() {
		if err := application.RestartService("apache"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		refreshBtn.OnTapped()
	})

	startMariaBtn := widget.NewButton("Start MariaDB", func() {
		if err := application.StartService("mariadb"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		refreshBtn.OnTapped()
	})
	stopMariaBtn := widget.NewButton("Stop MariaDB", func() {
		if err := application.StopService("mariadb"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		refreshBtn.OnTapped()
	})
	restartMariaBtn := widget.NewButton("Restart MariaDB", func() {
		if err := application.RestartService("mariadb"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		refreshBtn.OnTapped()
	})

	phpVersionSelect := widget.NewSelect([]string{}, func(value string) {
		if err := application.SwitchPHPVersion(value); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		refreshBtn.OnTapped()
	})

	versions, err := application.ListPHPVersions()
	if err == nil {
		phpVersionSelect.Options = versions
	}

	refreshBtn.OnTapped()

	apacheActions := container.NewHBox(startApacheBtn, stopApacheBtn, restartApacheBtn)
	mariadbActions := container.NewHBox(startMariaBtn, stopMariaBtn, restartMariaBtn)

	apacheSection := container.NewVBox(
		widget.NewLabel("Apache"),
		apacheStatus,
		apacheActions,
	)

	mariadbSection := container.NewVBox(
		widget.NewLabel("MariaDB"),
		mariadbStatus,
		mariadbActions,
	)

	phpSection := container.NewVBox(
		widget.NewLabel("PHP"),
		phpVersion,
		phpVersionSelect,
	)

	return container.NewVBox(
		title,
		layout.NewSpacer(),
		apacheSection,
		mariadbSection,
		phpSection,
		layout.NewSpacer(),
		refreshBtn,
	)
}

func resolveRootDir() string {
	if envRoot := os.Getenv("DEVBOX_ROOT"); envRoot != "" {
		return envRoot
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, "DevBox")
}
