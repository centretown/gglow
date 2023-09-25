package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ToolbarSelect struct {
	Chooser *widget.Select
	icon    *widget.Icon
}

func (t *ToolbarSelect) ToolbarObject() fyne.CanvasObject {
	hbox := container.NewHBox(t.icon, t.Chooser)
	return hbox
}

func NewToolbarSelect(ui *Ui) (t *ToolbarSelect) {
	t = &ToolbarSelect{}
	t.Chooser = NewFrameSelect(ui)
	t.icon = ui.frameIcon
	return t
}
