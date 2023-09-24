package ui

import (
	"glow-gui/res"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ToolbarSelect struct {
	Chooser *widget.Select
	icon    *widget.Icon
}

func (t *ToolbarSelect) ToolbarObject() fyne.CanvasObject {
	t.icon = widget.NewIcon(theme.DocumentIcon())
	hbox := container.NewHBox(t.icon, t.Chooser)
	return hbox
}

func NewToolbarSelect(ui *Ui) (t *ToolbarSelect) {
	t = &ToolbarSelect{}
	options := store.LookUpList()
	t.Chooser = widget.NewSelect(options, ui.OnChangeFrame)
	t.Chooser.PlaceHolder = res.ChooseEffectLabel.PlaceHolder()
	return t
}
