package ui

import (
	"gglow/fyglow/effectio"
	"gglow/fyglow/resource"
	"gglow/settings"
	"gglow/text"
	"os"

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

	effectSelect *EffectSelect

	frameEditor *FrameEditor
	layerEditor *LayerEditor

	editor    *fyne.Container
	splitView *container.Split
	view      int

	mainContainer *fyne.Container
	isMobile      bool
	mainMenu      *fyne.Menu
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
		mainMenu:    fyne.NewMenu(""),
		tree:        NewTreeSelector(effect),
	}

	window.SetContent(ui.BuildContent())
	return ui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
	ui.preferences.SetInt(settings.ContentView.String(), ui.view)
}

func (ui *Ui) layoutContent() *fyne.Container {
	ui.editor = container.NewBorder(ui.frameEditor.Container, nil, nil, nil,
		ui.layerEditor.Container)

	editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {})
	treeButton := widget.NewButtonWithIcon("", theme.ListIcon(), func() {})
	popUp := widget.NewPopUpMenu(ui.mainMenu, ui.window.Canvas())
	menuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})
	menuButton.OnTapped = func() {
		popUp.Move(menuButton.Position().Add(fyne.Delta{DX: 0,
			DY: menuButton.MinSize().Height}))
		popUp.Show()
	}

	selectorMenu := container.NewBorder(nil, nil,
		container.NewHBox(menuButton, editButton, treeButton), nil, ui.effectSelect.SelectEntry)

	if ui.isMobile {
		dropDown := dialog.NewCustom(text.EditLabel.String(), "hide", ui.editor, ui.window)
		editButton.OnTapped = func() {
			dropDown.Resize(ui.window.Canvas().Size())
			dropDown.Show()
		}
		ui.mainContainer = container.NewBorder(selectorMenu, nil, nil, nil, ui.playContainer)
	} else {

		ui.view = ui.preferences.IntWithFallback(settings.ContentView.String(), 1)
		ui.splitView = container.NewHSplit(ui.editor, ui.playContainer)

		setSplit := func(view int) {
			ui.view = view
			switch view {
			case 0:
				ui.mainContainer.Objects = []fyne.CanvasObject{selectorMenu, ui.playContainer}
				return
			case 1:
				ui.splitView.Leading = ui.editor
			case 2:
				ui.splitView.Leading = ui.tree
			}
			ui.splitView.Refresh()
			ui.mainContainer.Objects = []fyne.CanvasObject{selectorMenu, ui.splitView}
		}

		ui.mainContainer = container.NewBorder(selectorMenu, nil, nil, nil, ui.playContainer)
		setSplit(ui.view)

		editButton.OnTapped = func() {
			if ui.view == 1 {
				setSplit(0)
			} else {
				setSplit(1)
			}
		}
		treeButton.OnTapped = func() {
			if ui.view == 2 {
				setSplit(0)
			} else {
				setSplit(2)
			}
		}
	}

	return ui.mainContainer
}

func (ui *Ui) BuildContent() *fyne.Container {
	columns, rows := ui.getLightPreferences()
	color := ui.theme.Color(resource.LightStripBackground, ui.theme.GetVariant())
	ui.strip = NewLightStrip(columns*rows, rows, color)
	ui.sourceStrip.Set(ui.strip)

	ui.stripLayout = NewLightStripLayout(ui.window, ui.app.Preferences(), ui.sourceStrip, color)
	ui.stripPlayer = NewLightStripPlayer(ui.sourceStrip, ui.effect, ui.stripLayout)
	ui.stripTools = container.New(layout.NewCenterLayout(), ui.stripPlayer)

	ui.effectSelect = NewEffectSelect(ui.effect)

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

	ui.addShortCuts()
	return ui.layoutContent()
}

func (ui *Ui) getLightPreferences() (columns, rows int) {
	columns = ui.preferences.IntWithFallback(settings.StripColumns.String(),
		resource.StripColumnsDefault)
	rows = ui.preferences.IntWithFallback(settings.StripRows.String(),
		resource.StripRowsDefault)
	return
}

func (ui *Ui) addShortCuts() {
	exit := func() {
		os.Exit(0)
	}

	manage := func() {
		manageDialog := NewManager(ui.effect, ui.window)
		manageDialog.Resize(ui.window.Canvas().Size())
		manageDialog.Show()
	}

	AddGlobalShortCut(ui.window,
		&GlobalShortCut{Shortcut: CtrlQ, Action: exit})
	AddGlobalShortCut(ui.window,
		&GlobalShortCut{Shortcut: CtrlM, Action: manage})

	ui.mainMenu.Items = append(ui.mainMenu.Items, &fyne.MenuItem{IsSeparator: true},
		&fyne.MenuItem{Label: "manage",
			Shortcut: CtrlM, Action: manage})
	ui.mainMenu.Items = append(ui.mainMenu.Items, &fyne.MenuItem{IsSeparator: true},
		&fyne.MenuItem{Label: text.QuitLabel.String(),
			Shortcut: CtrlQ, Action: exit})
}
