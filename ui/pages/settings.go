//go:build cgo

package pages

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	apppkg "github.com/akilama471/WorkBench/internal/app"
)

func BuildSettingsPage(app *apppkg.Application) fyne.CanvasObject {
	title := widget.NewLabel("Settings")
	title.TextStyle = fyne.TextStyle{Bold: true}

	rootDirLabel := widget.NewLabel(fmt.Sprintf("WorkBench Root: %s", app.Paths().Root()))

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		rootDirLabel,
	)
}
