package ui

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
)

type ColorPatch struct {
	background *canvas.Rectangle
	color      color.RGBA
}

func NewColorPatch(color color.RGBA) *ColorPatch {
	cp := &ColorPatch{
		color:      color,
		background: canvas.NewRectangle(color),
	}

	// icon := res
	// button := widget.NewButtonWithIcon("", icon, func() {})
	return cp
}
