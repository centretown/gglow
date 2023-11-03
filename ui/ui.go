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

	mainContainer *fyne.Container
	isMobile      bool
}

func NewUi(app fyne.App, window fyne.Window, model *data.Model) *Ui {
	ui := &Ui{
		window:      window,
		app:         app,
		preferences: app.Preferences(),
		model:       model,
		sourceStrip: binding.NewUntyped(),
		isMobile:    app.Driver().Device().IsMobile(),
	}
	return ui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
}

func (ui *Ui) layoutContent() *fyne.Container {
	// sep := widget.NewSeparator()
	// effectBox := container.NewBorder(ui.effectToolbar, sep, nil, nil, ui.effectSelect)

	layerBox := container.NewBorder(ui.layerToolbar, nil, nil, nil, ui.layerSelect)
	editor := container.NewBorder(layerBox, nil, nil, nil, ui.layerEditor.Container)

	var sideBarContainer fyne.CanvasObject = editor
	if ui.isMobile {
		dropDown := widget.NewModalPopUp(editor, ui.window.Canvas())
		toolbarItem := widget.NewToolbarAction(theme.MenuIcon(), func() {
			if dropDown.Hidden {
				dropDown.Show()
			} else {
				dropDown.Hide()
			}
		})
		dropDown.Hide()
		sideBarContainer = dropDown
		tools := widget.NewToolbar(toolbarItem)
		ui.mainContainer = container.NewBorder(nil, nil, tools, nil, ui.playContainer)
		return ui.mainContainer
	}

	menuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if editor.Hidden {
			editor.Show()
		} else {
			editor.Hide()
		}
	})
	// editor.Hide()

	tools := container.NewBorder(nil, nil, menuButton, nil, ui.effectSelect)
	sidebar := container.NewBorder(nil, nil, nil, nil, sideBarContainer)
	ui.mainContainer = container.NewBorder(tools, nil, sidebar, nil, ui.playContainer)

	return ui.mainContainer
}

func (ui *Ui) BuildContent() *fyne.Container {

	columns, rows := ui.getLightPreferences()
	ui.strip = NewLightStrip(columns*rows, rows)
	ui.sourceStrip.Set(ui.strip)

	ui.stripLayout = NewLightStripLayout(ui.window, ui.app.Preferences(), ui.sourceStrip)
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
