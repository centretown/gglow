package main

import (
	"fmt"
	"glow-gui/effects"
	"glow-gui/resources"
	"glow-gui/sqlio"
	"glow-gui/storageio"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID(resources.AppID)
	preferences := app.Preferences()
	fh := storageio.NewStorageHandler(preferences)
	defer fh.OnExit()

	f := func(fh *storageio.StorageHandler, sqlh *sqlio.SqlHandler) error {
		err := sqlh.CreateNewDatabase()
		if err != nil {
			return err
		}

		err = WriteDatabase(fh, sqlh)
		if err != nil {
			return err
		}
		sqlh.OnExit()
		return nil
	}

	err := f(fh, sqlio.NewSqlLiteHandler())
	if err != nil {
		return
	}
	f(fh, sqlio.NewMySqlHandler())
	if err != nil {
		return
	}
	fmt.Println("Complete!")
}

func WriteDatabase(fh *storageio.StorageHandler, sqlh *sqlio.SqlHandler) error {
	sqlh.Folder = fh.Current.Name()
	list := fh.RootList()

	fmt.Println(sqlh.Folder)
	for _, item := range list {
		if fh.IsFolder(item) {
			items, err := fh.RefreshKeys(item)
			if err != nil {
				fyne.LogError("fh RefreshKeys", err)
				return err
			}

			_, err = sqlh.RefreshKeys(fh.Current.Name())
			if err != nil {
				fyne.LogError("sqlh RefreshKeys", err)
				return err
			}

			err = WriteFolder(items, fh, sqlh)
			if err != nil {
				fyne.LogError("write folder", err)
				return err
			}
			fh.RootList()
		}
	}
	return nil
}

func WriteFolder(list []string, source effects.IoHandler, dest effects.IoHandler) error {
	for _, item := range list {
		if !source.IsFolder(item) {
			frame, err := source.ReadEffect(item)
			if err != nil {
				return err
			}

			err = dest.WriteEffect(item, frame)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
