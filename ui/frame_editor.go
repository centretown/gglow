package ui

import (
	"glow-gui/fields"
	"glow-gui/glow"
	"glow-gui/resources"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FrameEditor struct {
	*fyne.Container
	model       fields.Model
	layerSelect *widget.Select
	fields      *fields.FrameFields
	rateBounds  *IntEntryBounds
	rateBox     *RangeIntBox
	tools       *FrameTools
}

func NewFrameEditor(model fields.Model, window fyne.Window,
	sharedTools *SharedTools) *FrameEditor {

	fe := &FrameEditor{
		model:       model,
		layerSelect: NewLayerSelect(model),
		rateBounds:  RateBounds,
		fields:      fields.NewFrameFields(),
	}

	fe.layerSelect = NewLayerSelect(fe.model)
	ratelabel := widget.NewLabel(resources.RateLabel.String())
	fe.rateBox = NewRangeIntBox(fe.fields.Interval, fe.rateBounds)
	frm := container.New(layout.NewFormLayout(), ratelabel, fe.rateBox.Container)
	fe.Container = container.NewBorder(nil, fe.layerSelect, nil, nil, frm)

	fe.tools = NewFrameTools(model, window)
	sharedTools.AddItems(fe.tools.Items()...)
	model.OnApply(fe.apply)

	fe.fields.Interval.AddListener(binding.NewDataListener(func() {
		frame := fe.model.GetFrame()
		interval, _ := fe.fields.Interval.Get()
		if interval != int(frame.Interval) {
			fe.model.SetChanged()
		}
	}))

	fe.model.AddFrameListener(binding.NewDataListener(fe.setFields))

	return fe
}

func (fe *FrameEditor) setFields() {
	frame := fe.model.GetFrame()
	fe.fields.FromFrame(frame)
	fe.rateBox.Entry.SetText(strconv.FormatInt(int64(frame.Interval), 10))
}

func (fe *FrameEditor) apply(frame *glow.Frame) {
	fe.fields.ToFrame(frame)
}
