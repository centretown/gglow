package ui

import (
	"glow-gui/fields"
	"glow-gui/resources"
	"glow-gui/settings"

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
	model       fields.Model
	theme       *settings.GlowTheme
	sourceStrip binding.Untyped

	strip         *LightStrip
	stripLayout   *LightStripLayout
	stripPlayer   *LightStripPlayer
	stripTools    *fyne.Container
	playContainer *fyne.Container

	effectSelect *widget.Select

	toolbar *SharedTools

	frameEditor *FrameEditor
	layerEditor *LayerEditor

	editor  *fyne.Container
	split   *container.Split
	isSplit bool

	mainContainer *fyne.Container
	isMobile      bool
}

func NewUi(app fyne.App, window fyne.Window, model fields.Model, theme *settings.GlowTheme) *Ui {
	ui := &Ui{
		window:      window,
		app:         app,
		preferences: app.Preferences(),
		theme:       theme,
		model:       model,
		sourceStrip: binding.NewUntyped(),
		isMobile:    app.Driver().Device().IsMobile(),
	}

	window.SetContent(ui.BuildContent())
	return ui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
	ui.preferences.SetString(settings.Effect.String(), ui.model.EffectName())
	ui.preferences.SetBool(settings.ContentSplit.String(), ui.isSplit)
}

func (ui *Ui) layoutContent() *fyne.Container {

	toolsLayout := container.New(layout.NewCenterLayout(), ui.toolbar)
	ui.editor = container.NewBorder(toolsLayout, nil, nil, nil,
		container.NewBorder(ui.frameEditor.Container, nil, nil, nil, ui.layerEditor.Container))

	menuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})
	selectorMenu := container.NewBorder(nil, nil, menuButton, nil, ui.effectSelect)

	if ui.isMobile {
		dropDown := dialog.NewCustom(resources.EditLabel.String(), "hide", ui.editor, ui.window)
		menuButton.OnTapped = func() {
			dropDown.Resize(ui.window.Canvas().Size())
			dropDown.Show()
		}
		ui.mainContainer = container.NewBorder(selectorMenu, nil, nil, nil, ui.playContainer)
	} else {

		ui.isSplit = ui.preferences.BoolWithFallback(settings.ContentSplit.String(), true)
		ui.split = container.NewHSplit(ui.editor, ui.playContainer)

		setSplit := func(isSplit bool) {
			if isSplit {
				ui.mainContainer.Objects = []fyne.CanvasObject{selectorMenu, ui.split}
			} else {
				ui.mainContainer.Objects = []fyne.CanvasObject{selectorMenu, ui.playContainer}
			}
		}

		ui.mainContainer = container.NewBorder(selectorMenu, nil, nil, nil, ui.playContainer)
		setSplit(ui.isSplit)

		menuButton.OnTapped = func() {
			ui.isSplit = !ui.isSplit
			setSplit(ui.isSplit)
		}
	}

	return ui.mainContainer
}

func (ui *Ui) BuildContent() *fyne.Container {
	columns, rows := ui.getLightPreferences()
	color := ui.theme.Color(settings.LightStripBackground, ui.theme.GetVariant())
	ui.strip = NewLightStrip(columns*rows, rows, color)
	ui.sourceStrip.Set(ui.strip)

	ui.stripLayout = NewLightStripLayout(ui.window, ui.app.Preferences(), ui.sourceStrip, color)
	ui.stripPlayer = NewLightStripPlayer(ui.sourceStrip, ui.model, ui.stripLayout)
	ui.stripTools = container.New(layout.NewCenterLayout(), ui.stripPlayer)

	ui.effectSelect = NewEffectSelect(ui.model)

	ui.toolbar = NewSharedTools(ui.model)
	ui.frameEditor = NewFrameEditor(ui.model, ui.window, ui.toolbar)
	ui.layerEditor = NewLayerEditor(ui.model, ui.window, ui.toolbar)
	ui.toolbar.Refresh()

	ui.playContainer = container.NewBorder(ui.stripTools, nil, nil, nil, ui.strip)

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
	columns = ui.preferences.IntWithFallback(settings.StripColumns.String(),
		settings.StripColumnsDefault)
	rows = ui.preferences.IntWithFallback(settings.StripRows.String(),
		settings.StripRowsDefault)
	return
}
