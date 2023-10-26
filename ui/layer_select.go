package ui

import (
	"glow-gui/data"
	"glow-gui/resources"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewLayerSelect(model *data.Model) (sel *widget.Select) {
	sel = widget.NewSelect([]string{}, func(s string) {})
	sel.PlaceHolder = resources.LayersLabel.PlaceHolder() + "..."
	// sel.Alignment = fyne.TextAlignCenter
	sel.OnChanged = func(s string) {
		model.SetCurrentLayer(sel.SelectedIndex())
	}
	model.LayerSummaryList.AddListener(binding.NewDataListener(func() {
		summaries, _ := model.LayerSummaryList.Get()
		sel.SetOptions(summaries)
		sel.SetSelectedIndex(0)
	}))
	return
}
