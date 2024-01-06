package fyui

import (
	"gglow/fyio"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LayerTools struct {
	*widget.Toolbar
}

func NewLayerTools(effect *fyio.EffectIo) *LayerTools {
	lt := &LayerTools{
		Toolbar: widget.NewToolbar(),
	}

	lt.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})))

	lt.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {})))
	return lt
}
