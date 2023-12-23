package ui

import (
	"gglow/effects"
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
	effect      effects.Effect
	layerSelect *widget.Select
	fields      *effects.FrameFields
	rateBounds  *IntEntryBounds
	rateBox     *RangeIntBox
	tools       *FrameTools
}

func NewFrameEditor(effect effects.Effect, window fyne.Window,
	sharedTools *SharedTools) *FrameEditor {

	fe := &FrameEditor{
		effect:      effect,
		layerSelect: NewLayerSelect(effect),
		rateBounds:  RateBounds,
		fields:      effects.NewFrameFields(),
	}

	fe.layerSelect = NewLayerSelect(fe.effect)
	ratelabel := widget.NewLabel(resources.RateLabel.String())
	fe.rateBox = NewRangeIntBox(fe.fields.Interval, fe.rateBounds)
	frm := container.New(layout.NewFormLayout(), ratelabel, fe.rateBox.Container)
	fe.Container = container.NewBorder(nil, fe.layerSelect, nil, nil, frm)

	fe.tools = NewFrameTools(effect, window)
	sharedTools.AddItems(fe.tools.Items()...)
	effect.OnApply(fe.apply)

	fe.fields.Interval.AddListener(binding.NewDataListener(func() {
		frame := fe.effect.GetFrame()
		interval, _ := fe.fields.Interval.Get()
		if interval != int(frame.Interval) {
			fe.effect.SetChanged()
		}
	}))

	fe.effect.AddFrameListener(binding.NewDataListener(fe.setFields))

	return fe
}

func (fe *FrameEditor) setFields() {
	frame := fe.effect.GetFrame()
	fe.fields.FromFrame(frame)
	fe.rateBox.Entry.SetText(strconv.FormatInt(int64(frame.Interval), 10))
}

func (fe *FrameEditor) apply(frame *glow.Frame) {
	fe.fields.ToFrame(frame)
}
