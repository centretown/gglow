package ui

import (
	"glow-gui/data"
	"glow-gui/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	window      fyne.Window
	app         fyne.App
	preferences fyne.Preferences
	model       *data.Model
	theme       *resources.GlowTheme
	sourceStrip binding.Untyped

	strip         *LightStrip
	stripLayout   *LightStripLayout
	stripPlayer   *LightStripPlayer
	stripTools    *fyne.Container
	playContainer *fyne.Container

	effectSelect  *widget.Select
	effectToolbar *FrameTools

	layerSelect  *widget.Select
	layerToolbar *LayerTools

	layerEditor *LayerEditor
	editor      *fyne.Container
	split       *container.Split
	isSplit     bool

	mainContainer *fyne.Container
	isMobile      bool
}

func NewUi(app fyne.App, window fyne.Window, model *data.Model, theme *resources.GlowTheme) *Ui {
	ui := &Ui{
		window:      window,
		app:         app,
		preferences: app.Preferences(),
		theme:       theme,
		model:       model,
		sourceStrip: binding.NewUntyped(),
		isMobile:    app.Driver().Device().IsMobile(),
	}
	return ui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
	ui.preferences.SetString(resources.Effect.String(), ui.model.EffectName)
	ui.preferences.SetBool(resources.ContentSplit.String(), ui.isSplit)
}

func (ui *Ui) layoutContent() *fyne.Container {

	layerBox := container.NewBorder(ui.layerToolbar, nil, nil, nil, ui.layerSelect)
	ui.editor = container.NewBorder(layerBox, nil, nil, nil, ui.layerEditor.Container)
	menuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})
	tools := container.NewBorder(nil, nil, menuButton, nil, ui.effectSelect)

	if ui.isMobile {
		dropDown := dialog.NewCustom("edit", "hide", ui.editor, ui.window)
		menuButton.OnTapped = func() {
			dropDown.Resize(ui.window.Canvas().Size())
			dropDown.Show()
		}
		ui.mainContainer = container.NewBorder(tools, nil, nil, nil, ui.playContainer)
	} else {

		ui.isSplit = ui.preferences.BoolWithFallback(resources.ContentSplit.String(), true)
		ui.split = container.NewHSplit(ui.editor, ui.playContainer)

		showSplit := func(isSplit bool) {
			if isSplit {
				ui.mainContainer.Objects = []fyne.CanvasObject{tools, ui.split}
			} else {
				ui.mainContainer.Objects = []fyne.CanvasObject{tools, ui.playContainer}
			}
		}

		ui.mainContainer = container.NewBorder(tools, nil, nil, nil, ui.playContainer)
		showSplit(ui.isSplit)

		menuButton.OnTapped = func() {
			ui.isSplit = !ui.isSplit
			showSplit(ui.isSplit)
		}
	}
	return ui.mainContainer
}

func (ui *Ui) BuildContent() *fyne.Container {
	columns, rows := ui.getLightPreferences()
	color := ui.theme.Color(resources.LightStripBackground, ui.theme.GetVariant())
	ui.strip = NewLightStrip(columns*rows, rows, color)
	ui.sourceStrip.Set(ui.strip)

	ui.stripLayout = NewLightStripLayout(ui.window, ui.app.Preferences(), ui.sourceStrip, color)
	ui.stripPlayer = NewLightStripPlayer(ui.sourceStrip, ui.model.Frame, ui.stripLayout)
	ui.stripTools = container.New(layout.NewCenterLayout(), ui.stripPlayer)

	ui.effectSelect = NewEffectSelect(ui.model)
	ui.effectToolbar = NewFrameTools(ui.model)

	ui.layerSelect = NewLayerSelect(ui.model)
	ui.layerToolbar = NewLayerTools(ui.model)
	ui.layerEditor = NewLayerEditor(ui.model, ui.window, ui.layerToolbar)

	ui.playContainer = container.NewBorder(widget.NewSeparator(), ui.stripTools, nil, nil, ui.strip)
	ui.sourceStrip.AddListener(binding.NewDataListener(func() {
		strip, _ := ui.sourceStrip.Get()
		ui.strip = strip.(*LightStrip)
		ui.playContainer.Objects = []fyne.CanvasObject{ui.stripTools, ui.strip}
		ui.playContainer.Refresh()
		ui.stripPlayer.ResetStrip()
	}))

	return ui.layoutContent()
}

func (ui *Ui) getLightPreferences() (columns, rows int) {
	columns = ui.preferences.IntWithFallback(resources.StripColumns.String(),
		resources.StripColumnsDefault)
	rows = ui.preferences.IntWithFallback(resources.StripRows.String(),
		resources.StripRowsDefault)
	return
}
