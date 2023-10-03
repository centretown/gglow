package ui

import (
	"glow-gui/res"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ColorDialog(window fyne.Window) *fyne.Container {
	hueLabel := widget.NewLabel(res.HueLabel.String())
	hueEntry := widget.NewEntry()
	hueEntry.SetPlaceHolder(res.HueLabel.PlaceHolder())

	label2 := widget.NewLabel(res.SaturationLabel.String())
	value2 := widget.NewEntry()
	value2.SetPlaceHolder(res.SaturationLabel.PlaceHolder())

	button1 := widget.NewButton("Choose Color...", func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {

		}, window)
		picker.Advanced = true
		picker.Show()
	})
	grid := container.New(layout.NewVBoxLayout(),
		hueLabel, hueEntry, label2, value2, button1)
	return grid
}
