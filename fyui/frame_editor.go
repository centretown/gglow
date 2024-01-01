package fyui

import (
	"gglow/effectio"
	"gglow/glow"
	"gglow/iohandler"
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
	effect      iohandler.EffectIoHandler
	layerSelect *widget.Select
	fields      *effectio.FrameFields
	rateBounds  *IntEntryBounds
	rateBox     *RangeIntBox
	tools       *FrameTools
	isEditing   bool
}

func NewFrameEditor(effect iohandler.EffectIoHandler, window fyne.Window,
	sharedTools *SharedTools) *FrameEditor {

	fe := &FrameEditor{
		effect:      effect,
		layerSelect: NewLayerSelect(effect),
		rateBounds:  RateBounds,
		fields:      effectio.NewFrameFields(),
	}

	fe.layerSelect = NewLayerSelect(fe.effect)
	ratelabel := widget.NewLabel(resources.RateLabel.String())
	fe.rateBox = NewRangeIntBox(fe.fields.Interval, fe.rateBounds)
	frm := container.New(layout.NewFormLayout(), ratelabel, fe.rateBox.Container)
	fe.Container = container.NewBorder(nil, fe.layerSelect, nil, nil, frm)

	fe.tools = NewFrameTools(effect, window)
	sharedTools.AddItems(fe.tools.Items()...)
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
