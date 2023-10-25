package ui

import (
	"glow-gui/resources"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ColorDialog(window fyne.Window) *fyne.Container {
	hueLabel := widget.NewLabel(resources.HueLabel.String())
	hueEntry := widget.NewEntry()
	hueEntry.SetPlaceHolder(resources.HueLabel.PlaceHolder())

	label2 := widget.NewLabel(resources.SaturationLabel.String())
	value2 := widget.NewEntry()
	value2.SetPlaceHolder(resources.SaturationLabel.PlaceHolder())

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

func NewLabelCheck(labelText string, field binding.Bool) (label *widget.Label, check *widget.Check) {
	label = widget.NewLabel(labelText)
	check = widget.NewCheckWithData("", field)
	return
}
