package ui

import (
	"fmt"
	"glow-gui/data"
	"glow-gui/res"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
)

type Ui struct {
	window        fyne.Window
	app           fyne.App
	model         *data.Model
	mainContainer *fyne.Container

	// title       *widget.Label
	strip       *LightStrip
	stripPlayer *LightStripPlayer
	// layerList   *LayerList

	playContainer *fyne.Container
	frameView     *fyne.Container
	// layerView     *fyne.Container

	// effectsIcon *widget.Icon
	// frameIcon   *widget.Icon
}

func NewUi(app fyne.App, window fyne.Window, model *data.Model) *Ui {
	ui := &Ui{
		window: window,
		app:    app,
		model:  model,
	}
	return ui
}

func (ui *Ui) OnExit() {
	ui.stripPlayer.OnExit()
}

func (ui *Ui) BuildContent() *fyne.Container {
	err := res.LoadGridIcons("dark")
	if err != nil {
		fmt.Println(err)
	}

	ui.strip = NewLightStrip(res.StripLength, res.StripRows, res.StripInterval)
	ui.stripPlayer = NewLightStripPlayer(ui.strip, ui.model.Frame)

	stripTools := container.New(layout.NewCenterLayout(), ui.stripPlayer)
	// ui.playContainer = container.NewVBox(titleBox, ui.strip, stripTools)
	ui.playContainer = container.NewBorder(nil, stripTools, nil, nil, ui.strip)

	selectors := container.NewVBox(NewFrameSelect(ui.model), NewLayerSelect(ui.model))
	form := NewLayerForm(ui.model, ui.changeViewFrame)
	ui.frameView = container.NewBorder(selectors, nil, nil, nil,
		form.AppTabs)

	ui.mainContainer = container.NewBorder(nil, ui.frameView, nil, nil, ui.playContainer)
	ui.AddListeners()
	return ui.mainContainer
}

func (ui *Ui) AddListeners() {
	ui.model.Title.AddListener(binding.NewDataListener(func() {
		title := ui.model.GetTitle()
		ui.SetWindowTitle(fmt.Sprintf("%s %s-%s",
			res.GlowLabel.String(), res.EffectsLabel.String(), title))
	}))

	// ui.model.Frame.AddListener(binding.NewDataListener(func() {
	// 	face, _ := ui.model.Frame.Get()
	// 	if face != nil {
	// 		frame := face.(*glow.Frame)
	// 		ui.stripPlayer.SetFrame(frame)
	// 	}
	// }))
}

func (ui *Ui) setView(ctr *fyne.Container) {
	ui.mainContainer.Objects = []fyne.CanvasObject{
		ui.playContainer, ctr}
}

func (ui *Ui) changeViewFrame() {
	ui.setView(ui.frameView)
}

// func (ui *Ui) changeViewLayer() {
// 	ui.setView(ui.layerView)
// }

// func (ui *Ui) SetTitle(title string) {
// 	ui.title.SetText(title)
// }

func (ui *Ui) SetWindowTitle(title string) {
	ui.window.SetTitle(title)
}
