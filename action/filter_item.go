package action

import "gglow/iohandler"

type FilterItem struct {
	Folder  string
	Effects []string
}

func NewFilterItem() *FilterItem {
	fi := &FilterItem{
		Effects: make([]string, 0),
	}
	return fi
}

func BuildFilterItem(folder string, h iohandler.IoHandler) (filter *FilterItem, err error) {
	filter = NewFilterItem()
	filter.Folder = folder
	var effects []string
	effects, err = h.ListEffects(folder)
	if err != nil {
		return filter, err
	}
	for _, effect := range effects {
		if iohandler.IsFolder(effect) {
			continue
		}
		filter.Effects = append(filter.Effects, effect)
	}
	return filter, err
}

func BuildFilterItems(h iohandler.IoHandler) (list []*FilterItem, err error) {
	list = make([]*FilterItem, 0)
	folders, err := h.ListFolders()
	if err != nil {
		return
	}

	for _, folder := range folders {
		var filter *FilterItem
		filter, err = BuildFilterItem(folder, h)
		if err != nil {
			return
		}
		list = append(list, filter)
	}

	return
}
