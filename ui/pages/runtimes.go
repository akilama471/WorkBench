//go:build cgo

package pages

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	apppkg "github.com/akilama471/WorkBench/internal/app"
)

func BuildRuntimesPage(app *apppkg.Application) fyne.CanvasObject {
	title := widget.NewLabel("Runtimes")
	title.TextStyle = fyne.TextStyle{Bold: true}

	phpTitle := widget.NewLabel("PHP")
	phpTitle.TextStyle = fyne.TextStyle{Bold: true}

	activeVersion, _ := app.CurrentPHPVersion()
	activeLabel := widget.NewLabel(fmt.Sprintf("Active: %s", activeVersion))

	versions, _ := app.ListPHPVersions()
	versionSelect := widget.NewSelect(versions, func(selected string) {
		if err := app.SwitchPHPVersion(selected); err != nil {
			activeLabel.SetText(fmt.Sprintf("Error: %v", err))
		} else {
			activeLabel.SetText(fmt.Sprintf("Active: %s", selected))
		}
	})
	if activeVersion != "" && activeVersion != "none" {
		versionSelect.SetSelected(activeVersion)
	}

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		phpTitle,
		activeLabel,
		versionSelect,
	)
}
