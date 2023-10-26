package ui

import (
	"glow-gui/data"
	"glow-gui/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	window      fyne.Window
	app         fyne.App
	preferences fyne.Preferences
	model       *data.Model
	sourceStrip binding.Untyped

	mainContainer *fyne.Container
	stripPlayer   *LightStripPlayer
}

func NewUi(app fyne.App, window fyne.Window, model *data.Model) *Ui {
	ui := &Ui{
		window:      window,
		app:         app,
		preferences: app.Preferences(),
		model:       model,
		sourceStrip: binding.NewUntyped(),
	}
	return ui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
}

func (ui *Ui) BuildContent() *fyne.Container {

	columns, rows := ui.getLightPreferences()
	strip := NewLightStrip(columns*rows, rows)
	ui.sourceStrip.Set(strip)

	stripLayout := NewLightStripLayout(ui.window, ui.app.Preferences(), ui.sourceStrip)
	ui.stripPlayer = NewLightStripPlayer(ui.sourceStrip, ui.model.Frame, stripLayout)
	stripTools := container.New(layout.NewCenterLayout(), ui.stripPlayer)

	effectSelect := NewEffectSelect(ui.model)
	effectMenuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})
	effectBox := container.NewBorder(nil, nil, effectMenuButton, nil, effectSelect)

	layerSelect := NewLayerSelect(ui.model)
	layertMenuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})
	layerBox := container.NewBorder(nil, nil, layertMenuButton, nil, layerSelect)

	selectors := container.NewVBox(effectBox, layerBox)

	layerEditor := NewLayerEditor(ui.model, ui.window)
	editor := container.NewBorder(selectors, nil, nil, nil, layerEditor.Container)

	playContainer := container.NewBorder(widget.NewSeparator(), stripTools, nil, nil, strip)
	ui.sourceStrip.AddListener(binding.NewDataListener(func() {
		i, _ := ui.sourceStrip.Get()
		strip = i.(*LightStrip)
		playContainer.Objects = []fyne.CanvasObject{stripTools, strip}
		playContainer.Refresh()
		ui.stripPlayer.ResetStrip()
	}))

	ui.mainContainer = container.NewBorder(editor, nil, nil, nil, playContainer)
	return ui.mainContainer
}

func (ui *Ui) getLightPreferences() (columns, rows int) {
	columns = ui.preferences.IntWithFallback(resources.StripColumns.String(),
		resources.StripColumnsDefault)
	rows = ui.preferences.IntWithFallback(resources.StripRows.String(),
		resources.StripRowsDefault)
	// interval = ui.preferences.FloatWithFallback(resources.StripInterval.String(),
	// 	resources.StripIntervalDefault)
	return
}
