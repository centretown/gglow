package main

import (
	"flag"
	"fmt"
	"gglow/fyglow/effectio"
	"gglow/fyglow/resource"
	"gglow/fyglow/ui"
	"gglow/iohandler"
	"gglow/settings"
	"gglow/store"
	"gglow/text"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	pathUsage   = "path to data accessor"
	pathDefault = "accessor.glow"
)

var accessor = &iohandler.Accessor{
	Driver: "sqlite3",
	Path:   "glow.db",
}

var accessPath string

func init() {
	// flag.StringVar(&accessPath, "p", "", pathUsage+" (short form)")
	flag.StringVar(&accessPath, "", "", pathUsage)
}

func main() {
	var err error
	app := app.NewWithID(resource.AppID)
	preferences := app.Preferences()

	storageHandler, accessor := loadStorage(preferences)
	fmt.Println(accessPath, accessor.Driver, accessor.Path)

	app.SetIcon(resource.DarkGander())

	theme := resource.NewGlowTheme(preferences)
	app.Settings().SetTheme(theme)
	window := app.NewWindow(text.GlowLabel.String())
	effect := effectio.NewEffect(storageHandler, preferences, accessor)
	ui := ui.NewUi(app, window, effect, theme)

	window.SetCloseIntercept(func() {
		accessor.Folder = effect.FolderName()
		accessor.Effect = effect.EffectName()
		err = iohandler.SaveAccessor(accessPath, accessor)
		if err != nil {
			fyne.LogError("SaveAccessor", err)
		} else {
			preferences.SetString(settings.AccessFile.String(), accessPath)
		}

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

func loadStorage(preferences fyne.Preferences) (iohandler.IoHandler, *iohandler.Accessor) {
	flag.Parse()
	accessPath = flag.Arg(0)
	if accessPath == "" {
		accessPath = preferences.StringWithFallback(settings.AccessFile.String(), "")
	}

	if accessPath == "" {
		accessPath = pathDefault
	}

	info, err := os.Stat(accessPath)
	if err == nil {
		if info.IsDir() {
			fyne.LogError("loadStorage",
				fmt.Errorf("path '%s' is a folder", accessPath))
			os.Exit(1)
		}
		accessor, err = iohandler.LoadAccessor(accessPath)
		if err != nil {
			fyne.LogError("load accessor file", err)
			os.Exit(1)
		}
	}

	storeHandler, err := store.NewIoHandler(accessor)
	if err == nil {
		_, err = storeHandler.SetRootCurrent()
	}

	if err != nil {
		fyne.LogError("loadStorage", err)
		os.Exit(1)
	}

	return storeHandler, accessor
}
