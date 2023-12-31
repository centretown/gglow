package fyui

import (
	"gglow/fyio"
	"gglow/fyresource"
	"gglow/resources"
	"gglow/settings"

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
	effect      *fyio.EffectIo
	theme       *fyresource.GlowTheme
	sourceStrip binding.Untyped

	strip         *LightStrip
	stripLayout   *LightStripLayout
	stripPlayer   *LightStripPlayer
	stripTools    *fyne.Container
	playContainer *fyne.Container

	effectSelect *EffectSelect

	frameEditor *FrameEditor
	layerEditor *LayerEditor

	editor  *fyne.Container
	split   *container.Split
	isSplit bool

	mainContainer *fyne.Container
	isMobile      bool
}

func NewUi(app fyne.App, window fyne.Window, effect *fyio.EffectIo, theme *fyresource.GlowTheme) *Ui {
	ui := &Ui{
		window:      window,
		app:         app,
		preferences: app.Preferences(),
		theme:       theme,
		effect:      effect,
		sourceStrip: binding.NewUntyped(),
		isMobile:    app.Driver().Device().IsMobile(),
	}

	AddGlobalShortCut(window,
		&GlobalShortCut{Shortcut: CtrlS,
			Apply: func() {
				// fmt.Println("SAVE")
				effect.SaveEffect()
			},
			Enabled: effect.HasChanged})

	window.SetContent(ui.BuildContent())

	return ui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
	ui.preferences.SetBool(settings.ContentSplit.String(), ui.isSplit)
}

func (ui *Ui) layoutContent() *fyne.Container {

	ui.editor = container.NewBorder(ui.frameEditor.Container, nil, nil, nil,
		ui.layerEditor.Container)

	menuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})
	selectorMenu := container.NewBorder(ui.effectSelect.Folder, nil, menuButton, nil, ui.effectSelect.SelectEntry)

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
	color := ui.theme.Color(fyresource.LightStripBackground, ui.theme.GetVariant())
	ui.strip = NewLightStrip(columns*rows, rows, color)
	ui.sourceStrip.Set(ui.strip)

	ui.stripLayout = NewLightStripLayout(ui.window, ui.app.Preferences(), ui.sourceStrip, color)
	ui.stripPlayer = NewLightStripPlayer(ui.sourceStrip, ui.effect, ui.stripLayout)
	ui.stripTools = container.New(layout.NewCenterLayout(), ui.stripPlayer)

	// ui.effectTitle = binding.NewSprintf()
	ui.effectSelect = NewEffectSelect(ui.effect)

	ui.frameEditor = NewFrameEditor(ui.effect, ui.window)
	ui.layerEditor = NewLayerEditor(ui.effect, ui.window)

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
		fyresource.StripColumnsDefault)
	rows = ui.preferences.IntWithFallback(settings.StripRows.String(),
		fyresource.StripRowsDefault)
	return
}
