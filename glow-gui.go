package main

import (
	"fmt"
	"glow-gui/data"
	"glow-gui/res"
	"glow-gui/store"
	"glow-gui/ui"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	app := app.NewWithID(res.AppID)
	icon, err := res.GooseNoirImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	app.Settings().SetTheme(theme.DarkTheme())

	err = store.Setup()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	window := app.NewWindow(res.GlowLabel.String() + " " +
		res.EffectsLabel.String())
	ui := ui.NewUi(app, window, data.NewModel())
	defer ui.OnExit()

	window.SetContent(ui.BuildContent())
	// window.Resize(res.WindowSize)
	window.ShowAndRun()
}
