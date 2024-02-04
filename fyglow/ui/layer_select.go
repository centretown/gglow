package ui

import (
	"fmt"
	"gglow/fyglow/effectio"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewLayerSelect(effect *effectio.EffectIo) (sel *widget.Select) {
	var preset bool
	listener := binding.NewDataListener(func() {
		preset = true
		fmt.Println("LayerListener", preset)
		sel.SetOptions(effect.SummaryList())
		sel.SetSelectedIndex(effect.LayerIndex())
	})

	sel = widget.NewSelect([]string{}, func(s string) {})
	effect.AddFrameListener(listener)
	effect.AddLayerListener(listener)
	sel.OnChanged = onLayerChanged(sel, effect, &preset)
	return
}

func onLayerChanged(sel *widget.Select, effect *effectio.EffectIo, preset *bool) func(string) {
	return func(string) {
		fmt.Println("OnChanged", *preset)
		if *preset {
			*preset = false
			return
		}
		effect.SetCurrentLayer(sel.SelectedIndex())
	}
}
