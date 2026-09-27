package gui

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/akilama471/WorkBench/internal/winsvc"
	"github.com/sqweek/dialog"
)

type ServiceVersionControl struct {
	Versions []string
	Sel      string
	BtnSet   widget.Clickable
	Msg      string
	Loaded   bool

	DropdownOpen bool
	BtnToggle    widget.Clickable
	ItemBtns     []widget.Clickable
}

type OptionsUI struct {
	ui *UI

	// Tabs
	tabIndex   int
	btnImport  widget.Clickable
	btnGeneral widget.Clickable

	// Import Tab
	installZipPath widget.Editor
	btnBrowse      widget.Clickable
	installType    widget.Enum
	btnInstall     widget.Clickable
	installMsg     string
	isInstalling   bool

	// General Tab
	controls    map[string]*ServiceVersionControl
	generalList widget.List

	// Schedule Tab
	btnSchedule      widget.Clickable
	autoStartEnabled widget.Bool
	autoStartApache  widget.Bool
	autoStartMariaDB widget.Bool
	autoStartMySQL   widget.Bool
	autoStartDelay   widget.Editor
	btnApplySchedule widget.Clickable
	scheduleMsg      string
}

func (ui *UI) openOptionsWindow() {
	w := new(app.Window)
	w.Option(app.Title("Options - WorkBench"), app.Size(unit.Dp(500), unit.Dp(400)))

	opt := &OptionsUI{
		ui:       ui,
		tabIndex: 0,
		controls: map[string]*ServiceVersionControl{
			"apache":  {},
			"mariadb": {},
			"mysql":   {},
			"php":     {},
		},
	}
	opt.generalList.Axis = layout.Vertical
	opt.installZipPath.SingleLine = true
	opt.installType.Value = "apache"
	opt.autoStartEnabled.Value = IsAutoStartEnabled()

	// Read existing service config if enabled
	delayStr := "10"
	if opt.autoStartEnabled.Value {
		delayMs, services, err := winsvc.GetServiceConfig()
		if err == nil {
			delayStr = fmt.Sprintf("%d", delayMs/1000)
			for _, s := range services {
				if s == "apache" {
					opt.autoStartApache.Value = true
				} else if s == "mariadb" {
					opt.autoStartMariaDB.Value = true
				} else if s == "mysql" {
					opt.autoStartMySQL.Value = true
				}
			}
		}
	} else {
		// Default selections
		opt.autoStartApache.Value = true
		opt.autoStartMariaDB.Value = true
	}

	opt.autoStartDelay.SingleLine = true
	opt.autoStartDelay.SetText(delayStr)

	go func() {
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return // close only this window
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// Handle events
				if opt.btnImport.Clicked(gtx) {
					opt.tabIndex = 0
				}
				if opt.btnGeneral.Clicked(gtx) {
					opt.tabIndex = 1
				}
				if opt.btnSchedule.Clicked(gtx) {
					opt.tabIndex = 2
				}

				if opt.btnApplySchedule.Clicked(gtx) {
					if opt.autoStartEnabled.Value {
						var selectedServices []string
						if opt.autoStartApache.Value {
							selectedServices = append(selectedServices, "apache")
						}
						if opt.autoStartMariaDB.Value {
							selectedServices = append(selectedServices, "mariadb")
						}
						if opt.autoStartMySQL.Value {
							selectedServices = append(selectedServices, "mysql")
						}

						err := EnableAutoStart(opt.autoStartDelay.Text(), selectedServices)
						if err != nil {
							opt.scheduleMsg = fmt.Sprintf("Error: %v", err)
						} else {
							opt.scheduleMsg = "Auto Start Enabled successfully."
						}
					} else {
						err := DisableAutoStart()
						if err != nil {
							opt.scheduleMsg = fmt.Sprintf("Error: %v", err)
						} else {
							opt.scheduleMsg = "Auto Start Disabled successfully."
						}
					}
				}

				if opt.btnBrowse.Clicked(gtx) {
					go func() {
						filename, err := dialog.File().Filter("ZIP Files", "zip").Title("Select Package ZIP").Load()
						if err == nil && filename != "" {
							opt.installZipPath.SetText(filename)
							w.Invalidate()
						}
					}()
				}

				if opt.btnInstall.Clicked(gtx) && !opt.isInstalling {
					zipPath := opt.installZipPath.Text()
					serviceType := opt.installType.Value
					if zipPath != "" {
						opt.installMsg = "Installing..."
						opt.isInstalling = true
						// trigger frame to show message before blocking/starting goroutine
						e.Frame(gtx.Ops)
						go func(sType, p string) {
							version, err := ui.backend.InstallPackage(sType, p)
							if err != nil {
								opt.installMsg = fmt.Sprintf("Error: %v", err)
							} else {
								opt.installMsg = fmt.Sprintf("Success! Installed %s %s", sType, version)
								ui.refreshStatus()
								ui.window.Invalidate() // update main window too
							}
							opt.isInstalling = false
							w.Invalidate()
						}(serviceType, zipPath)
						continue
					} else {
						opt.installMsg = "Please enter a valid ZIP path."
					}
				}

				for svcID, ctrl := range opt.controls {
					if ctrl.BtnSet.Clicked(gtx) {
						ver := ctrl.Sel
						if ver != "" && ver != "none" {
							err := ui.backend.SwitchServiceVersion(svcID, ver)
							if err != nil {
								ctrl.Msg = fmt.Sprintf("Error: %v", err)
							} else {
								ctrl.Msg = fmt.Sprintf("Default set to %s", ver)
								ui.refreshStatus()
								ui.window.Invalidate()
							}
						}
					}
				}

				opt.layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()
}

func (opt *OptionsUI) layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(gtx,
					layout.Rigid(material.Button(opt.ui.theme, &opt.btnImport, "Import").Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(material.Button(opt.ui.theme, &opt.btnGeneral, "General").Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(material.Button(opt.ui.theme, &opt.btnSchedule, "Schedule").Layout),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if opt.tabIndex == 0 {
					return opt.layoutImportTab(gtx)
				} else if opt.tabIndex == 1 {
					return opt.layoutGeneralTab(gtx)
				}
				return opt.layoutScheduleTab(gtx)
			}),
		)
	})
}

