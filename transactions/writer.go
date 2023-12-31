package transactions

import (
	"fmt"
	"gglow/iohandler"
)

func (action *Action) isSelected(selections ...string) bool {
	if action.filterMap == nil || len(action.filterMap) == 0 {
		return true
	}

	selectionLength := len(selections)
	if selectionLength == 0 {
		return true
	}

	em, ok := action.filterMap[selections[0]]
	if !ok {
		return false
	}

	if selectionLength < 2 {
		return true
	}

	if len(em) == 0 {
		return true
	}

	_, ok = em[selections[1]]
	return ok
}

func (action *Action) buildFilterMap() {
	action.filterMap = make(map[string]map[string]bool)
	for _, filter := range action.Filters {
		effectMap := make(map[string]bool)
		for _, effect := range filter.Effects {
			effectMap[effect] = false
		}
		action.filterMap[filter.Folder] = effectMap
	}
}

func (action *Action) WriteDatabase(dataIn iohandler.IoHandler, dataOut iohandler.OutHandler) error {
	folders, err := dataIn.RootFolder()
	if err != nil {
		err = fmt.Errorf("dataIn RootFolder %s", err)
		return err
	}

	action.buildFilterMap()

	for _, folder := range folders {
		if dataIn.IsFolder(folder) && action.isSelected(folder) {
			action.AddNote(fmt.Sprintf("add folder %s", folder))
			items, err := dataIn.SetFolder(folder)
			if err != nil {
				err = fmt.Errorf("dataIn SetFolder %s", err)
				return err
			}

			_, err = dataOut.SetFolder(folder)
			if err != nil {
				err = fmt.Errorf("dataOut SetFolder %s", err)
				return err
			}

			err = action.WriteFolder(folder, items, dataIn, dataOut)
			if err != nil {
				err = fmt.Errorf("dataOut WriteFolder %s", err)
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
		if !dataIn.IsFolder(item) && action.isSelected(folder, item) {
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
