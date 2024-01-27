package ui

import (
	"gglow/fyglow/effectio"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type EffectSelect struct {
	*widget.Select
	effect  *effectio.EffectIo
	lookup  map[string]string
	options []string
	auto    bool
}

func NewEffectSelect(effect *effectio.EffectIo) *EffectSelect {
	fs := &EffectSelect{
		effect:  effect,
		options: make([]string, 0),
		lookup:  make(map[string]string),
	}

	fs.Select = widget.NewSelect(effect.ListCurrent(), fs.onChange)

	effect.AddFrameListener(binding.NewDataListener(func() {
		fs.auto = true
		fs.Select.SetSelected(effect.EffectName())
	}))

	effect.AddFolderListener(binding.NewDataListener(func() {
		fs.options = fs.effect.ListCurrent()
		fs.Select.SetOptions(fs.options)
		fs.buildLookup()
	}))
	return fs
}

func (fs *EffectSelect) buildLookup() {
	fs.lookup = make(map[string]string)
	for _, s := range fs.options {
		fs.lookup[strings.ToLower(s)] = s
	}
}

func (fs *EffectSelect) onChange(title string) {
	if fs.auto {
		fs.auto = false
		return
	}

	if fs.effect.IsFolder(title) {
		fs.options = fs.effect.LoadFolder(title)
		fs.buildLookup()
		fs.Select.SetOptions(fs.options)
	} else {
		fs.effect.LoadEffect(title)
	}
}
