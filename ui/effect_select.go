package ui

import (
	"glow-gui/data"
	"glow-gui/resources"
	"glow-gui/store"

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
	fs.Select.PlaceHolder = resources.EffectsLabel.PlaceHolder() + "..."
	// fs.Select.Alignment = fyne.TextAlignCenter

	return fs.Select
}

func (fs *EffectSelect) onChange(frameName string) {
	err := fs.model.LoadFrame(frameName)
	if err != nil {
		//todo popup message
		return
	}
}
