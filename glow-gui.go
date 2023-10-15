package main

import (
	"fmt"
	"glow-gui/data"
	"glow-gui/resources"
	"glow-gui/store"
	"glow-gui/ui"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	app := app.NewWithID(resources.AppID)
	icon, err := resources.GooseNoirImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	pref := app.Preferences()
	pref.FloatWithFallback("StripLength", resources.StripLength)
	pref.FloatWithFallback("StripRows", resources.StripRows)
	pref.FloatWithFallback("StripInterval", resources.StripInterval)
	// pref.StringWithFallback("ThemeVariant", "DarkTheme")

	app.Settings().SetTheme(theme.DarkTheme())

	err = store.Setup()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	window := app.NewWindow(resources.GlowLabel.String() + " " +
		resources.EffectsLabel.String())
	ui := ui.NewUi(app, window, data.NewModel())
	defer ui.OnExit()

	window.SetContent(ui.BuildContent())
	// window.Resize(res.WindowSize)
	window.ShowAndRun()
}
