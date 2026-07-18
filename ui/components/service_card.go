//go:build cgo

package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/akilama471/WorkBench/internal/service"
)

func ServiceCard(name string, status service.Status, onStart, onStop, onRestart func()) fyne.CanvasObject {
	statusLabel := widget.NewLabel(status.String())

	startBtn := widget.NewButton("Start", onStart)
	stopBtn := widget.NewButton("Stop", onStop)
	restartBtn := widget.NewButton("Restart", onRestart)

	actions := container.NewHBox(startBtn, stopBtn, restartBtn)

	return container.NewVBox(
		widget.NewLabel(name),
		statusLabel,
		actions,
	)
}
