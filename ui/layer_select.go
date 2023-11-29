package ui

import (
	"glow-gui/control"
	"glow-gui/resources"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewLayerSelect(model *control.Model) (sel *widget.Select) {
	sel = widget.NewSelect([]string{}, func(s string) {})
	sel.PlaceHolder = resources.LayersLabel.PlaceHolder() + "..."
	// sel.Alignment = fyne.TextAlignCenter
	sel.OnChanged = func(s string) {
		model.SetCurrentLayer(sel.SelectedIndex())
	}
	model.AddLayerListener(binding.NewDataListener(func() {
		summaries, _ := model.LayerSummaryList.Get()
		sel.SetOptions(summaries)
		sel.SetSelectedIndex(model.LayerIndex())
	}))
	return
}
