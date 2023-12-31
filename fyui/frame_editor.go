package fyui

import (
	"gglow/fyio"
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
	effect *fyio.EffectIo
	// layerSelect *widget.Select
	fields     *fyio.FrameFields
	rateBounds *IntEntryBounds
	rateBox    *RangeIntBox
	isEditing  bool
}

func NewFrameEditor(effect *fyio.EffectIo, window fyne.Window) *FrameEditor {

	fe := &FrameEditor{
		effect: effect,
		// layerSelect: NewLayerSelect(effect),
		rateBounds: RateBounds,
		fields:     fyio.NewFrameFields(),
	}

	tools := container.NewCenter(NewFrameTools(effect, window))
	// fe.layerSelect = NewLayerSelect(fe.effect)
	ratelabel := widget.NewLabel(resources.RateLabel.String())
	fe.rateBox = NewRangeIntBox(fe.fields.Interval, fe.rateBounds)
	frm := container.New(layout.NewFormLayout(), ratelabel, fe.rateBox.Container)
	fe.Container = container.NewBorder(tools, nil, nil, nil, frm)
	// fe.Container = container.NewBorder(tools, fe.layerSelect, nil, nil, frm)

	effect.OnSave(fe.apply)

	fe.fields.Interval.AddListener(binding.NewDataListener(func() {
		frame := fe.effect.GetFrame()
		interval, _ := fe.fields.Interval.Get()
		if interval != int(frame.Interval) {
			fe.setChanged()
		}
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
	fe.isEditing = true
}

func (fe *FrameEditor) apply(frame *glow.Frame) {
	fe.fields.ToFrame(frame)
}
