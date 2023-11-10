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

	var options []string
	options = append(options, model.Store.LookUpList()...)
	fs.Select = widget.NewSelect(options, fs.onChange)
	// fs.Select.PlaceHolder = model.EffectName
	// fs.Select.Alignment = fyne.TextAlignCenter
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
		var options = []string{}
		options = append(options, fs.model.Store.LookUpList()...)
		fs.updateList(options)
		return
	}

	fs.model.LoadFrame(frameName)
}

func (fs *EffectSelect) updateList(options []string) {
	fs.Select.Options = options
	fs.Select.Refresh()
}
