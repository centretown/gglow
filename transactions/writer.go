package transactions

import (
	"fmt"
)

func (action *Action) WriteDatabase(dataIn ActionHandler, dataOut ActionHandler) error {
	list := dataIn.KeyList()

	for _, item := range list {
		if dataIn.IsFolder(item) {
			action.AddNote(fmt.Sprintf("add folder %s", item))
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

			err = action.WriteFolder(items, dataIn, dataOut)
			if err != nil {
				fmt.Println("dataOut WriteFolder", err)
				return err
			}
			dataIn.Refresh()
		}
	}
	return nil
}

func (action *Action) WriteFolder(list []string, source ActionHandler, dest ActionHandler) error {
	err := dest.WriteFolder(dest.FolderName())
	if err != nil {
		return err
	}

	for _, item := range list {
		if !source.IsFolder(item) {
			action.AddNote(fmt.Sprintf("add effect %s", item))
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
