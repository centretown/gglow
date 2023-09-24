package ui

import (
	"fyne.io/fyne/v2/widget"
)

type LayerItem struct {
	*widget.Label
}

func NewLayerItem() (item *LayerItem) {
	item = &LayerItem{}
	item.Label = widget.NewLabel("layer")
	return
}
