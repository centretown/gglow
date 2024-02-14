package glow

import (
	"fmt"
	"image"
	"image/color"
)

type LinearGradient struct {
	Origin      Origin
	Orientation Orientation
	Stops       []color.NRGBA
}

func NewLinearGradient(origin Origin, orientation Orientation, stops []color.NRGBA) *LinearGradient {
	lg := &LinearGradient{
		Origin:      origin,
		Orientation: orientation,
		Stops:       stops,
	}
	return lg
}

type Extent struct {
	Begin, End, Inc int
}

func (lg *LinearGradient) Draw(dst *image.NRGBA) {
	var xext, yext Extent

	b := dst.Rect.Bounds()
	if lg.Origin == TopLeft || lg.Origin == BottomLeft {
		xext = Extent{Begin: b.Min.X, End: b.Max.X, Inc: 1}
	} else {
		xext = Extent{Begin: b.Max.X, End: b.Min.X, Inc: -1}
	}

	if lg.Origin == TopLeft || lg.Origin == TopRight {
		yext = Extent{Begin: b.Min.Y, End: b.Max.Y, Inc: 1}
	} else {
		yext = Extent{Begin: b.Max.Y, End: b.Min.Y, Inc: -1}
	}

	if lg.Stops == nil || len(lg.Stops) == 0 {
		lg.Stops = append(lg.Stops, color.NRGBA{255, 255, 255, 255})
	}

	switch lg.Orientation {
	case Horizontal:
		lg.DrawHorizontal(dst, xext, yext)
	case Vertical:
		lg.DrawVertical(dst, xext, yext)
	case Diagonal:
		lg.DrawDiagonal(dst, xext, yext)
	}
	fmt.Println(lg.Stops)
}

func (lg *LinearGradient) DrawHorizontal(dst *image.NRGBA, xext, yext Extent) {
	var (
		length        = dst.Bounds().Dy()
		delta  *Delta = NewDelta(lg.Stops, length)
	)

	i := 0
	for y := yext.Begin; y != yext.End; y += yext.Inc {
		cc := delta.Point(i)
		for x := xext.Begin; x != xext.End; x += xext.Inc {
			dst.SetNRGBA(x, y, cc)
		}
		i++
	}
}

func (lg *LinearGradient) DrawVertical(dst *image.NRGBA, xext, yext Extent) {
	var (
		length        = dst.Bounds().Dx()
		delta  *Delta = NewDelta(lg.Stops, length)
	)

	i := 0
	for x := xext.Begin; x != xext.End; x += xext.Inc {
		cc := delta.Point(i)
		for y := yext.Begin; y != yext.End; y += yext.Inc {
			dst.SetNRGBA(x, y, cc)
		}
		i++
	}
}

func (lg *LinearGradient) DrawAngle(dst *image.NRGBA, xext, yext Extent) {
	var (
		height, width        = dst.Bounds().Dy(), dst.Bounds().Dx()
		length               = height * width
		delta         *Delta = NewDelta(lg.Stops, length)
	)
	i := 0
	for y := yext.Begin; y != yext.End; y += yext.Inc {
		j := 0
		for x := xext.Begin; x != xext.End; x += xext.Inc {
			cc := delta.Point((i*width + j*height) / 2)
			dst.SetNRGBA(x, y, cc)
			j++
		}
		i++
	}
}

func (lg *LinearGradient) DrawDiagonal(dst *image.NRGBA, xext, yext Extent) {
	var (
		height, width        = dst.Bounds().Dy(), dst.Bounds().Dx()
		length               = height * width
		delta         *Delta = NewDelta(lg.Stops, length)
	)
	i := 0

	rise := height
	run := width

	for y := yext.Begin; y != yext.End; y += yext.Inc {
		j := 0
		for x := xext.Begin; x != xext.End; x += xext.Inc {
			cc := delta.Point((i*run + j*rise) / 2)
			dst.SetNRGBA(x, y, cc)
			j++
		}
		i++
	}
}