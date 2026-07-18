//go:build cgo

package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	apppkg "github.com/akilama471/WorkBench/internal/app"
)

func BuildSettingsPage(app *apppkg.Application) fyne.CanvasObject {
	return widget.NewLabel("Settings")
}
