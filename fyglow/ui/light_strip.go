package ui

import (
	"gglow/fyglow/resource"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	minStripWidth  float32 = 320
	minStripHeight float32 = 120
	maxRows        int     = 20
	maxCols        int     = 50
)

type LightStrip struct {
	widget.BaseWidget
	background *canvas.Rectangle
	image      *canvas.Image
	colorOff   color.RGBA
	lights     []*canvas.Circle
	length     int
	rows       int
}

func NewLightStrip(length, rows int, background color.Color) *LightStrip {
	strip := &LightStrip{
		background: canvas.NewRectangle(color.Transparent),
		image:      canvas.NewImageFromResource(resource.DarkGander()),
		length:     length,
		rows:       rows,
	}

	strip.colorOff = color.RGBA{48, 24, 16, 0}
	strip.buildLights()
	strip.ExtendBaseWidget(strip)
	return strip
}

func (strip *LightStrip) Length() uint16 {
	return uint16(strip.length)
}
func (strip *LightStrip) Rows() uint16 {
	return uint16(strip.rows)
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
		strip.lights[i] = canvas.NewCircle(strip.colorOff)
	}
}

type lightStripRenderer struct {
	objects []fyne.CanvasObject
	strip   *LightStrip
}

func (strip *LightStrip) CreateRenderer() fyne.WidgetRenderer {
	objects := make([]fyne.CanvasObject, 0, len(strip.lights)+1)
	objects = append(objects, strip.background)

	for _, l := range strip.lights {
		objects = append(objects, l)
	}

	lsr := lightStripRenderer{
		objects: objects,
		strip:   strip,
	}

	return &lsr
}

func (lsr *lightStripRenderer) Layout(size fyne.Size) {
	lsr.strip.background.Resize(size)
	lsr.strip.background.Refresh()

	rows := int(lsr.strip.rows)
	cols := int(lsr.strip.length) / rows

	cellSize := min(size.Width/float32(cols), size.Height/float32(rows))
	cellSize = float32(math.Floor(float64(cellSize)))

	diameter := float32(math.Ceil(float64(cellSize / 2)))
	circleSize := fyne.Size{Width: diameter, Height: diameter}

	pad := theme.InnerPadding()
	xOrigin := (size.Width-cellSize*float32(cols))/2 + pad
	yOrigin := (size.Height-cellSize*float32(rows))/2 + theme.Padding()

	getPos := func(row, col int) (x, y float32) {
		x = float32(col)*cellSize + xOrigin
		y = float32(row)*cellSize + yOrigin
		return x, y
	}

	for i, light := range lsr.strip.lights {
		light.Resize(circleSize)
		x, y := getPos(i/cols, i%cols)
		light.Move(fyne.Position{X: x, Y: y})
	}
}

func (lsr *lightStripRenderer) MinSize() (size fyne.Size) {
	size.Width = minStripWidth
	size.Height = minStripHeight
	return
}

func (lsr *lightStripRenderer) Refresh() {
	lsr.Layout(lsr.strip.BaseWidget.Size())
}

func (lsr *lightStripRenderer) Objects() []fyne.CanvasObject {
	return lsr.objects
}

func (lsr *lightStripRenderer) Destroy() {}
