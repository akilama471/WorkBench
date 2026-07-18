//go:build cgo

package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func RuntimeSelector(versions []string, active string, onSelect func(string)) fyne.CanvasObject {
	selectWidget := widget.NewSelect(versions, onSelect)
	if active != "" {
		selectWidget.SetSelected(active)
	}
	return selectWidget
}
