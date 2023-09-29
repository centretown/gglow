package ui

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/res"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	frame glow.Frame

	title       *widget.Label
	strip       *LightStrip
	stripPlayer *LightStripPlayer
	layerList   *LayerList

	window fyne.Window
	app    fyne.App

	effectsIcon *widget.Icon
	frameIcon   *widget.Icon
}

func NewUi(app fyne.App, window fyne.Window) *Ui {
	gui := &Ui{
		window: window,
		app:    app,
	}
	return gui
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

	ui.title = widget.NewLabel(res.EffectsLabel.String())
	ui.title.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	titleBox := container.New(layout.NewCenterLayout(),
		container.NewHBox(ui.effectsIcon, ui.title))

	ui.strip = NewLightStrip(res.StripLength, res.StripRows, res.StripInterval)
	ui.stripPlayer = NewLightStripPlayer(ui.strip)

	toolBox := container.New(layout.NewCenterLayout(), ui.stripPlayer)

	selector := container.NewBorder(nil, nil, ui.frameIcon, nil, NewFrameSelect(ui))

	top := container.NewVBox(titleBox, ui.strip, toolBox, selector)

	ui.layerList = NewLayerList(ui.window, &ui.frame)

	return container.NewBorder(top, nil, nil, nil, ui.layerList.List)
}

func (ui *Ui) SetTitle(title string) {
	ui.title.SetText(title)
}

func (ui *Ui) SetWindowTitle(title string) {
	ui.window.SetTitle(title)
}

func (ui *Ui) SetFrame(frame *glow.Frame) {
	ui.frame = *frame
	ui.frame.Setup(ui.strip.Length(),
		ui.strip.Rows(),
		ui.strip.Interval())
}

func (ui *Ui) OnChangeFrame(frameName string) {
	uri, err := store.LookupURI(frameName)
	if err != nil {
		return
	}
	frame := &glow.Frame{}
	err = store.LoadFrameURI(uri, frame)
	if err != nil {
		fyne.LogError(fmt.Sprintf("unable to load frame %s", uri.Name()), err)
		return
	}

	ui.SetFrame(frame)
	ui.SetWindowTitle(fmt.Sprintf("%s %s-%s",
		res.GlowLabel.String(), res.EffectsLabel.String(), frameName))
	ui.layerList.SetFrame(&ui.frame)
	ui.stripPlayer.SetFrame(&ui.frame)
}
