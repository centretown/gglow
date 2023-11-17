package ui

import (
	"glow-gui/glow"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ColorPatch struct {
	widget.BaseWidget
	tapped    func() `json:"-"`
	rectangle *canvas.Rectangle
	disabled  bool
	colorHSV  glow.HSV
}

func NewColorPatch() (patch *ColorPatch) {
	var hsv glow.HSV
	hsv.FromColor(theme.DisabledColor())
	patch = NewColorPatchWithColor(hsv, nil)
	patch.disabled = true
	return
}

func NewColorPatchWithColor(hsv glow.HSV, tapped func()) *ColorPatch {
	cp := &ColorPatch{
		rectangle: canvas.NewRectangle(hsv.ToRGB()),
		tapped:    tapped,
	}
	cp.colorHSV = hsv
	cp.ExtendBaseWidget(cp)
	return cp
}

func (cp *ColorPatch) SetTapped(tapped func()) {
	cp.tapped = tapped
}

func (cp *ColorPatch) SetDisabled(b bool) {
	cp.disabled = b
	cp.setFill(theme.DisabledColor())
}

func (cp *ColorPatch) Disabled() bool {
	return cp.disabled
}

func (cp *ColorPatch) GetHSVColor() glow.HSV {
	return cp.colorHSV
}

func (cp *ColorPatch) GetColor() color.Color {
	return cp.rectangle.FillColor
}

func (cp *ColorPatch) Copy(source *ColorPatch) {
	cp.disabled = source.disabled
	cp.colorHSV = source.colorHSV
	if cp.disabled {
		cp.SetDisabled(true)
	} else {
		cp.setFill(cp.colorHSV.ToRGB())
	}
}

func (cp *ColorPatch) SetHSVColor(hsv glow.HSV) {
	cp.disabled = false
	cp.colorHSV = hsv
	cp.setFill(hsv.ToRGB())
}

func (cp *ColorPatch) SetColor(color color.Color) {
	cp.disabled = false
	cp.colorHSV.FromColor(color)
	cp.setFill(color)
}

func (cp *ColorPatch) setFill(color color.Color) {
	cp.rectangle.FillColor = color
	cp.rectangle.Refresh()
}

func (cp *ColorPatch) Tapped(_ *fyne.PointEvent) {
	if cp.tapped != nil {
		cp.tapped()
	}
}

func (cp *ColorPatch) TappedSecondary(_ *fyne.PointEvent) {
}

type patchRenderer struct {
	objects []fyne.CanvasObject
	patch   *ColorPatch
}

func (cp *ColorPatch) CreateRenderer() fyne.WidgetRenderer {
	pr := &patchRenderer{
		objects: []fyne.CanvasObject{cp.rectangle},
		patch:   cp,
	}
	cp.ExtendBaseWidget(cp)
	return pr
}

func (pr *patchRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{pr.patch.rectangle}
}

func (pr *patchRenderer) Destroy() {}

func (pr *patchRenderer) Refresh() {
	pr.patch.rectangle.Refresh()
}

func (pr *patchRenderer) Layout(size fyne.Size) {
	pr.patch.rectangle.Resize(size)
}

func (pr *patchRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(theme.IconInlineSize() + theme.Padding())
}
