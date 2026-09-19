package gui

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
    "gioui.org/unit"

	backend "github.com/akilama471/WorkBench/internal/app"
)

type UI struct {
	backend *backend.Application
	theme   *material.Theme
	window  *app.Window

	// Buttons
	btnStartApache widget.Clickable
	btnStopApache  widget.Clickable
	btnStartMaria  widget.Clickable
	btnStopMaria   widget.Clickable

	// Status Cache
	apacheStatus string
	mariaStatus  string
	phpStatus    string

	// Install Package Section
	installZipPath widget.Editor
	installType    widget.Enum
	btnInstall     widget.Clickable
	installMsg     string
}

func NewUI(b *backend.Application) *UI {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	
	ui := &UI{
		backend: b,
		theme:   th,
		window:  new(app.Window),
	}
	ui.installZipPath.SingleLine = true
	ui.installType.Value = "apache"
	return ui
}

func (ui *UI) Run() {
	ui.window.Option(app.Title("WorkBench"), app.Size(unit.Dp(600), unit.Dp(400)))
	ui.refreshStatus()

	go func() {
		if err := ui.loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (ui *UI) refreshStatus() {
	if svc, err := ui.backend.ServiceManager.GetService("apache"); err == nil {
		if !svc.IsInstalled() {
			ui.apacheStatus = "Not Installed"
		} else {
			ui.apacheStatus = string(svc.Status())
		}
	} else {
		ui.apacheStatus = "Unknown"
	}

	if svc, err := ui.backend.ServiceManager.GetService("mariadb"); err == nil {
		if !svc.IsInstalled() {
			ui.mariaStatus = "Not Installed"
		} else {
			ui.mariaStatus = string(svc.Status())
		}
	} else {
		ui.mariaStatus = "Unknown"
	}

	if php, err := ui.backend.CurrentPHPVersion(); err == nil && php != "" {
		ui.phpStatus = php
	} else {
		ui.phpStatus = "None"
	}
}

func (ui *UI) loop() error {
	var ops op.Ops
	for {
		switch e := ui.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Handle Button Clicks
			if ui.btnStartApache.Clicked(gtx) {
				ui.backend.StartService("apache")
				ui.refreshStatus()
			}
			if ui.btnStopApache.Clicked(gtx) {
				ui.backend.StopService("apache")
				ui.refreshStatus()
			}
			if ui.btnStartMaria.Clicked(gtx) {
				ui.backend.StartService("mariadb")
				ui.refreshStatus()
			}
			if ui.btnStopMaria.Clicked(gtx) {
				ui.backend.StopService("mariadb")
				ui.refreshStatus()
			}
			if ui.btnInstall.Clicked(gtx) {
				zipPath := ui.installZipPath.Text()
				serviceType := ui.installType.Value
				if zipPath != "" {
					ui.installMsg = "Installing..."
					// Update UI before starting long running task
					e.Frame(gtx.Ops) 
					go func(sType, p string) {
						version, err := ui.backend.InstallPackage(sType, p)
						if err != nil {
							ui.installMsg = fmt.Sprintf("Error: %v", err)
						} else {
							ui.installMsg = fmt.Sprintf("Success! Installed %s %s", sType, version)
							ui.refreshStatus()
						}
						ui.window.Invalidate()
					}(serviceType, zipPath)
					continue
				} else {
					ui.installMsg = "Please enter a valid ZIP path."
				}
			}

			ui.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (ui *UI) layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H4(ui.theme, "WorkBench Dashboard")
				return title.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.serviceRow(gtx, "Apache HTTP Server", ui.apacheStatus, &ui.btnStartApache, &ui.btnStopApache)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.serviceRow(gtx, "MariaDB Database", ui.mariaStatus, &ui.btnStartMaria, &ui.btnStopMaria)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				phpText := material.Body1(ui.theme, fmt.Sprintf("Active PHP Version: %s", ui.phpStatus))
				return phpText.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(32)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H5(ui.theme, "Install Package")
				return title.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(material.RadioButton(ui.theme, &ui.installType, "apache", "Apache").Layout),
					layout.Rigid(material.RadioButton(ui.theme, &ui.installType, "mariadb", "MariaDB").Layout),
					layout.Rigid(material.RadioButton(ui.theme, &ui.installType, "mysql", "MySQL").Layout),
					layout.Rigid(material.RadioButton(ui.theme, &ui.installType, "php", "PHP").Layout),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				ed := material.Editor(ui.theme, &ui.installZipPath, "Absolute Path to ZIP File")
				return ed.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(material.Button(ui.theme, &ui.btnInstall, "Install from ZIP").Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(material.Body2(ui.theme, ui.installMsg).Layout),
				)
			}),
			layout.Flexed(1, layout.Spacer{}.Layout), // Push everything up
		)
	})
}

func (ui *UI) serviceRow(gtx layout.Context, name, status string, startBtn, stopBtn *widget.Clickable) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(material.Body1(ui.theme, name).Layout),
				layout.Rigid(material.Body2(ui.theme, "Status: "+status).Layout),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(ui.theme, startBtn, "Start")
					if status == "Running" {
						gtx = gtx.Disabled()
					}
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(ui.theme, stopBtn, "Stop")
					if status != "Running" {
						gtx = gtx.Disabled()
					}
					return btn.Layout(gtx)
				}),
			)
		}),
	)
}
