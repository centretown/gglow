package fyui

import (
	"gglow/fyio"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LayerTools struct {
	InsertButton *ButtonItem
	RemoveButton *ButtonItem
}

func NewLayerTools(effect *fyio.EffectIo) *LayerTools {
	lt := &LayerTools{}

	lt.InsertButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), lt.add))
	lt.RemoveButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), lt.remove))
	return lt
}

func (lt *LayerTools) Items() (items []widget.ToolbarItem) {
	items = []widget.ToolbarItem{
		lt.InsertButton,
		lt.RemoveButton,
	}
	return
}

func (lt *LayerTools) add() {
}

func (lt *LayerTools) remove() {
}
