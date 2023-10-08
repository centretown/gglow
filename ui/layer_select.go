package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewLayerSelect(model *data.Model) (sel *widget.Select) {
	sel = widget.NewSelect([]string{}, func(s string) {})
	sel.PlaceHolder = "choose layer..."
	sel.Alignment = fyne.TextAlignCenter
	sel.OnChanged = func(s string) {
		index := sel.SelectedIndex()
		if index >= 0 && index < len(sel.Options) {
			model.SetCurrentLayer(index)
		}
	}
	model.LayerSummaryList.AddListener(binding.NewDataListener(func() {
		summaries, _ := model.LayerSummaryList.Get()
		sel.SetOptions(summaries)
		sel.SetSelectedIndex(0)
	}))
	model.LayerIndex.AddListener(binding.NewDataListener(func() {
		index, _ := model.LayerIndex.Get()
		sel.SetSelectedIndex(index)
	}))
	return
}
