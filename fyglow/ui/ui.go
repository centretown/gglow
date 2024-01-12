package ui

import (
	"gglow/fyglow/effectio"
	"gglow/fyglow/resource"
	"gglow/resources"
	"gglow/settings"
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

	editor  *fyne.Container
	split   *container.Split
	isSplit bool

	mainContainer *fyne.Container
	isMobile      bool
	mainMenu      *fyne.Menu
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
	}

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

	editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {})
	popUp := widget.NewPopUpMenu(ui.mainMenu, ui.window.Canvas())
	menuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {})
	menuButton.OnTapped = func() {
		popUp.Move(menuButton.Position().Add(fyne.Delta{DX: 0,
			DY: menuButton.MinSize().Height}))
		popUp.Show()
	}

	selectorMenu := container.NewBorder(nil, nil,
		container.NewHBox(menuButton, editButton), nil, ui.effectSelect.SelectEntry)

	if ui.isMobile {
		dropDown := dialog.NewCustom(resources.EditLabel.String(), "hide", ui.editor, ui.window)
		editButton.OnTapped = func() {
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

		editButton.OnTapped = func() {
			ui.isSplit = !ui.isSplit
			setSplit(ui.isSplit)
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

	AddGlobalShortCut(ui.window,
		&GlobalShortCut{Shortcut: CtrlQ,
			Action: exit,
		})
	ui.mainMenu.Items = append(ui.mainMenu.Items, &fyne.MenuItem{IsSeparator: true},
		&fyne.MenuItem{Label: "Quit", Shortcut: CtrlQ, Action: exit})
}

// 	fileMenu := &fyne.Menu{
// 		Label: "File",
// 		Items: []*fyne.MenuItem{
// 			{Label: "New Folder",
// 				Icon:   theme.FolderNewIcon(),
// 				Action: func() { fmt.Println("New Folder") },
// 			},
// 			{Label: "Trash Folder",
// 				Icon:   theme.DeleteIcon(),
// 				Action: func() { fmt.Println("Trash Folder") },
// 			},
// 			{IsSeparator: true},
// 			{Label: "New Effect",
// 				Icon:   theme.ContentAddIcon(),
// 				Action: func() { fmt.Println("New Effect") },
// 			},
// 			{Label: "Save Effect",
// 				Icon:   theme.DocumentSaveIcon(),
// 				Action: func() { fmt.Println("Save Effect") },
// 			},
// 			{Label: "Remove Effect",
// 				Icon:   theme.ContentRemoveIcon(),
// 				Action: func() { fmt.Println("Remove Effect") },
// 			},
// 		},
// 	}
// 	editMenu := &fyne.Menu{
// 		Label: "Edit",
// 		Items: []*fyne.MenuItem{
// 			{Label: "Cut",
// 				Icon:   theme.ContentCutIcon(),
// 				Action: func() { fmt.Println("Cut") },
// 			},
// 			{IsSeparator: true},
// 			{Label: "Copy",
// 				Icon:   theme.ContentCopyIcon(),
// 				Action: func() { fmt.Println("Copy") },
// 			},
// 			{Label: "Paste",
// 				Icon:   theme.ContentPasteIcon(),
// 				Action: func() { fmt.Println("Paste") },
// 			},
// 			{IsSeparator: true},
// 			{Label: "New Layer",
// 				Icon:   theme.ContentAddIcon(),
// 				Action: func() { fmt.Println("New Layer") },
// 			},
// 			{Label: "Remove Layer",
// 				Icon:   theme.ContentRemoveIcon(),
// 				Action: func() { fmt.Println("Remove Layer") },
// 			},
// 		}}
// 	main := fyne.NewMainMenu()
// 	return main
// }
