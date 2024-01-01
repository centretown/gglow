package fyui

import (
	"gglow/iohandler"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type EffectSelect struct {
	*widget.Select
	effect iohandler.EffectIoHandler
	auto   bool
}

func NewEffectSelect(effect iohandler.EffectIoHandler) *widget.Select {
	fs := &EffectSelect{
		effect: effect,
	}
	fs.Select = widget.NewSelect(effect.ListCurrentFolder(), fs.onChange)
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
		opt := fs.effect.LoadFolder(title)
		fs.Select.SetOptions(opt)
		return
	}
	if fs.auto {
		fs.auto = false
	} else {
		fs.effect.LoadEffect(title)
	}
}
