package ui

import (
	"gglow/fyglow/effectio"
	"gglow/glow"
	"gglow/resources"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FrameEditor struct {
	*fyne.Container
	effect *effectio.EffectIo
	// layerSelect *widget.Select
	fields     *effectio.FrameFields
	folderName binding.String
	effectName binding.String
	rateBounds *IntEntryBounds
	rateBox    *RangeIntBox
	isEditing  bool
}

func NewFrameEditor(effect *effectio.EffectIo, window fyne.Window, menu *fyne.Menu) *FrameEditor {

	fe := &FrameEditor{
		effect:     effect,
		folderName: binding.NewString(),
		effectName: binding.NewString(),
		rateBounds: RateBounds,
		fields:     effectio.NewFrameFields(),
	}

	tools := container.NewCenter(NewFrameTools(effect, window, menu))
	folderLabel := widget.NewLabelWithData(fe.folderName)
	effectLabel := widget.NewLabelWithData(fe.effectName)
	ratelabel := widget.NewLabel(resources.RateLabel.String())
	fe.rateBox = NewRangeIntBox(fe.fields.Interval, fe.rateBounds)
	frm := container.New(layout.NewFormLayout(),
		folderLabel, effectLabel,
		ratelabel, fe.rateBox.Container)
	fe.Container = container.NewBorder(tools, nil, nil, nil, frm)

	effect.OnSave(fe.apply)

	fe.fields.Interval.AddListener(binding.NewDataListener(func() {
		frame := fe.effect.GetFrame()
		interval, _ := fe.fields.Interval.Get()
		if interval != int(frame.Interval) {
			fe.setChanged()
		}
	}))

	effect.AddFolderListener(binding.NewDataListener(func() {
		fe.folderName.Set(effect.FolderName())
	}))
	fe.effect.AddFrameListener(binding.NewDataListener(fe.setFields))
	return fe
}

func (fe *FrameEditor) setChanged() {
	if fe.isEditing {
		fe.effect.SetChanged()
	}
}

func (fe *FrameEditor) setFields() {
	fe.isEditing = false
	frame := fe.effect.GetFrame()
	fe.fields.FromFrame(frame)
	fe.rateBox.Entry.SetText(strconv.FormatInt(int64(frame.Interval), 10))
	fe.folderName.Set(fe.effect.FolderName())
	fe.effectName.Set(fe.effect.EffectName())
	fe.isEditing = true
}

func (fe *FrameEditor) apply(frame *glow.Frame) {
	fe.fields.ToFrame(frame)
}
