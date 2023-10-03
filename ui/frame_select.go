package ui

import (
	"glow-gui/data"
	"glow-gui/res"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type FrameSelect struct {
	*widget.Select
	model *data.Model
}

func NewFrameSelect(model *data.Model) *widget.Select {

	fs := &FrameSelect{
		model: model,
	}

	options := store.LookUpList()
	fs.Select = widget.NewSelect(options, fs.onChange)
	fs.Select.PlaceHolder = res.ChooseEffectLabel.PlaceHolder()
	fs.Select.Alignment = fyne.TextAlignCenter

	return fs.Select
}

func (fs *FrameSelect) onChange(frameName string) {
	err := fs.model.LoadFrame(frameName)
	if err != nil {
		//todo popup message
		return
	}
}
