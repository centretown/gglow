package store

import (
	"fmt"
	"glow-gui/effects"
	"glow-gui/settings"
	"glow-gui/sqlio"
	"glow-gui/storageio"

	"fyne.io/fyne/v2"
)

const (
	PathHistory int = iota
	FolderHistory
	EffectHistory
)

func DataSource(config *settings.Configuration,
	preferences fyne.Preferences,
	refresh bool) (store effects.IoHandler, err error) {

	if config.Driver == "sqlite" {
		config.Driver = "sqlite3"
	}

	if preferences != nil {
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
	}

	switch config.Driver {
	case "file":
		store, err = storageio.NewStorageHandler(config.Path)

	case "sqlite3", "mysql":
		store, err = sqlio.NewSqlHandler(config.Driver, config.Path)

	default:
		err = fmt.Errorf("undefined storage method %s", config.Driver)
	}

	if err != nil {
		return
	}

	if refresh {
		_, err = store.Refresh()
	}
	return
}
