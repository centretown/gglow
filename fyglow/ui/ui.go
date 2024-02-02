package ui

import (
	"gglow/fyglow/effectio"
	"gglow/fyglow/resource"
	"gglow/settings"
	"gglow/text"

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
	effect      *effectio.EffectIo
	theme       *resource.GlowTheme
	sourceStrip binding.Untyped

	strip         *LightStrip
	stripLayout   *LightStripLayout
	stripPlayer   *LightStripPlayer
	stripTools    *fyne.Container
	playContainer *fyne.Container

	frameEditor *FrameEditor
	layerEditor *LayerEditor

	editor    *fyne.Container
	splitView *container.Split
	view      int

	mainContainer *fyne.Container
	isMobile      bool
	mainMenu      *fyne.Menu
	folderName    binding.String
	effectName    binding.String
	tree          *widget.Tree
}

func NewUi(app fyne.App, window fyne.Window, effect *effectio.EffectIo, theme *resource.GlowTheme) *Ui {
	ui := &Ui{
		window:      window,
		app:         app,
		preferences: app.Preferences(),
		theme:       theme,
		effect:      effect,
		sourceStrip: binding.NewUntyped(),
		isMobile:    app.Driver().Device().IsMobile(),
		folderName:  binding.NewString(),
		effectName:  binding.NewString(),
	}

	ui.mainMenu = BuildMenu(effect, window)
	ui.tree = NewTreeSelector(effect)
	window.SetContent(ui.BuildContent())
	return ui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
	ui.preferences.SetInt(settings.ContentView.String(), ui.view)
	ui.preferences.SetFloat(settings.SplitOffset.String(), ui.splitView.Offset)
}

func (ui *Ui) BuildContent() *fyne.Container {
	columns, rows := ui.getLightPreferences()
	color := ui.theme.Color(resource.LightStripBackground, ui.theme.GetVariant())
	ui.strip = NewLightStrip(columns*rows, rows, color)
	ui.sourceStrip.Set(ui.strip)

	ui.stripLayout = NewLightStripLayout(ui.window, ui.app.Preferences(), ui.sourceStrip, color)
	ui.stripPlayer = NewLightStripPlayer(ui.sourceStrip, ui.effect, ui.stripLayout)
	ui.stripTools = container.New(layout.NewCenterLayout(), ui.stripPlayer)

	ui.frameEditor = NewFrameEditor(ui.effect, ui.window, ui.mainMenu)
	ui.layerEditor = NewLayerEditor(ui.effect, ui.window, ui.mainMenu)

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

func (ui *Ui) layoutContent() *fyne.Container {
	ui.editor = container.NewBorder(ui.frameEditor.Container, nil, nil, nil,
		ui.layerEditor.Container)

	editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {})
	listButton := widget.NewButtonWithIcon("", theme.ListIcon(), func() {})
	mainMenu := widget.NewPopUpMenu(ui.mainMenu, ui.window.Canvas())
	menuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})
	menuButton.Importance = widget.HighImportance
	menuButton.OnTapped = func() {
		mainMenu.Move(menuButton.Position().Add(fyne.Delta{DX: 0,
			DY: menuButton.MinSize().Height}))
		mainMenu.Show()
	}

	lbl := binding.NewSprintf("%s/%s", ui.folderName, ui.effectName)
	effectLabel := widget.NewLabelWithData(lbl)
	effectLabel.Importance = widget.SuccessImportance
	ui.effect.AddFrameListener(binding.NewDataListener(func() {
		ui.folderName.Set(ui.effect.FolderName())
		ui.effectName.Set(ui.effect.EffectName())
	}))

	viewTools := container.NewBorder(nil, nil,
		container.NewHBox(menuButton, listButton, editButton), effectLabel)
	if ui.isMobile {
		dropDown := dialog.NewCustom(text.EditLabel.String(), "hide", ui.editor, ui.window)
		editButton.OnTapped = func() {
			dropDown.Resize(ui.window.Canvas().Size())
			dropDown.Show()
		}
		ui.mainContainer = container.NewBorder(viewTools, nil, nil, nil, ui.playContainer)
	} else {
		ui.layoutDesktop(viewTools, editButton, listButton)
	}

	return ui.mainContainer
}

func (ui *Ui) layoutDesktop(selectorMenu fyne.CanvasObject,
	editButton, listButton *widget.Button) {

	ui.view = ui.preferences.IntWithFallback(settings.ContentView.String(), LIST_VIEW)
	ui.splitView = container.NewHSplit(ui.editor, ui.playContainer)
	splitOffset := ui.preferences.FloatWithFallback(settings.SplitOffset.String(), 0)
	ui.splitView.SetOffset(splitOffset)

	ui.mainContainer = container.NewBorder(selectorMenu, nil, nil, nil, ui.playContainer)
	ui.setDesktopView(ui.view, selectorMenu)

	editButton.OnTapped = func() {
		if ui.view == EDIT_VIEW {
			ui.setDesktopView(PLAY_VIEW, selectorMenu)
		} else {
			ui.setDesktopView(EDIT_VIEW, selectorMenu)
		}
	}
	listButton.OnTapped = func() {
		if ui.view == LIST_VIEW {
			ui.setDesktopView(PLAY_VIEW, selectorMenu)
		} else {
			ui.setDesktopView(LIST_VIEW, selectorMenu)
		}
	}
}

const (
	PLAY_VIEW = iota
	EDIT_VIEW
	LIST_VIEW
)

func (ui *Ui) setDesktopView(view int, selectorMenu fyne.CanvasObject) {
	ui.view = view
	switch view {
	case PLAY_VIEW:
		ui.mainContainer.Objects = []fyne.CanvasObject{selectorMenu, ui.playContainer}
		return
	case EDIT_VIEW:
		ui.splitView.Leading = ui.editor
	case LIST_VIEW:
		ui.splitView.Leading = ui.tree
	}
	ui.splitView.Refresh()
	ui.mainContainer.Objects = []fyne.CanvasObject{selectorMenu, ui.splitView}
}

func (ui *Ui) getLightPreferences() (columns, rows int) {
	columns = ui.preferences.IntWithFallback(settings.StripColumns.String(),
		resource.StripColumnsDefault)
	rows = ui.preferences.IntWithFallback(settings.StripRows.String(),
		resource.StripRowsDefault)
	return
}
