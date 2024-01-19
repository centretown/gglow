package ui

import (
	"gglow/action"
	"gglow/fyglow/effectio"
	"gglow/iohandler"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/yaml.v3"
)

func BuildAction(data binding.BoolTree, effect *effectio.EffectIo, drivers []string, path string) (act *action.Action) {
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

	for _, driver := range drivers {
		output := &iohandler.Accessor{
			Driver:   driver,
			Path:     path,
			Database: path,
		}
		act.Outputs = append(act.Outputs, output)
	}

	folders := data.ChildIDs(binding.DataTreeRootID)
	for _, id := range folders {
		filter := newFilter(id, data.ChildIDs(id))
		if len(filter.Effects) > 0 {
			act.FilterItems = append(act.FilterItems, filter)
		}
	}

	return
}

func BuildBoolTree(effect *effectio.EffectIo) binding.BoolTree {
	data := binding.NewBoolTree()
	folders, _ := effect.ListFolder(iohandler.DOTS)
	for _, folder := range folders {
		data.Append(binding.DataTreeRootID, folder.Key, false)
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

func ShowActionResults(act *action.Action, window fyne.Window) {
	buf, _ := yaml.Marshal(act)
	seg := string(buf)
	rich := widget.NewRichTextWithText(seg)
	scroll := container.NewScroll(rich)
	var title string
	if act.HasErrors() {
		title = "Action has errors. Check the log."
	} else {
		title = "Action was successful!"
	}

	dlg := dialog.NewCustom(title, "Close", scroll, window)
	dlg.Resize(window.Canvas().Size())
	dlg.Show()
}
