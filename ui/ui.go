package ui

import (
	"fmt"
	"glow-gui/data"
	"glow-gui/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type Ui struct {
	window        fyne.Window
	app           fyne.App
	model         *data.Model
	mainContainer *fyne.Container

	strip       *LightStrip
	stripPlayer *LightStripPlayer

	playContainer *fyne.Container
	frameView     *fyne.Container
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
	err := resources.LoadGridIcons("dark")
	if err != nil {
		fmt.Println(err)
	}

	ui.strip = NewLightStrip(resources.StripLength, resources.StripRows, resources.StripInterval)
	ui.stripPlayer = NewLightStripPlayer(ui.strip, ui.model.Frame)

	stripTools := container.New(layout.NewCenterLayout(), ui.stripPlayer)
	ui.playContainer = container.NewBorder(nil, stripTools, nil, nil, ui.strip)

	selectors := container.NewVBox(NewFrameSelect(ui.model), NewLayerSelect(ui.model))
	form := NewLayerForm(ui.model)
	ui.frameView = container.NewBorder(selectors, nil, nil, nil,
		form.AppTabs)

	ui.mainContainer = container.NewBorder(ui.frameView, nil, nil, nil, ui.playContainer)
	return ui.mainContainer
}
