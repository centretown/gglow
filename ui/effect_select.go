package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type EffectSelect struct {
	*widget.Select
	model *data.Model
}

func NewEffectSelect(model *data.Model) *widget.Select {

	fs := &EffectSelect{
		model: model,
	}

	fs.Select = widget.NewSelect(model.Store.KeyList, fs.onChange)
	model.Frame.AddListener(binding.NewDataListener(func() {
		selected := fs.Select.Selected
		if selected != model.EffectName {
			fs.Select.SetSelected(model.EffectName)
		}
	}))
	return fs.Select
}

func (fs *EffectSelect) onChange(name string) {
	store := fs.model.Store
	if store.IsFolder(name) {
		store.RefreshKeys(name)
		fs.Select.Options = fs.model.Store.KeyList
		fs.Select.Refresh()
		return
	}
	fs.model.LoadFrame(name)
}
