package ui

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/res"
	"glow-gui/store"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	frame glow.Frame

	mainBox     *fyne.Container
	strip       *LightStrip
	title       *widget.Label
	stripPlayer *LightStripPlayer
	layerList   *widget.List

	window fyne.Window
	app    fyne.App

	effectsIcon *widget.Icon
	frameIcon   *widget.Icon
	layerIcon   *widget.Icon
	gridIcon    *widget.Icon
	chromaIcon  *widget.Icon
}

func NewUi(app fyne.App, window fyne.Window) *Ui {
	gui := &Ui{
		window: window,
		app:    app,
	}
	return gui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
}

func (ui *Ui) BuildContent() *fyne.Container {
	err := res.LoadGridIcons("dark")
	if err != nil {
		fmt.Println(err)
	}

	ui.effectsIcon = res.NewAppIcon(res.EffectsIcon)
	ui.frameIcon = res.NewAppIcon(res.FrameIcon)
	ui.layerIcon = widget.NewIcon(theme.ListIcon())
	ui.gridIcon = widget.NewIcon(theme.GridIcon())
	ui.chromaIcon = widget.NewIcon(theme.ColorChromaticIcon())

	ui.title = widget.NewLabel(res.GlowEffectsLabel.String())
	ui.title.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	titleBox := container.New(layout.NewCenterLayout(),
		container.NewHBox(ui.effectsIcon, ui.title))

	ui.strip = NewLightStrip(res.StripLength, res.StripRows, res.StripInterval)
	ui.stripPlayer = NewLightStripPlayer(ui.strip)

	toolBox := container.New(layout.NewCenterLayout(), ui.stripPlayer)

	selector := container.NewBorder(nil, nil, ui.frameIcon, nil, NewFrameSelect(ui))

	top := container.NewVBox(titleBox, ui.strip, toolBox, selector)

	ui.layerList = NewLayerList(&ui.frame)

	ui.mainBox = container.NewBorder(top, nil, nil, nil, ui.layerList)
	return ui.mainBox
}

func (ui *Ui) SetTitle(title string) {
	ui.title.SetText(title)
}

func (ui *Ui) SetWindowTitle(title string) {
	ui.window.SetTitle(title)
}

func (ui *Ui) SetFrame(frame *glow.Frame) {
	ui.frame = *frame
	ui.frame.Setup(ui.strip.Length(),
		ui.strip.Rows(),
		ui.strip.Interval())
}

func (ui *Ui) OnChangeFrame(frameName string) {
	uri, err := store.LookupURI(frameName)
	if err != nil {
		return
	}
	frame := &glow.Frame{}
	err = store.LoadFrameURI(uri, frame)
	if err != nil {
		return
	}

	ui.SetFrame(frame)
	ui.SetWindowTitle(res.GlowEffectsLabel.String() + " - " + frameName)
	ui.stripPlayer.SetFrame(&ui.frame)
	ui.layerList.Refresh()
}

// func (ui *Ui) LightInfo() *fyne.Container {
// 	format := func(id res.LabelID, min, max float64,
// 		boundf *float64) (slider *widget.Slider,
// 		label *widget.Label, entry *widgetx.NumericalEntry) {

// 		boundValue := binding.BindFloat(boundf)
// 		slider = widget.NewSliderWithData(min, max, boundValue)

// 		entry = widgetx.NewNumericalEntry()

// 		label = widget.NewLabel(id.String())
// 		return
// 	}

// 	grid := container.NewGridWithColumns(3)
// 	slider, label, entry := format(res.LengthLabel, 10, 100, &ui.strip.length)
// 	grid.Add(label)
// 	grid.Add(slider)
// 	grid.Add(entry)

// 	slider, label, entry = format(res.RowsLabel, 1, 10, &ui.strip.rows)
// 	grid.Add(label)
// 	grid.Add(slider)
// 	grid.Add(entry)

// 	slider, label, entry = format(res.IntervalLabel, 1, 1000, &ui.strip.interval)
// 	grid.Add(label)
// 	grid.Add(slider)
// 	grid.Add(entry)

// 	return grid
// }

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
