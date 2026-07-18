//go:build cgo

package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func StatusBadge(status string) fyne.CanvasObject {
	return widget.NewLabel(status)
}
