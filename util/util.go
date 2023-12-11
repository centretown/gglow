package main

import (
	"fmt"
	"glow-gui/resources"
	"glow-gui/storageio"

	"glow-gui/effects"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID(resources.AppID)
	preferences := app.Preferences()
	fh := storageio.NewStorageHandler(preferences)
	eff := effects.NewEffectIo(fh, preferences)

	list := eff.KeyList()
	for _, item := range list {
		fmt.Println(item)
	}
}
