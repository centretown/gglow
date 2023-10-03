package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewLayerSelect(model *data.Model) (sel *widget.Select) {
	sel = widget.NewSelect([]string{}, func(s string) {
	})
	sel.PlaceHolder = "choose layer..."
	sel.Alignment = fyne.TextAlignCenter
	return
}
