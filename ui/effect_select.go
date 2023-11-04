package ui

import (
	"glow-gui/data"
	"glow-gui/store"

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

	options := store.LookUpList()
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
	err := fs.model.LoadFrame(frameName)
	if err != nil {
		//todo popup message
		return
	}
}
