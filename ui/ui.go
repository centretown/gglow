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

	effectSelect     *widget.Select
	effectToolbar    *FrameTools
	effectMenuButton *widget.Button

	layerSelect      *widget.Select
	layerToolbar     *LayerTools
	layertMenuButton *widget.Button

	layerEditor *LayerEditor

	mainContainer *fyne.Container
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

func (ui *Ui) LayoutContent() *fyne.Container {
	effectBox := container.NewBorder(nil, ui.effectToolbar, ui.effectMenuButton, nil, ui.effectSelect)
	layerBox := container.NewBorder(nil, ui.layerToolbar, ui.layertMenuButton, nil, ui.layerSelect)
	selectors := container.NewVBox(effectBox, layerBox)
	editor := container.NewVBox(selectors, ui.layerEditor.Container)
	// ui.mainContainer = container.NewBorder(editor, nil, nil, nil, ui.playContainer)
	ui.mainContainer = container.NewBorder(nil, nil, editor, nil, ui.playContainer)
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
	ui.effectMenuButton = widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if ui.effectToolbar.Hidden {
			ui.effectToolbar.Show()
		} else {
			ui.effectToolbar.Hide()
		}
	})

	ui.layerSelect = NewLayerSelect(ui.model)
	ui.layerToolbar = NewLayerTools(ui.model)
	ui.layertMenuButton = widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if ui.layerToolbar.Hidden {
			ui.layerToolbar.Show()
		} else {
			ui.layerToolbar.Hide()
		}
	})

	ui.layerEditor = NewLayerEditor(ui.model, ui.window)
	ui.playContainer = container.NewBorder(widget.NewSeparator(), ui.stripTools, nil, nil, ui.strip)
	ui.sourceStrip.AddListener(binding.NewDataListener(func() {
		strip, _ := ui.sourceStrip.Get()
		ui.strip = strip.(*LightStrip)
		ui.playContainer.Objects = []fyne.CanvasObject{ui.stripTools, ui.strip}
		ui.playContainer.Refresh()
		ui.stripPlayer.ResetStrip()
	}))

	return ui.LayoutContent()
}

func (ui *Ui) getLightPreferences() (columns, rows int) {
	columns = ui.preferences.IntWithFallback(resources.StripColumns.String(),
		resources.StripColumnsDefault)
	rows = ui.preferences.IntWithFallback(resources.StripRows.String(),
		resources.StripRowsDefault)
	return
}
