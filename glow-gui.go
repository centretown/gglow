package main

import (
	"glow-gui/control"
	"glow-gui/data"
	"glow-gui/resources"
	"glow-gui/settings"
	"glow-gui/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID(resources.AppID)

	icon, err := resources.DarkGanderImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	preferences := app.Preferences()
	theme := settings.NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)

	store := data.NewStore(preferences)
	model := control.NewController(store)

	window := app.NewWindow(resources.GlowLabel.String())
	ui := ui.NewUi(app, window, model, theme)

	window.SetCloseIntercept(func() {
		store.OnExit()
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

	window.ShowAndRun()
}
