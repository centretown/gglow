package ui

import (
	"glow-gui/glow"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LightStrip struct {
	widget.DisableableWidget
	objects    []fyne.CanvasObject
	background *canvas.Rectangle
	frame      *glow.Frame
}

func NewLightStrip(frame *glow.Frame) *LightStrip {
	strip := LightStrip{
		frame:      frame,
		background: canvas.NewRectangle(theme.BackgroundColor()),
	}
	strip.updateFrame()
	return &strip
}

func (strip *LightStrip) SetFrame(frame *glow.Frame) {
	strip.frame = frame
	strip.updateFrame()
}

func (strip *LightStrip) updateFrame() {
	strip.objects = make([]fyne.CanvasObject, strip.frame.Length)
	for i := range strip.objects {
		circle := canvas.NewCircle(color.White)
		circle.StrokeColor = color.Black
		circle.StrokeWidth = 1
		strip.objects[i] = circle
	}
}

func (strip *LightStrip) CreateRenderer() fyne.WidgetRenderer {
	strip.ExtendBaseWidget(strip)
	rs := lightStripRenderer{
		strip: strip,
	}

	return &rs
}

type lightStripRenderer struct {
	strip *LightStrip
}

func (lr *lightStripRenderer) Layout(size fyne.Size) {
	lr.strip.background.Resize(size)
	rows := lr.strip.frame.Rows
	cols := lr.strip.frame.Length / rows
	width := size.Width / float32(cols) / 2
	circleSize := fyne.Size{Width: width, Height: width}

	for i, obj := range lr.strip.objects {
		obj.Resize(circleSize)
		col := uint16(i) % cols
		row := uint16(i) % rows
		obj.Move(fyne.Position{
			X: float32(col) * width * 2,
			Y: float32(row) * width * 2})
	}
}

func (lr *lightStripRenderer) MinSize() fyne.Size {
	return fyne.Size{Width: 1, Height: 1}
}

func (lr *lightStripRenderer) Refresh() {
	lr.strip.background.Refresh()
	canvas.Refresh(&lr.strip.DisableableWidget)
	lr.Layout(lr.strip.DisableableWidget.Size())
}

func (lr *lightStripRenderer) Objects() []fyne.CanvasObject {
	return lr.strip.objects
}

func (lr *lightStripRenderer) Destroy() {}
