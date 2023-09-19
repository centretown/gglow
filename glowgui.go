package main

import (
	"glow-gui/res"
	"glow-gui/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID(res.AppID)
	icon, err := res.GooseNoirImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	window := app.NewWindow(res.WindowTitle)
	uiContent := ui.ContentManager{}
	window.SetContent(uiContent.BuildContent(window))

	window.Resize(res.WindowSize)
	window.ShowAndRun()
}
