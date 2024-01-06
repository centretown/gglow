package fyui

import (
	"gglow/fyio"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type EffectSelect struct {
	*widget.Select
	effect *fyio.EffectIo
	auto   bool
}

func NewEffectSelect(effect *fyio.EffectIo) *widget.Select {
	fs := &EffectSelect{
		effect: effect,
	}
	fs.Select = widget.NewSelect(effect.ListCurrentFolder(),
		fs.onChange)

	effect.AddFrameListener(binding.NewDataListener(func() {
		selected := fs.Select.Selected
		if selected != effect.EffectName() {
			fs.auto = true
			fs.Select.SetSelected(effect.EffectName())
		}
	}))

	effect.AddFolderListener(binding.NewDataListener(func() {
		ls := fs.effect.ListCurrentFolder()
		fs.Select.SetOptions(ls)
	}))
	return fs.Select
}

func (fs *EffectSelect) onChange(title string) {
	if fs.auto {
		fs.auto = false
		return
	}
	if fs.effect.IsFolder(title) {
		ls := fs.effect.LoadFolder(title)
		fs.Select.SetOptions(ls)
	} else {
		fs.effect.LoadEffect(title)
	}
}
