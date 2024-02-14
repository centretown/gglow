package glow

import (
	"image/color"
)

type DeltaSegment struct {
	C          color.NRGBA
	Begin, End int
	R, G, B, A int
}

func NewDeltaSegment(from, to color.NRGBA) *DeltaSegment {
	return &DeltaSegment{
		C: from,
		R: int(to.R) - int(from.R),
		G: int(to.G) - int(from.G),
		B: int(to.B) - int(from.B),
		A: int(to.A) - int(from.A),
	}
}

func (dlt *DeltaSegment) Point(absolute int) color.NRGBA {
	relative := absolute - dlt.Begin
	length := dlt.End - dlt.Begin
	length += B2I(length == 0)
	// length := dlt.End - dlt.Begin
	return color.NRGBA{
		R: dlt.C.R + uint8(dlt.R*relative/length),
		G: dlt.C.G + uint8(dlt.G*relative/length),
		B: dlt.C.B + uint8(dlt.B*relative/length),
		A: dlt.C.A + uint8(dlt.A*relative/length),
	}
}

type Delta struct {
	segments []*DeltaSegment
	length   int
	count    int
}

func NewDelta(stops []color.NRGBA, length int) *Delta {
	delta := &Delta{length: length}
	var count int = len(stops)

	if count == 0 {
		delta.segments = make([]*DeltaSegment, 0)
		return delta
	}

	if count < 2 {
		delta.segments = make([]*DeltaSegment, 1)
		delta.segments[0] = &DeltaSegment{C: stops[0]}
		delta.count = len(delta.segments)
		return delta
	}

	count--
	var (
		segments = make([]*DeltaSegment, count)
		begin    int
	)

	for i := range segments {
		d := NewDeltaSegment(stops[i], stops[i+1])
		d.Begin = begin
		d.End = (i + 1) * length / count
		segments[i] = d
		begin = d.End + 1
	}
	segments[count-1].End = length

	delta.segments = segments
	delta.count = len(delta.segments)
	return delta
}

func (dlt *Delta) Point(i int) color.NRGBA {
	index := i * dlt.count / dlt.length
	// compensate for rare integer division
	index -= B2I(dlt.segments[index].Begin > i)
	return dlt.segments[index].Point(i)
}
