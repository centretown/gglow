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

	layerSettingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {})
	layerSelect := NewLayerSelect(ui.model)

	layerForm := NewLayerForm(ui.model)

	stripColumns := ui.preferences.FloatWithFallback(resources.StripColumns.String(),
		resources.StripColumnsDefault)
	stripRows := ui.preferences.FloatWithFallback(resources.StripRows.String(),
		resources.StripRowsDefault)
	stripInterval := ui.preferences.FloatWithFallback(resources.StripInterval.String(),
		resources.StripIntervalDefault)

	strip := NewLightStrip(stripColumns*stripRows, stripRows, stripInterval)
	ui.sourceStrip.Set(strip)

	ui.stripPlayer = NewLightStripPlayer(ui.sourceStrip, ui.model.Frame)
	stripTools := container.New(layout.NewCenterLayout(), ui.stripPlayer)

	lightStripSettings := NewLightStripSettings(ui.window, ui.app.Preferences(), ui.sourceStrip)
	lightSettingsButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		lightStripSettings.CustomDialog.Resize(layerForm.Size())
		lightStripSettings.CustomDialog.Show()
	})
	effectSelect := NewEffectSelect(ui.model)
	effectBox := container.NewBorder(nil, nil, lightSettingsButton, nil, effectSelect)

	layerBox := container.NewBorder(nil, nil, layerSettingsButton, nil, layerSelect)
	selectors := container.NewVBox(effectBox, layerBox)
	editor := container.NewBorder(selectors, nil, nil, nil,
		layerForm.AppTabs)

	playContainer := container.NewBorder(nil, stripTools, nil, nil, strip)
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
