package ui

import (
	"gglow/fyglow/effectio"
	"gglow/resources"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewLayerSelect(effect *effectio.EffectIo) (sel *widget.Select) {
	sel = widget.NewSelect([]string{}, func(s string) {})
	sel.PlaceHolder = resources.LayersLabel.PlaceHolder() + "..."
	sel.OnChanged = func(s string) {
		index := sel.SelectedIndex()
		effect.SetCurrentLayer(index)
	}
	effect.AddFrameListener(binding.NewDataListener(func() {
		summaries := effect.SummaryList()
		sel.SetOptions(summaries)
		sel.SetSelectedIndex(effect.LayerIndex())
	}))
	return
}
