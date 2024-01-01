package fyui

import (
	"gglow/iohandler"
	"gglow/resources"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewLayerSelect(effect iohandler.EffectIoHandler) (sel *widget.Select) {
	sel = widget.NewSelect([]string{}, func(s string) {})
	sel.PlaceHolder = resources.LayersLabel.PlaceHolder() + "..."
	// sel.Alignment = fyne.TextAlignCenter
	sel.OnChanged = func(s string) {
		effect.SetCurrentLayer(sel.SelectedIndex())
	}
	effect.AddLayerListener(binding.NewDataListener(func() {
		summaries := effect.SummaryList()
		sel.SetOptions(summaries)
		sel.SetSelectedIndex(effect.LayerIndex())
	}))
	return
}
