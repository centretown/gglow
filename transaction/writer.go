package transaction

import (
	"fmt"
	"glow-gui/effects"
)

func WriteDatabase(dataIn effects.IoHandler, dataOut effects.IoHandler) error {
	list := dataIn.KeyList()

	for _, item := range list {
		fmt.Println(item, "item")
		if dataIn.IsFolder(item) {
			items, err := dataIn.RefreshFolder(item)
			if err != nil {
				fmt.Println("dataIn RefreshFolder", err)
				return err
			}

			_, err = dataOut.RefreshFolder(item)
			if err != nil {
				fmt.Println("dataOut RefreshFolder", err)
				return err
			}

			err = WriteFolder(items, dataIn, dataOut)
			if err != nil {
				fmt.Println("dataOut WriteFolder", err)
				return err
			}
			dataIn.Refresh()
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
