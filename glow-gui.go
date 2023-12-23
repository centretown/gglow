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

var config settings.Configuration

func init() {
	flag.StringVar(&config.Driver, "s", driverDefault, driverUsage+" (short form)")
	flag.StringVar(&config.Driver, "storage", driverDefault, driverUsage)
	flag.StringVar(&config.Path, "p", pathDefault, pathUsage+" (short form)")
	flag.StringVar(&config.Path, "path", pathDefault, pathUsage)
	flag.StringVar(&config.Folder, "f", folderDefault, folderUsage+" (short form)")
	flag.StringVar(&config.Folder, "folder", folderDefault, folderUsage)
	flag.StringVar(&config.Effect, "e", effectDefault, effectUsage+" (short form)")
	flag.StringVar(&config.Effect, "effect", effectDefault, effectUsage)
}

const (
	driverDefault = "sqlite3"
	driverUsage   = "storage driver (sqlite3, mysql, file)"
	pathUsage     = "path to data"
	pathDefault   = ""
	folderUsage   = "folder to access"
	folderDefault = ""
	effectUsage   = "effect to read"
	effectDefault = ""
)

const (
	PathHistory int = iota
	FolderHistory
	EffectHistory
)

func main() {

	flag.Parse()
	fmt.Println("using storage method", config.Driver, config.Path)

	app := app.NewWithID(resources.AppID)
	preferences := app.Preferences()
	icon, err := resources.DarkGanderImage.Load()
	if err == nil {
		app.SetIcon(icon)
	}

	history := preferences.StringListWithFallback(config.Driver, []string{"", "", ""})
	if config.Path == "" {
		config.Path = history[PathHistory]
	}
	if config.Folder == "" {
		config.Folder = history[FolderHistory]
	}
	if config.Effect == "" {
		config.Effect = history[EffectHistory]
	}

	//storage
	store, err := store.NewHandler(&config)
	if err == nil {
		_, err = store.Refresh()

	}
	if err != nil {
		fyne.LogError("storage", err)
		os.Exit(1)
	}

	effect := effects.NewEffectIo(store, preferences, &config)

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
