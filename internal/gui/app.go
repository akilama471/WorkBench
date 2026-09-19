package gui

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	btnStartMysql  widget.Clickable
	btnStopMysql   widget.Clickable

	// Status Cache
	apacheStatus string
	apacheVer    string
	mariaStatus  string
	mariaVer     string
	mysqlStatus  string
	mysqlVer     string
	phpStatus    string

	// Top Right Button
	btnOptions widget.Clickable
}

func NewUI(b *backend.Application) *UI {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	
	ui := &UI{
		backend: b,
		theme:   th,
		window:  new(app.Window),
	}
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
			st := svc.Status()
			if string(st) == "Running" {
				ui.apacheStatus = fmt.Sprintf("Running (Port %d)", svc.Port())
			} else {
				ui.apacheStatus = string(st)
			}
		}
	} else {
		ui.apacheStatus = "Unknown"
	}
	ui.apacheVer, _ = ui.backend.CurrentServiceVersion("apache")

	if svc, err := ui.backend.ServiceManager.GetService("mariadb"); err == nil {
		if !svc.IsInstalled() {
			ui.mariaStatus = "Not Installed"
		} else {
			st := svc.Status()
			if string(st) == "Running" {
				ui.mariaStatus = fmt.Sprintf("Running (Port %d)", svc.Port())
			} else {
				ui.mariaStatus = string(st)
			}
		}
	} else {
		ui.mariaStatus = "Unknown"
	}
	ui.mariaVer, _ = ui.backend.CurrentServiceVersion("mariadb")

	if svc, err := ui.backend.ServiceManager.GetService("mysql"); err == nil {
		if !svc.IsInstalled() {
			ui.mysqlStatus = "Not Installed"
		} else {
			st := svc.Status()
			if string(st) == "Running" {
				ui.mysqlStatus = fmt.Sprintf("Running (Port %d)", svc.Port())
			} else {
				ui.mysqlStatus = string(st)
			}
		}
	} else {
		ui.mysqlStatus = "Unknown"
	}
	ui.mysqlVer, _ = ui.backend.CurrentServiceVersion("mysql")

	if php, err := ui.backend.CurrentServiceVersion("php"); err == nil && php != "" {
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
				ui.window.Invalidate()
			}
			if ui.btnStopApache.Clicked(gtx) {
				ui.backend.StopService("apache")
				ui.refreshStatus()
				ui.window.Invalidate()
			}
			if ui.btnStartMaria.Clicked(gtx) {
				ui.backend.StartService("mariadb")
				ui.refreshStatus()
				ui.window.Invalidate()
			}
			if ui.btnStopMaria.Clicked(gtx) {
				ui.backend.StopService("mariadb")
				ui.refreshStatus()
				ui.window.Invalidate()
			}
			if ui.btnStartMysql.Clicked(gtx) {
				ui.backend.StartService("mysql")
				ui.refreshStatus()
				ui.window.Invalidate()
			}
			if ui.btnStopMysql.Clicked(gtx) {
				ui.backend.StopService("mysql")
				ui.refreshStatus()
				ui.window.Invalidate()
			}
			if ui.btnOptions.Clicked(gtx) {
				ui.openOptionsWindow()
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
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(material.H4(ui.theme, "WorkBench Dashboard").Layout),
					layout.Rigid(material.Button(ui.theme, &ui.btnOptions, "Options").Layout),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.serviceRow(gtx, "Apache HTTP Server", ui.apacheStatus, ui.apacheVer, &ui.btnStartApache, &ui.btnStopApache)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.serviceRow(gtx, "MariaDB Database", ui.mariaStatus, ui.mariaVer, &ui.btnStartMaria, &ui.btnStopMaria)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.serviceRow(gtx, "MySQL Database", ui.mysqlStatus, ui.mysqlVer, &ui.btnStartMysql, &ui.btnStopMysql)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				phpText := material.Body1(ui.theme, fmt.Sprintf("Active PHP Version: %s", ui.phpStatus))
				return phpText.Layout(gtx)
			}),
			layout.Flexed(1, layout.Spacer{}.Layout), // Push everything up
		)
	})
}

func (ui *UI) serviceRow(gtx layout.Context, name, status, version string, startBtn, stopBtn *widget.Clickable) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(material.Body1(ui.theme, name).Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if version != "" && version != "none" {
						return material.Body2(ui.theme, "Status: "+status+" (v"+version+")").Layout(gtx)
					}
					return material.Body2(ui.theme, "Status: "+status).Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(ui.theme, startBtn, "Start")
					if strings.HasPrefix(status, "Running") {
						gtx = gtx.Disabled()
					}
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(ui.theme, stopBtn, "Stop")
					if !strings.HasPrefix(status, "Running") {
						gtx = gtx.Disabled()
					}
					return btn.Layout(gtx)
				}),
			)
		}),
	)
}
