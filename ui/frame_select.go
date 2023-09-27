package ui

import (
	"glow-gui/res"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewFrameSelect(ui *Ui) (sel *widget.Select) {
	options := store.LookUpList()
	sel = widget.NewSelect(options, ui.OnChangeFrame)
	sel.PlaceHolder = res.ChooseEffectLabel.PlaceHolder()
	sel.Alignment = fyne.TextAlignCenter
	return
}
