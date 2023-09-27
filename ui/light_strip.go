package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	minDiameter     float32 = 10
	maxDiameter     float32 = 100
	defaultDiameter float32 = 14
)

type LightStrip struct {
	widget.DisableableWidget
	background *canvas.Rectangle
	colorOff   color.RGBA
	lights     []*canvas.Circle
	length     float64
	rows       float64
	interval   float64
	diameter   float32
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
		diameter:   defaultDiameter,
	}

	strip.colorOff = color.RGBA{48, 24, 16, 255}
	strip.background.CornerRadius = defaultDiameter
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
		circle := strip.newLight(strip.colorOff)
		strip.lights[i] = circle
	}
}

func (strip *LightStrip) newLight(color color.RGBA) (circle *canvas.Circle) {
	circle = canvas.NewCircle(color)
	return
}

func (strip *LightStrip) CreateRenderer() fyne.WidgetRenderer {
	padding := theme.Padding()
	objects := make([]fyne.CanvasObject, 0, len(strip.lights)+1)
	objects = append(objects, strip.background)

	for _, l := range strip.lights {
		objects = append(objects, l)
	}

	rs := lightStripRenderer{
		objects: objects,
		strip:   strip,
		space:   padding * 2,
		padding: padding * 2,
	}

	return &rs
}

type lightStripRenderer struct {
	objects  []fyne.CanvasObject
	strip    *LightStrip
	diameter float32
	space    float32
	padding  float32
}

func (r *lightStripRenderer) getPos(row, col int, xSpace, ySpace, padding float32) (x, y float32) {
	x = float32(col)*(r.diameter+xSpace) + padding
	y = float32(row)*(r.diameter+ySpace) + r.padding
	return
}

func (r *lightStripRenderer) Layout(size fyne.Size) {
	r.strip.BaseWidget.Resize(size)
	r.strip.background.Resize(size)

	rows := int(r.strip.rows)
	cols := int(r.strip.length) / rows

	r.diameter = defaultDiameter
	circleSize := fyne.Size{Width: r.diameter, Height: r.diameter}

	xSpace := size.Width / float32(cols+1)
	xSpace -= r.diameter
	padding := xSpace + r.diameter/2

	for i, light := range r.strip.lights {
		light.Resize(circleSize)
		x, y := r.getPos(i/cols, i%cols, xSpace, r.space, padding)
		light.Move(fyne.Position{X: x, Y: y})
	}
}

func (r *lightStripRenderer) MinSize() (size fyne.Size) {
	rows := int(r.strip.rows)
	cols := int(r.strip.length) / rows
	x, y := r.getPos(rows, cols, r.space, r.space, r.padding)
	size.Width = x
	size.Height = y
	return
}

func (r *lightStripRenderer) Refresh() {
	r.Layout(r.strip.BaseWidget.Size())
	canvas.Refresh(&r.strip.BaseWidget)
	r.strip.background.Refresh()
}

func (r *lightStripRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *lightStripRenderer) Destroy() {}
