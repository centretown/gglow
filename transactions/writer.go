package transactions

import (
	"fmt"
	"gglow/iohandler"
)

func (action *Action) WriteDatabase(dataIn iohandler.IoHandler, dataOut iohandler.OutHandler) error {
	list := dataIn.ListCurrentFolder()

	for _, item := range list {
		if dataIn.IsFolder(item) {
			action.AddNote(fmt.Sprintf("add folder %s", item))
			items, err := dataIn.SetFolder(item)
			if err != nil {
				fmt.Println("dataIn RefreshFolder", err)
				return err
			}

			_, err = dataOut.SetFolder(item)
			if err != nil {
				fmt.Println("dataOut RefreshFolder", err)
				return err
			}

			err = action.WriteFolder(item, items, dataIn, dataOut)
			if err != nil {
				fmt.Println("dataOut WriteFolder", err)
				return err
			}
			dataIn.RootFolder()
		}
	}
	return nil
}

func (action *Action) WriteFolder(folder string, items []string, dataIn iohandler.IoHandler,
	dataOut iohandler.OutHandler) error {

	err := dataOut.WriteFolder(folder)
	if err != nil {
		return err
	}

	for _, item := range items {
		if !dataIn.IsFolder(item) {
			action.AddNote(fmt.Sprintf("add effect %s", item))
			frame, err := dataIn.ReadEffect(item)
			if err != nil {
				return err
			}

			err = dataOut.WriteEffect(item, frame)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
