package main

import (
	"glow-gui/effects"
	"glow-gui/resources"
	"glow-gui/settings"
	"glow-gui/sqlio"
	"glow-gui/storageio"
	"glow-gui/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const defaultEffectPath = "/home/dave/src/glow-gui/cabinet/json/"

func main() {
	app := app.NewWithID(resources.AppID)

	icon, err := resources.DarkGanderImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	preferences := app.Preferences()
	theme := settings.NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)

	var store effects.IoHandler
	var sqlOn = true

	if sqlOn {
		store = sqlio.NewMySqlHandler()
	} else {
		rootPath := preferences.StringWithFallback(settings.EffectPath.String(),
			defaultEffectPath)
		store = storageio.NewStorageHandler(rootPath)
	}

	eff := effects.NewEffectIo(store, preferences)

	window := app.NewWindow(resources.GlowLabel.String())
	ui := ui.NewUi(app, window, eff, theme)

	window.SetCloseIntercept(func() {
		eff.OnExit()
		ui.OnExit()
		size := window.Canvas().Size()
		preferences.SetInt(settings.ContentWidth.String(), int(size.Width))
		preferences.SetInt(settings.ContentHeight.String(), int(size.Height))
		window.Close()
	})

	width := preferences.IntWithFallback(settings.ContentWidth.String(), 0)
	height := preferences.IntWithFallback(settings.ContentHeight.String(), 0)
	if height > 0 && width > 0 {
		window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	}

	window.Show()
	eff.SetActive()
	app.Run()
}
