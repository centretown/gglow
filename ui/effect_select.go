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

	fs.Select = widget.NewSelect(model.Store.LookUpList(), fs.onChange)
	model.Frame.AddListener(binding.NewDataListener(func() {
		selected := fs.Select.Selected
		if selected != model.EffectName {
			fs.Select.SetSelected(model.EffectName)
		}
	}))
	return fs.Select
}

func (fs *EffectSelect) onChange(frameName string) {
	store := fs.model.Store
	if store.IsFolder(frameName) {
		store.RefreshLookupList(frameName)
		fs.updateList(fs.model.Store.LookUpList())
		return
	}

	fs.model.LoadFrame(frameName)
}

func (fs *EffectSelect) updateList(options []string) {
	fs.Select.Options = options
	fs.Select.Refresh()
}
