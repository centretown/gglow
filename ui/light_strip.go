package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	minStripWidth  float32 = 320
	minStripHeight float32 = 120
)

type LightStrip struct {
	widget.DisableableWidget
	background *canvas.Rectangle
	colorOff   color.RGBA
	lights     []*canvas.Circle
	length     float64
	rows       float64
	interval   float64
}

func (strip *LightStrip) Length() uint16 {
	return uint16(strip.length)
}
func (strip *LightStrip) Rows() uint16 {
	return uint16(strip.rows)
}
func (strip *LightStrip) Interval() uint32 {
	return uint32(strip.interval)
}

// glow.Light interface
func (strip *LightStrip) Get(i uint16) color.RGBA {
	r, g, b, a := strip.lights[i].FillColor.RGBA()
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// glow.Light interface
func (strip *LightStrip) Set(i uint16, color color.RGBA) {
	c := strip.lights[i]
	c.FillColor = color
	c.Refresh()
}

func NewLightStrip(length, rows, interval float64) *LightStrip {
	strip := &LightStrip{
		background: canvas.NewRectangle(theme.ShadowColor()),
		length:     length,
		rows:       rows,
		interval:   interval,
	}

	strip.colorOff = color.RGBA{48, 24, 16, 255}
	strip.background.CornerRadius = 14
	strip.buildLights()
	strip.ExtendBaseWidget(strip)
	return strip
}

func (strip *LightStrip) TurnOff() {
	for i := range strip.lights {
		l := strip.lights[i]
		l.FillColor = strip.colorOff
		l.Refresh()
	}
}

func (strip *LightStrip) buildLights() {
	strip.lights = make([]*canvas.Circle, int(strip.length))
	for i := range strip.lights {
		circle := canvas.NewCircle(strip.colorOff)
		strip.lights[i] = circle
	}
}

type lightStripRenderer struct {
	objects  []fyne.CanvasObject
	strip    *LightStrip
	diameter float32
	space    float32
	padding  float32
}

func (strip *LightStrip) CreateRenderer() fyne.WidgetRenderer {
	padding := theme.Padding()
	objects := make([]fyne.CanvasObject, 0, len(strip.lights)+1)
	objects = append(objects, strip.background)

	for _, l := range strip.lights {
		objects = append(objects, l)
	}

	lrs := lightStripRenderer{
		objects: objects,
		strip:   strip,
		space:   padding * 2,
		padding: padding * 2,
	}

	return &lrs
}

func (lsr *lightStripRenderer) getPos(row, col int, xSpace, ySpace,
	xPadding, yPadding float32) (x, y float32) {

	x = float32(col)*(lsr.diameter+xSpace) + xPadding
	y = float32(row)*(lsr.diameter+ySpace) + yPadding
	return
}

func (lsr *lightStripRenderer) Layout(size fyne.Size) {
	lsr.strip.BaseWidget.Resize(size) //before or after
	lsr.strip.background.Resize(size) //before or after

	rows := int(lsr.strip.rows)
	cols := int(lsr.strip.length) / rows

	lsr.diameter = lsr.calculateDiameter(size, float32(rows), float32(cols))
	circleSize := fyne.Size{Width: lsr.diameter, Height: lsr.diameter}

	xSpace := size.Width / float32(cols+1)
	xSpace -= lsr.diameter
	xPadding := xSpace + lsr.diameter/2
	ySpace := size.Height / float32(rows)
	ySpace -= lsr.diameter
	yPadding := ySpace + lsr.diameter/2

	for i, light := range lsr.strip.lights {
		light.Resize(circleSize)
		x, y := lsr.getPos(i/cols, i%cols, xSpace, lsr.space, xPadding, yPadding)
		light.Move(fyne.Position{X: x, Y: y})
	}
}

func (lsr *lightStripRenderer) calculateDiameter(size fyne.Size, rows, cols float32) float32 {
	width := size.Width / cols
	height := size.Height / rows
	return min(width, height) / 2
}

func (lsr *lightStripRenderer) MinSize() (size fyne.Size) {
	size.Width = minStripWidth
	size.Height = minStripHeight
	return
}

func (lsr *lightStripRenderer) Refresh() {
	lsr.Layout(lsr.strip.BaseWidget.Size())
	canvas.Refresh(&lsr.strip.BaseWidget)
	lsr.strip.background.Refresh()
}

func (lsr *lightStripRenderer) Objects() []fyne.CanvasObject {
	return lsr.objects
}

func (lsr *lightStripRenderer) Destroy() {}
