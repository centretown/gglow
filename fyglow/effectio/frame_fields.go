package effectio

import (
	"gglow/glow"

	"fyne.io/fyne/v2/data/binding"
)

type FrameFields struct {
	Interval binding.Int
}

func NewFrameFields() *FrameFields {
	fld := &FrameFields{
		Interval: binding.NewInt(),
	}
	return fld
}

func (fld *FrameFields) FromFrame(frame *glow.Frame) {
	fld.Interval.Set(int(frame.Interval))
}

func (fld *FrameFields) ToFrame(frame *glow.Frame) {
	var i int
	i, _ = fld.Interval.Get()
	frame.Interval = uint32(i)
}
