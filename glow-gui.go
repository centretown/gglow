package main

import (
	"flag"
	"fmt"
	"glow-gui/effects"
	"glow-gui/resources"
	"glow-gui/settings"
	"glow-gui/store"
	"glow-gui/ui"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var parsed settings.Configuration

func init() {
	settings.ParseCommandLine(&parsed)
}

func main() {

	flag.Parse()
	fmt.Println("using storage method", parsed.Method, parsed.Path)

	app := app.NewWithID(resources.AppID)
	preferences := app.Preferences()
	icon, err := resources.DarkGanderImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	//storage
	store, err := store.DataSource(&parsed, preferences, true)
	if err != nil {
		fyne.LogError("storage", err)
		os.Exit(1)
	}
	effect := effects.NewEffectIo(store, preferences, &parsed)

	//window
	theme := settings.NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)
	window := app.NewWindow(resources.GlowLabel.String())
	ui := ui.NewUi(app, window, effect, theme)

	window.SetCloseIntercept(func() {
		effect.OnExit()
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
	effect.SetActive()
	app.Run()
}
