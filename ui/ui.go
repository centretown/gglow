package ui

import (
	"fmt"
	"glow-gui/data"
	"glow-gui/res"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	window        fyne.Window
	app           fyne.App
	model         *data.Model
	mainContainer *fyne.Container

	title       *widget.Label
	strip       *LightStrip
	stripPlayer *LightStripPlayer
	layerList   *LayerList

	playContainer *fyne.Container
	frameView     *fyne.Container
	layerView     *fyne.Container

	effectsIcon *widget.Icon
	frameIcon   *widget.Icon
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

	ui.effectsIcon = res.NewAppIcon(res.EffectsIcon)
	ui.frameIcon = res.NewAppIcon(res.FrameIcon)

	ui.title = widget.NewLabelWithData(ui.model.Title)
	ui.title.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	titleBox := container.New(layout.NewCenterLayout(),
		container.NewHBox(ui.effectsIcon, ui.title))

	ui.strip = NewLightStrip(res.StripLength, res.StripRows, res.StripInterval)
	ui.stripPlayer = NewLightStripPlayer(ui.strip, ui.model.Frame)

	stripTools := container.New(layout.NewCenterLayout(), ui.stripPlayer)
	ui.playContainer = container.NewVBox(titleBox, ui.strip, stripTools)

	toLayerButton := widget.NewButtonWithIcon("", theme.GridIcon(),
		ui.changeViewLayer)
	frameSelector := container.NewBorder(nil, nil, toLayerButton, nil,
		NewFrameSelect(ui.model))
	ui.layerList = NewLayerList(ui.model, ui.changeViewLayer)
	ui.frameView = container.NewBorder(frameSelector, nil, nil, nil,
		ui.layerList.List)

	toFrameButton := widget.NewButtonWithIcon("", theme.ListIcon(),
		ui.changeViewFrame)
	layerSelector := container.NewBorder(nil, nil, toFrameButton, nil,
		NewLayerSelect(ui.model))
	form := NewLayerForm(ui.model, ui.changeViewFrame)
	ui.layerView = container.NewBorder(layerSelector, nil, nil, nil,
		form.AppTabs)

	ui.mainContainer = container.NewBorder(ui.playContainer, nil, nil, nil, ui.frameView)
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

func (ui *Ui) changeViewLayer() {
	ui.setView(ui.layerView)
}

func (ui *Ui) SetTitle(title string) {
	ui.title.SetText(title)
}

func (ui *Ui) SetWindowTitle(title string) {
	ui.window.SetTitle(title)
}
