package ui

import (
	"gglow/action"
	"gglow/fyglow/effectio"
	"gglow/iohandler"
	"strings"

	"fyne.io/fyne/v2/data/binding"
)

func BuildAction(data binding.BoolTree, effect *effectio.EffectIo) (act *action.Action) {
	newFilter := func(folder string, ids []string) (item action.FilterItem) {
		item.Folder = folder
		item.Effects = make([]string, 0)
		for _, id := range ids {
			val, _ := data.GetValue(id)
			if val {
				item.Effects = append(item.Effects,
					strings.TrimPrefix(id, folder+"/"))
			}
		}
		return
	}

	act = action.NewAction()
	act.Method = "clone"
	act.Input = effect.Accessor
	output := &iohandler.Accessor{
		Driver:   "code",
		Path:     "/home/dave/src/gglow/generated_test_tranactions_7",
		Database: "/home/dave/src/gglow/generated_test_tranactions_7",
	}
	act.Outputs = append(act.Outputs, output)

	folders := data.ChildIDs("")
	for _, id := range folders {
		filter := newFilter(id, data.ChildIDs(id))
		if len(filter.Effects) > 0 {
			act.FilterItems = append(act.FilterItems, filter)
		}
	}

	// out, _ := yaml.Marshal(act)
	// fmt.Println(string(out))
	return
}

func BuildBoolTree(effect *effectio.EffectIo) binding.BoolTree {
	data := binding.NewBoolTree()
	folders, _ := effect.ListFolder(iohandler.DOTS)
	for _, folder := range folders {
		data.Append("", folder.Key, false)
	}

	for _, folder := range folders {
		ls, _ := effect.ListFolder(folder.Key)
		for _, l := range ls {
			if l.Key != iohandler.DOTS {
				val := l.Value + "/" + l.Key
				data.Append(folder.Key, val, false)
			}
		}
	}
	return data
}
