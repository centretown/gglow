package ui

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

const (
	minStripWidth  float32 = 240
	minStripHeight float32 = 240
)

type LightStrip struct {
	widget.BaseWidget
	background *canvas.Rectangle
	image      *canvas.Image
	colorOff   color.NRGBA
	lights     *image.NRGBA
	length     int
	rows       int
	cols       int
}

func NewLightStrip(length, rows int, background color.Color) *LightStrip {
	strip := &LightStrip{
		background: canvas.NewRectangle(color.Black),
		length:     length,
		rows:       rows,
		cols:       length / rows,
	}

	strip.colorOff = color.NRGBA{48, 24, 16, 255}
	strip.buildLights()
	strip.image = canvas.NewImageFromImage(strip.lights)
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
func (strip *LightStrip) Get(i uint16) color.NRGBA {
	var x, y int = int(i) % strip.cols, int(i) / strip.cols
	n := strip.lights.NRGBAAt(x, y)
	return n
}

// glow.Light interface
func (strip *LightStrip) Set(i uint16, c color.NRGBA) {
	var x, y int = int(i) % strip.cols, int(i) / strip.cols
	strip.lights.SetNRGBA(x, y, c)
}

func (strip *LightStrip) TurnOff() {
	c := strip.colorOff
	for x := 0; x < strip.cols; x++ {
		for y := 0; y < strip.rows; y++ {
			strip.lights.SetNRGBA(x, y, c)
		}
	}
}

func (strip *LightStrip) buildLights() {
	rect := image.Rect(0, 0, strip.cols, strip.rows)
	strip.lights = image.NewNRGBA(rect)
	strip.TurnOff()
}

type lightStripRenderer struct {
	objects []fyne.CanvasObject
	strip   *LightStrip
}

func (strip *LightStrip) CreateRenderer() fyne.WidgetRenderer {
	lsr := lightStripRenderer{
		objects: []fyne.CanvasObject{strip.background, strip.image},
		strip:   strip,
	}

	return &lsr
}

func (lsr *lightStripRenderer) Layout(size fyne.Size) {
	lsr.strip.background.Resize(size)
	lsr.strip.background.Refresh()
	lsr.strip.image.Resize(size)
	lsr.strip.image.Refresh()
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
