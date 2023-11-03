package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SelectItem struct {
	*widget.Select
}

func NewSelectItem(sel *widget.Select) (si *SelectItem) {
	si = &SelectItem{
		Select: sel,
	}
	return
}

func (si *SelectItem) ToolbarObject() fyne.CanvasObject {
	return si.Select
}
