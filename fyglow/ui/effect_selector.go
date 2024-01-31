package ui

import (
	"gglow/fyglow/effectio"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// type EffectSelect struct {
// 	*widget.Select
// 	effect  *effectio.EffectIo
// 	auto    bool
// }

func NewEffectSelector(effect *effectio.EffectIo) *widget.Select {
	var auto bool
	selector := widget.NewSelect(effect.ListCurrent(), OnEffectChanged(effect, &auto))

	effect.AddFrameListener(binding.NewDataListener(func() {
		auto = true
		selector.SetSelected(effect.EffectName())
	}))

	effect.AddFolderListener(binding.NewDataListener(func() {
		// list := []string{iohandler.AsFolder()}
		// list = append(list, effect.ListCurrent()...)
		selector.SetOptions(effect.ListCurrent())
	}))
	return selector
}

func OnEffectChanged(effect *effectio.EffectIo, auto *bool) func(title string) {

	return func(title string) {
		if *auto {
			*auto = false
			return
		}
		effect.Select(title)
	}
}