func (opt *OptionsUI) layoutImportTab(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H5(opt.ui.theme, "Import ZIP Package")
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(material.RadioButton(opt.ui.theme, &opt.installType, "apache", "Apache").Layout),
				layout.Rigid(material.RadioButton(opt.ui.theme, &opt.installType, "mariadb", "MariaDB").Layout),
				layout.Rigid(material.RadioButton(opt.ui.theme, &opt.installType, "mysql", "MySQL").Layout),
				layout.Rigid(material.RadioButton(opt.ui.theme, &opt.installType, "php", "PHP").Layout),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(opt.ui.theme, &opt.installZipPath, "Absolute Path to ZIP File")
					return ed.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(material.Button(opt.ui.theme, &opt.btnBrowse, "Browse").Layout),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(opt.ui.theme, &opt.btnInstall, "Install from ZIP")
					if opt.isInstalling {
						gtx = gtx.Disabled()
					}
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if opt.isInstalling {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(material.Loader(opt.ui.theme).Layout),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							layout.Rigid(material.Body2(opt.ui.theme, opt.installMsg).Layout),
						)
					}
					return material.Body2(opt.ui.theme, opt.installMsg).Layout(gtx)
				}),
			)
		}),
	)
}

func (opt *OptionsUI) layoutGeneralTab(gtx layout.Context) layout.Dimensions {
	services := []struct {
		ID   string
		Name string
	}{
		{"apache", "Apache"},
		{"mariadb", "MariaDB"},
		{"mysql", "MySQL"},
		{"php", "PHP"},
	}

	for _, s := range services {
		ctrl := opt.controls[s.ID]
		if !ctrl.Loaded {
			versions, _ := opt.ui.backend.ListServiceVersions(s.ID)
			if len(versions) == 0 {
				ctrl.Versions = []string{"none"}
			} else {
				ctrl.Versions = versions
			}
			current, _ := opt.ui.backend.CurrentServiceVersion(s.ID)
			if current != "" {
				ctrl.Sel = current
			} else {
				ctrl.Sel = "none"
			}
			ctrl.ItemBtns = make([]widget.Clickable, len(ctrl.Versions))
			ctrl.Loaded = true
		}

		if ctrl.BtnToggle.Clicked(gtx) {
			ctrl.DropdownOpen = !ctrl.DropdownOpen
		}
		for i := range ctrl.Versions {
			if ctrl.ItemBtns[i].Clicked(gtx) {
				ctrl.Sel = ctrl.Versions[i]
				ctrl.DropdownOpen = false
			}
		}
	}

	var allChildren []layout.Widget

	for _, s := range services {
		svc := s
		ctrl := opt.controls[svc.ID]

		allChildren = append(allChildren, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H6(opt.ui.theme, "Default "+svc.Name+" Version").Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Main Dropdown button
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(opt.ui.theme, &ctrl.BtnToggle, ctrl.Sel+" ▼")
							btn.Background = opt.ui.theme.Palette.ContrastBg
							return btn.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if !ctrl.DropdownOpen {
								return layout.Dimensions{}
							}
							var list []layout.FlexChild
							for i, v := range ctrl.Versions {
								idx := i
								ver := v
								list = append(list, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									itemBtn := material.Button(opt.ui.theme, &ctrl.ItemBtns[idx], ver)
									itemBtn.Background = opt.ui.theme.Palette.Bg
									itemBtn.Color = opt.ui.theme.Palette.Fg
									return itemBtn.Layout(gtx)
								}))
							}
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx, list...)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(material.Button(opt.ui.theme, &ctrl.BtnSet, "Set as Default").Layout),
						layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
						layout.Rigid(material.Body2(opt.ui.theme, ctrl.Msg).Layout),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
			)
		})
	}

	lst := material.List(opt.ui.theme, &opt.generalList)
	return lst.Layout(gtx, len(allChildren), func(gtx layout.Context, index int) layout.Dimensions {
		return allChildren[index](gtx)
	})
}

func (opt *OptionsUI) layoutScheduleTab(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H5(opt.ui.theme, "Windows Service Auto Start")
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(material.CheckBox(opt.ui.theme, &opt.autoStartEnabled, "Enable WorkBench Windows Service").Layout),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !opt.autoStartEnabled.Value {
				gtx = gtx.Disabled()
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(material.Body2(opt.ui.theme, "Select services to auto-start:").Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(material.CheckBox(opt.ui.theme, &opt.autoStartApache, "Apache HTTP Server").Layout),
				layout.Rigid(material.CheckBox(opt.ui.theme, &opt.autoStartMariaDB, "MariaDB").Layout),
				layout.Rigid(material.CheckBox(opt.ui.theme, &opt.autoStartMySQL, "MySQL").Layout),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !opt.autoStartEnabled.Value {
				gtx = gtx.Disabled()
			}
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(material.Body1(opt.ui.theme, "Delay before start (seconds): ").Layout),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(opt.ui.theme, &opt.autoStartDelay, "10")
					return ed.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(material.Button(opt.ui.theme, &opt.btnApplySchedule, "Apply").Layout),
				layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
				layout.Rigid(material.Body2(opt.ui.theme, opt.scheduleMsg).Layout),
			)
		}),
	)
}
