package main

import (
	"glow-gui/data"
	"glow-gui/resources"
	"glow-gui/store"
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
	theme := resources.NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)

	window := app.NewWindow(resources.GlowLabel.String() + " " +
		resources.EffectsLabel.String())

	store := store.NewStore()
	model := data.NewModel(store)
	ui := ui.NewUi(app, window, model, theme)

	window.SetContent(ui.BuildContent())
	window.SetCloseIntercept(func() {
		ui.OnExit()
		size := window.Content().Size()
		preferences.SetInt(resources.ContentWidth.String(), int(size.Width))
		preferences.SetInt(resources.ContentHeight.String(), int(size.Height))
		window.Close()
	})

	width := preferences.IntWithFallback(resources.ContentWidth.String(), 0)
	height := preferences.IntWithFallback(resources.ContentHeight.String(), 0)
	if height > 0 && width > 0 {
		window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	}

	window.Show()
	effect := preferences.StringWithFallback(resources.Effect.String(), "")
	if len(effect) > 0 {
		model.LoadFrame(effect)
	}
	app.Run()
}
