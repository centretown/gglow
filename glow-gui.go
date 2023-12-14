package main

import (
	"flag"
	"fmt"
	"glow-gui/effects"
	"glow-gui/resources"
	"glow-gui/settings"
	"glow-gui/sqlio"
	"glow-gui/storageio"
	"glow-gui/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	data_source_default = "sqlite"
	usage               = "storage method (sqlite, mysql, file)"
	defaultEffectPath   = "/home/dave/src/glow-gui/cabinet/json/"
)

var data_source string

func init() {
	flag.StringVar(&data_source, "s", data_source_default, usage+" (short form)")
	flag.StringVar(&data_source, "storage", data_source_default, usage)
}

func main() {

	flag.Parse()
	fmt.Println("using storage method", data_source)

	app := app.NewWithID(resources.AppID)
	icon, err := resources.DarkGanderImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	//storage
	preferences := app.Preferences()
	store := selectDataSource(preferences)
	effect := effects.NewEffectIo(store, preferences)

	//window
	window := app.NewWindow(resources.GlowLabel.String())
	theme := settings.NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)
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

func selectDataSource(preferences fyne.Preferences) (store effects.IoHandler) {
	switch data_source {
	case "file":
		rootPath := preferences.StringWithFallback(settings.EffectPath.String(),
			defaultEffectPath)
		store = storageio.NewStorageHandler(rootPath)

	case "mysql":
		store = sqlio.NewMySqlHandler()

	default:
		fallthrough
	case "sqlite":
		store = sqlio.NewSqlLiteHandler()
	}
	return
}
