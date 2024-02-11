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

type extent struct {
	beg, end, inc int
}

type delta struct {
	I, R, G, B, A int
}

func (lg *LinearGradient) Draw(dst *image.NRGBA) {
	var xext, yext extent

	b := dst.Rect.Bounds()
	if lg.Origin == TopLeft || lg.Origin == BottomLeft {
		xext = extent{beg: b.Min.X, end: b.Max.X, inc: 1}
	} else {
		xext = extent{beg: b.Max.X, end: b.Min.X, inc: -1}
	}

	if lg.Origin == TopLeft || lg.Origin == TopRight {
		yext = extent{beg: b.Min.Y, end: b.Max.Y, inc: 1}
	} else {
		yext = extent{beg: b.Max.Y, end: b.Min.Y, inc: -1}
	}

	if lg.Stops == nil || len(lg.Stops) == 0 {
		lg.Stops = append(lg.Stops, color.NRGBA{255, 255, 255, 255})
	}
	lg.DrawHorizontal(dst, xext, yext)
	fmt.Println(lg.Stops)
}

func (lg *LinearGradient) buildDeltas(dst *image.NRGBA, length int) []*delta {
	stopCount := len(lg.Stops)
	if stopCount < 2 {
		deltas := make([]*delta, 1)
		deltas[0] = &delta{I: length}
		return deltas
	}

	deltaCount := stopCount - 1
	ratio := length / deltaCount
	deltas := make([]*delta, deltaCount)
	for i := range deltas {

		from := lg.Stops[i]
		to := lg.Stops[i+1]

		d := &delta{
			I: (i + 1) * ratio,
			R: (int(to.R) - int(from.R)) << 8 / ratio,
			G: (int(to.G) - int(from.G)) << 8 / ratio,
			B: (int(to.B) - int(from.B)) << 8 / ratio,
			A: (int(to.A) - int(from.A)) << 8 / ratio,
		}
		deltas[i] = d
	}

	deltas[len(deltas)-1].I = length
	return deltas
}

func (lg *LinearGradient) DrawHorizontal(dst *image.NRGBA, xext, yext extent) {
	var (
		b      = dst.Bounds()
		length = (b.Max.Y - b.Min.Y)
		deltas []*delta
	)

	slen := len(lg.Stops)
	if slen > 0 {
		deltas = lg.buildDeltas(dst, length)
	}

	var (
		c, cc color.NRGBA
		delta *delta
	)
	deltasLength := length / len(deltas)
	deltaLast := len(deltas) - 1
	cc = lg.Stops[0]

	deltasIndex := 0
	delta = deltas[0]
	for i, y := 0, yext.beg; y < yext.end; y += yext.inc {

		for x := xext.beg; x < xext.end; x += xext.inc {
			dst.SetNRGBA(x, y, cc)
		}

		deltasIndex += B2I(i > delta.I)
		delta = deltas[deltasIndex]
		c = lg.Stops[deltasIndex]
		cc.R = c.R + uint8((y * delta.R >> 8))
		cc.G = c.G + uint8((y * delta.G >> 8))
		cc.B = c.B + uint8((y * delta.B >> 8))

		fmt.Println(i, deltasIndex, deltasLength, yext.end, cc, c)
		i++
	}
	fmt.Println(length, deltasLength, deltaLast)
	for _, d := range deltas {
		fmt.Println(d)
	}
}

func (lg *LinearGradient) DrawVertical(dst *image.NRGBA, xext, yext extent) {
	var (
		b      = dst.Bounds()
		length = (b.Max.X - b.Min.X)
		deltas []*delta
	)

	slen := len(lg.Stops)
	if slen > 0 {
		deltas = lg.buildDeltas(dst, length)
	}

	var (
		c, cc color.NRGBA
		delta *delta
	)

	delta = deltas[0]

	for y := yext.beg; y < yext.end; y += yext.inc {

		deltasIndex := 0
		cc = lg.Stops[0]

		for x := xext.beg; x < xext.end; x += xext.inc {
			dst.SetNRGBA(x, y, cc)

			deltasIndex += B2I(x > delta.I)
			delta = deltas[deltasIndex]
			c = lg.Stops[deltasIndex]
			cc.R = c.R + uint8((x * delta.R >> 8))
			cc.G = c.G + uint8((x * delta.G >> 8))
			cc.B = c.B + uint8((x * delta.B >> 8))
			fmt.Println("x,ldelta,end", x, deltasIndex, cc, c)
		}
	}
}

func (lg *LinearGradient) DrawDiagonal(dst *image.NRGBA, xext, yext extent) {
	for y := yext.beg; y < yext.end; y += yext.inc {
		for x := xext.beg; x < xext.end; x += xext.inc {
		}
	}
}
