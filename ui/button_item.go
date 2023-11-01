package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ButtonItem struct {
	*widget.Button
}

func NewButtonItem(button *widget.Button) (bi *ButtonItem) {
	button.Importance = widget.LowImportance
	bi = &ButtonItem{
		Button: button,
	}
	return
}

func (bi *ButtonItem) ToolbarObject() fyne.CanvasObject {
	return bi.Button
}
