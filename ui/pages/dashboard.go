//go:build cgo

package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	apppkg "github.com/akilama471/WorkBench/internal/app"
)

func BuildDashboard(app *apppkg.Application) fyne.CanvasObject {
	title := widget.NewLabel("WorkBench")
	title.TextStyle = fyne.TextStyle{Bold: true}

	apacheStatus := widget.NewLabel("Apache: Unknown")
	mariadbStatus := widget.NewLabel("MariaDB: Unknown")

	status := container.NewVBox(title, apacheStatus, mariadbStatus)

	return status
}
