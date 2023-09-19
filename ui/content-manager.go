package ui

import (
	"glow-gui/res"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ContentManager struct {
	mainContainer     *fyne.Container
	contentContainers [res.ContentCount]*fyne.Container
	currentID         res.ContentID
	sideBar           *fyne.Container
	appImage          *canvas.Image
}

func (m *ContentManager) BuildContent(window fyne.Window) *fyne.Container {
	m.contentContainers[res.FrameContent] = frameInfo()
	m.contentContainers[res.LayerContent] = layerInfo()
	m.contentContainers[res.ChromaContent] = colorDialog(window)
	m.sideBar = m.buildSideBar()
	m.mainContainer = container.NewBorder(nil, nil, m.sideBar, nil,
		m.contentContainers[m.currentID])
	m.appImage = canvas.NewImageFromFile(res.GooseNoirImage.String())
	return m.mainContainer
}

func (m *ContentManager) setContent(id res.ContentID) {
	m.mainContainer.Remove(m.contentContainers[m.currentID])
	m.currentID = id
	m.mainContainer.Add(m.contentContainers[id])
	m.mainContainer.Refresh()
}

func (m *ContentManager) buildSideBar() *fyne.Container {
	x := widget.NewRadioGroup(res.ContentLabels(), func(s string) {
		switch s {
		case res.FrameContent.String():
			m.setContent(res.FrameContent)
		case res.LayerContent.String():
			m.setContent(res.LayerContent)
		case res.ChromaContent.String():
			m.setContent(res.ChromaContent)
		}
	})

	x.Horizontal = false
	x.Selected = res.FrameContent.String()
	grid := container.NewVBox(x)

	return grid
}

func frameInfo() *fyne.Container {
	lengthLabel := widget.NewLabel(res.LengthLabel.String())
	lengthEntry := widget.NewEntry()
	lengthEntry.SetPlaceHolder(res.LengthLabel.PlaceHolder())

	rowsLabel := widget.NewLabel(res.RowsLabel.String())
	rowsEntry := widget.NewEntry()
	rowsEntry.SetPlaceHolder(res.RowsLabel.PlaceHolder())

	intervalLabel := widget.NewLabel(res.IntervalLabel.String())
	intervalEntry := widget.NewEntry()
	intervalEntry.SetPlaceHolder(res.IntervalLabel.PlaceHolder())
	container := container.New(layout.NewVBoxLayout(),
		lengthLabel, lengthEntry, rowsLabel, rowsEntry, intervalLabel, intervalEntry)
	return container
}

func layerInfo() *fyne.Container {
	lengthLabel := widget.NewLabel(res.LengthLabel.String())
	lengthEntry := widget.NewEntry()
	lengthEntry.SetPlaceHolder(res.LengthLabel.PlaceHolder())

	rowsLabel := widget.NewLabel(res.RowsLabel.String())
	rowsEntry := widget.NewEntry()
	rowsEntry.SetPlaceHolder(res.RowsLabel.PlaceHolder())

	container := container.New(layout.NewVBoxLayout(),
		lengthLabel, lengthEntry, rowsLabel, rowsEntry)
	return container
}

func colorDialog(window fyne.Window) *fyne.Container {
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
