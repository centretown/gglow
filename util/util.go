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

const defaultEffectPath = "/home/dave/src/glow-gui/cabinet/json/"

func main() {
	app.NewWithID(resources.AppID)
	fh := storageio.NewStorageHandler(defaultEffectPath)
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
	list := fh.KeyList()

	for _, item := range list {
		fmt.Println(item, "item")
		if fh.IsFolder(item) {
			items, err := fh.RefreshFolder(item)
			if err != nil {
				fyne.LogError("fh RefreshFolder", err)
				return err
			}

			_, err = sqlh.RefreshFolder(item)
			if err != nil {
				fyne.LogError("fh RefreshFolder", err)
				return err
			}

			err = WriteFolder(items, fh, sqlh)
			if err != nil {
				fyne.LogError("write folder", err)
				return err
			}
			fh.Refresh()
		}
	}
	return nil
}

func WriteFolder(list []string, source effects.IoHandler, dest effects.IoHandler) error {
	err := dest.WriteFolder(dest.FolderName())
	if err != nil {
		return err
	}

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
