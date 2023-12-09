package ui

import (
	"glow-gui/effects"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type EffectSelect struct {
	*widget.Select
	effect effects.Effect
	auto   bool
}

func NewEffectSelect(effect effects.Effect) *widget.Select {
	fs := &EffectSelect{
		effect: effect,
	}
	fs.Select = widget.NewSelect(effect.KeyList(), fs.onChange)
	effect.AddFrameListener(binding.NewDataListener(func() {
		selected := fs.Select.Selected
		if selected != effect.EffectName() {
			fs.auto = true
			fs.Select.SetSelected(effect.EffectName())
		}
	}))
	return fs.Select
}

func (fs *EffectSelect) onChange(title string) {
	if fs.effect.IsFolder(title) {
		fs.auto = false
		fs.Select.SetOptions(fs.effect.RefreshKeys(title))
		return
	}
	if fs.auto {
		fs.auto = false
	} else {
		fs.effect.ReadEffect(title)
	}
}
