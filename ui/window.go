//go:build cgo

package ui

import (
	"fyne.io/fyne/v2"
)

type Window struct {
	window fyne.Window
}

func NewWindow(window fyne.Window) *Window {
	return &Window{window: window}
}
